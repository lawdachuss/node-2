package proxy

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sardanioss/httpcloak"
)

// ProxyResult represents a tested proxy.
type ProxyResult struct {
	URL     string
	EgressIP string
	Country string
	OK      bool
}

// proxySources are URLs that return plain-text lists of SOCKS5 proxies.
var proxySources = []string{
	"https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&format=text&protocol=socks5&country=nl&timeout=5000",
	"https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&format=text&protocol=socks5&country=in&timeout=5000",
	"https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&format=text&protocol=socks5&country=de&timeout=5000",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",
	"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
	"https://cdn.jsdelivr.net/gh/proxyscrape/free-proxy-list@main/proxies/countries/nl/socks5/data.txt",
}

// cfg cache
var (
	cachedProxies []ProxyResult
	cacheTime     time.Time
	cacheMu       sync.Mutex
	cacheTTL      = 5 * time.Minute
)

// FetchProxies fetches SOCKS5 proxies from public lists, tests them using
// httpcloak (Chrome 146), and returns working proxies sorted by preference
// (NL first, then IN, DE, others).
func FetchProxies(ctx context.Context, limit int) ([]ProxyResult, error) {
	cacheMu.Lock()
	if len(cachedProxies) > 0 && time.Since(cacheTime) < cacheTTL {
		result := append([]ProxyResult{}, cachedProxies...)
		cacheMu.Unlock()
		return result, nil
	}
	cacheMu.Unlock()

	allURLs := fetchAllProxyURLs(ctx)
	if len(allURLs) == 0 {
		return nil, fmt.Errorf("no proxies fetched from any source")
	}

	// Only test a subset to avoid 5+ minute runs
	if len(allURLs) > 100 {
		allURLs = allURLs[:100]
	}

	fmt.Printf("[proxy] testing %d proxies for liveness...\n", len(allURLs))

	var mu sync.Mutex
	var alive []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for _, u := range allURLs {
		wg.Add(1)
		sem <- struct{}{}
		go func(proxyURL string) {
			defer wg.Done()
			defer func() { <-sem }()
			egress := checkProxyEgress(ctx, proxyURL)
			if egress != "" {
				mu.Lock()
				alive = append(alive, proxyURL)
				mu.Unlock()
			}
		}(u)
	}
	wg.Wait()

	fmt.Printf("[proxy] %d proxies alive, testing Chaturbate reachability...\n", len(alive))

	var results []ProxyResult
	var resultsMu sync.Mutex
	var wg2 sync.WaitGroup
	sem2 := make(chan struct{}, 10)

	for _, u := range alive {
		wg2.Add(1)
		sem2 <- struct{}{}
		go func(proxyURL string) {
			defer wg2.Done()
			defer func() { <-sem2 }()
			r := testChaturbateReachability(ctx, proxyURL)
			if r.OK {
				resultsMu.Lock()
				results = append(results, r)
				resultsMu.Unlock()
			}
		}(u)
	}
	wg2.Wait()

	fmt.Printf("[proxy] %d proxies reach Chaturbate\n", len(results))

	results = sortProxies(results)

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	cacheMu.Lock()
	cachedProxies = append([]ProxyResult{}, results...)
	cacheTime = time.Now()
	cacheMu.Unlock()

	return results, nil
}

func fetchAllProxyURLs(ctx context.Context) []string {
	client := &http.Client{Timeout: 15 * time.Second}
	seen := make(map[string]bool)
	var all []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, source := range proxySources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(reqCtx, "GET", src, nil)
			if err != nil {
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			lines := strings.Split(string(body), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				var proxyURL string
				if strings.HasPrefix(line, "socks5://") {
					proxyURL = line
				} else if strings.Contains(line, ":") {
					proxyURL = "socks5://" + line
				} else {
					continue
				}
				mu.Lock()
				if !seen[proxyURL] {
					seen[proxyURL] = true
					all = append(all, proxyURL)
				}
				mu.Unlock()
			}
		}(source)
	}
	wg.Wait()

	// Shuffle for random ordering
	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	return all
}

func checkProxyEgress(ctx context.Context, proxyURL string) string {
	opts := []httpcloak.Option{
		httpcloak.WithTimeout(10 * time.Second),
		httpcloak.WithProxy(proxyURL),
	}
	client := httpcloak.New("chrome-146-windows", opts...)
	if c, ok := interface{}(client).(interface{ Close() error }); ok {
		defer c.Close()
	}
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := client.Do(reqCtx, &httpcloak.Request{
		Method: "GET",
		URL:    "https://api.ipify.org",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
}

func testChaturbateReachability(ctx context.Context, proxyURL string) ProxyResult {
	opts := []httpcloak.Option{
		httpcloak.WithTimeout(15 * time.Second),
		httpcloak.WithProxy(proxyURL),
	}
	client := httpcloak.New("chrome-146-windows", opts...)
	if c, ok := interface{}(client).(interface{ Close() error }); ok {
		defer c.Close()
	}

	// First get the egress IP
	egressIP := checkProxyEgress(ctx, proxyURL)

	// Test Chaturbate homepage
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	resp, err := client.Do(reqCtx, &httpcloak.Request{
		Method: "GET",
		URL:    "https://chaturbate.com",
		Timeout: 15 * time.Second,
	})
	if err != nil {
		return ProxyResult{URL: proxyURL, EgressIP: egressIP, OK: false}
	}
	defer resp.Body.Close()

	statusOK := resp.StatusCode == 200

	// Check for face-id / captcha redirect
	blocked := false
	if loc, ok := resp.Headers["Location"]; ok && len(loc) > 0 {
		lower := strings.ToLower(loc[0])
		if strings.Contains(lower, "verify") || strings.Contains(lower, "captcha") ||
			strings.Contains(lower, "face") || strings.Contains(lower, "human") {
			blocked = true
		}
	}

	return ProxyResult{
		URL:      proxyURL,
		EgressIP: egressIP,
		Country:  lookupCountry(egressIP),
		OK:       statusOK && !blocked,
	}
}

func lookupCountry(ip string) string {
	if ip == "" {
		return ""
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=countryCode")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
}

func sortProxies(proxies []ProxyResult) []ProxyResult {
	preferred := map[string]int{"NL": 0, "IN": 1, "DE": 2}
	sorted := append([]ProxyResult{}, proxies...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			pi := preferred[sorted[i].Country]
			pj := preferred[sorted[j].Country]
			if pj < pi || (pj == pi && sorted[i].Country == "" && sorted[j].Country != "") {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

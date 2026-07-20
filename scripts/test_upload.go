//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/teacat/chaturbate-dvr/entity"
	"github.com/teacat/chaturbate-dvr/server"
	"github.com/teacat/chaturbate-dvr/uploader"
)

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		v = strings.Trim(v, "'\"")
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
	return s.Err()
}

func splitCS(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func mask(s string) string {
	if s == "" {
		return "<empty>"
	}
	if len(s) < 8 {
		return "<too-short>"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

type testLogger struct{}

func (testLogger) Info(format string, a ...any)  { log.Printf("INFO: "+format, a...) }
func (testLogger) Error(format string, a ...any) { log.Printf("ERROR: "+format, a...) }

func main() {
	envPath := ".env"
	if len(os.Args) > 1 {
		envPath = os.Args[1]
	}
	if err := loadDotEnv(envPath); err != nil {
		log.Fatalf("loading %s: %v", envPath, err)
	}

	videoPath := "videos/completed/357054_medium.mp4"
	if len(os.Args) > 2 {
		videoPath = os.Args[2]
	}
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		log.Fatalf("video not found: %s", videoPath)
	}

	server.Config = &entity.Config{
		VoeSXAPIKey:        os.Getenv("VOESX_API_KEY"),
		StreamtapeLogin:    os.Getenv("STREAMTAPE_LOGIN"),
		StreamtapeKey:      os.Getenv("STREAMTAPE_API_KEY"),
		MixdropEmail:       os.Getenv("MIXDROP_EMAIL"),
		MixdropToken:       os.Getenv("MIXDROP_KEY"),
		SeekStreamingKey:   os.Getenv("SEEKSTREAMING_KEY"),
		UpnshareKeys:       splitCS(os.Getenv("UPNSHARE_KEY")),
		NetuAPIKey:         os.Getenv("NETU_API_KEY"),
		LobFileAPIKey:      os.Getenv("LOBFILE_API_KEY"),
	}

	log.Printf("Testing upload of %s", videoPath)
	log.Printf("Configured hosts: GoFile (always), VOE.sx=%s Streamtape=%s Mixdrop=%s SeekStreaming=%s UPnShare=%s Netu=%s LobFile=%s",
		mask(server.Config.VoeSXAPIKey),
		mask(server.Config.StreamtapeKey),
		mask(server.Config.MixdropToken),
		mask(server.Config.SeekStreamingKey),
		fmt.Sprintf("%d key(s)", len(server.Config.UpnshareKeys)),
		mask(server.Config.NetuAPIKey),
		mask(server.Config.LobFileAPIKey),
	)

	upl := uploader.NewMultiHostUploader(
		server.Config.VoeSXAPIKey,
		server.Config.StreamtapeLogin,
		server.Config.StreamtapeKey,
		server.Config.MixdropEmail,
		server.Config.MixdropToken,
		server.Config.SeekStreamingKey,
		testLogger{},
		server.Config.UpnshareKeys,
		server.Config.LobFileAPIKey,
		server.Config.NetuAPIKey,
	)

	available := upl.AvailableHosts()
	log.Printf("Active hosts: %v", available)

	results := upl.UploadToAll(videoPath)
	allOK := true
	for _, r := range results {
		if r.Error != nil {
			allOK = false
			log.Printf("  [FAIL] %-14s %v", r.Host, r.Error)
		} else {
			log.Printf("  [ OK ] %-14s %s", r.Host, r.DownloadLink)
		}
	}
	if len(results) == 0 {
		log.Println("No hosts were active — set at least one API key.")
	}
	if allOK {
		fmt.Println("\nAll active hosts succeeded.")
	} else {
		fmt.Println("\nSome hosts failed (see above).")
	}
}

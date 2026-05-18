#!/bin/sh
apk add --no-cache curl jq
echo '[COOKIE] Starting Cloudflare cookie refresher...'

# Wait for Byparr load balancer to be healthy
until curl -fsS --max-time 5 http://byparr-lb/health > /dev/null 2>&1; do
  echo '[COOKIE] Waiting for Byparr load balancer...'
  sleep 5
done

echo '[COOKIE] Byparr load balancer is up'

# Wait for Byparr backend to respond before attempting API calls.
# We probe the /v1 endpoint with a lightweight request because the nginx
# /backend-health just proxies to Byparr's root / which returns empty.
echo '[COOKIE] Waiting for Byparr backend to accept requests...'
for i in $(seq 1 60); do
  if curl -sS --max-time 10 -X POST http://byparr-lb/v1 \
    -H 'Content-Type: application/json' \
    -d '{"cmd":"request.get","url":"https://chaturbate.com","maxTimeout":5}' > /dev/null 2>&1; then
    echo '[COOKIE] Byparr backend is accepting requests'
    break
  fi
  if [ "$i" -eq 60 ]; then
    echo '[COOKIE] Byparr backend not ready yet, will retry in main loop...'
  fi
  sleep 5
done

# Use shorter timeout for initial attempts so we fail fast when Byparr
# is still provisioning browser instances.
BYPAAR_TIMEOUT=30
attempt_num=0

while true; do
  attempt_num=$((attempt_num + 1))
  # Increase timeout after first few attempts to give Byparr more time
  if [ "$attempt_num" -gt 5 ]; then
    BYPAAR_TIMEOUT=90
  fi

  echo "[COOKIE] Getting cf_clearance from Byparr... ($(date -u +%H:%M:%S), attempt $attempt_num, timeout ${BYPAAR_TIMEOUT}s)"
  RESPONSE=$(curl -sS --fail --max-time $((BYPAAR_TIMEOUT + 10)) -X POST http://byparr-lb/v1 \
    -H 'Content-Type: application/json' \
    -d "{\"cmd\":\"request.get\",\"url\":\"https://chaturbate.com\",\"maxTimeout\":${BYPAAR_TIMEOUT}}")
  CF_COOKIE=$(echo "$RESPONSE" | jq -r '.solution.cookies[] | select(.name=="cf_clearance" or .name=="csrftoken") | .name + "=" + .value' 2>/dev/null | paste -sd ';' -)
  CF_USER_AGENT=$(echo "$RESPONSE" | jq -r '.solution.userAgent // empty' 2>/dev/null)
  if [ -n "$CF_COOKIE" ]; then
    echo '[COOKIE] Refreshed cookies (cf_clearance + csrftoken when present)'

    if [ -n "$CF_USER_AGENT" ]; then
      body=$(jq -n --arg cookies "$CF_COOKIE" --arg ua "$CF_USER_AGENT" '{cookies:$cookies, user_agent:$ua}')
    else
      body=$(jq -n --arg cookies "$CF_COOKIE" '{cookies:$cookies}')
    fi

    # Retry pushing to chaturbate-dvr for up to 5 minutes
    PUSHED=false
    for i in $(seq 1 60); do
      HTTP_CODE=$(curl -sS -o /tmp/cookie-push.json -w '%{http_code}' --max-time 15 -X POST http://chaturbate-dvr:8080/update_config \
        -H 'Content-Type: application/json' \
        -d "$body" 2>/dev/null || echo "000")
      if [ "$HTTP_CODE" = "200" ]; then
        echo '[COOKIE] Pushed to chaturbate-dvr (ok)'
        PUSHED=true
        break
      fi
      sleep 5
    done

    if [ "$PUSHED" != "true" ]; then
      echo '[COOKIE] Could not push cookies to chaturbate-dvr after 5 min, will retry in 30 min'
      cat /tmp/cookie-push.json 2>/dev/null || true
    fi
    sleep 1800
  else
    echo '[COOKIE] Failed to get cf_clearance, retrying in 15 seconds...'
    sleep 15
  fi
done

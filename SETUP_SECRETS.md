# GitHub Secrets Setup Guide

## Step-by-Step Instructions

### 1. Go to Your Repository Settings
1. Open https://github.com/vasud3v/record
2. Click **Settings** (top menu)
3. Click **Secrets and variables** → **Actions** (left sidebar)
4. Click **New repository secret**

---

### 2. Add SETTINGS_JSON Secret

**Name:** `SETTINGS_JSON`

**Value:** Copy and paste this ENTIRE content:

```json
{
  "cookies": "jZcKVhbRNIWIisSY0xRsgC_2J1RUc6BnEl.V7CdohIw-1777442158-1.2.1.1-tLXuiJSDY6mfuZ1_rP_i_OCKuYoK4IPTjESVDag73.qdehisLULg2WXXp_ui_GRv4YXjBsDU9Gl3I.AAX79Ka1R0W2hQfY7XIBNn_dnDNf_PbK6jJs2n5ixR5EycKo6BaEODQI30i0oFJY6YAhNb6dDN9tsT__AyMQsrNCpFqumvYNDACYrGadOfyi4T4YTqkWkMyscwEtkNTLt6QHAW5XZNIr7PdV0X9FN8TACzG2udofsiFJeadZNO7r2W24ot6cpQRRXNNWhZsfqsE8bNO6NR0i1Ulgcy_qsKPd72xmC0407ip0OS4Nbum9FS_bQWu3EMl_6106mQb9WSalICMw",
  "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Brave/1 Mobile/15E148 Safari/E7FBAF",
  "enable_gofile_upload": true,
  "enable_supabase": true,
  "turboviplay_api_key": "xizpCCPcnb",
  "voesx_api_key": "AF1YD2ExCqHrCSUjnwp9fPabywVRU1hwgFf8aKBzkx9gAV7S9Phdn9cS7ZAuWcN1",
  "streamtape_login": "ad687ba4675c26af3bd4",
  "streamtape_api_key": "WgMD3kVBWMsb66q",
  "supabase_url": "https://xhfbhgklqylmfmfjtgkq.supabase.co",
  "supabase_api_key": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InhoZmJoZ2tscXlsbWZtZmp0Z2txIiwicm9sZSI6ImFub24iLCJpYXQiOjE3Nzc3NDIxNTYsImV4cCI6MjA5MzMxODE1Nn0.xIPocBS1e1QhGm080ISgU63vHXLywIH-isk0757Z3Xw",
  "discord_webhook_url": "https://discord.com/api/webhooks/1497660499670863966/GjrVGaCdXCBgvYrnSM-pIepKHqSIA_HgyiIb7NM8rn4i9L5xHM7QlZJpqD60PAtQiOa-",
  "imgbb_api_key": "9f48b991edad5d980312c5f187c7ba7f",
  "finalize_mode": "remux",
  "ffmpeg_container": "mp4",
  "ffmpeg_encoder": "libx264",
  "ffmpeg_quality": 18,
  "ffmpeg_preset": "slow",
  "disk_warning_percent": 96,
  "disk_critical_percent": 98
}
```

Click **Add secret**

---

### 3. Add CHANNELS_JSON Secret

**Name:** `CHANNELS_JSON`

**Value:** Copy the entire content from `conf/channels.json` (it's very long, make sure you get ALL of it)

Click **Add secret**

---

### 4. Add Individual Secrets (Optional but Recommended)

These are already in SETTINGS_JSON, but adding them separately allows the workflow to use them as environment variables:

#### SUPABASE_URL
**Value:** `https://xhfbhgklqylmfmfjtgkq.supabase.co`

#### SUPABASE_API_KEY
**Value:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InhoZmJoZ2tscXlsbWZtZmp0Z2txIiwicm9sZSI6ImFub24iLCJpYXQiOjE3Nzc3NDIxNTYsImV4cCI6MjA5MzMxODE1Nn0.xIPocBS1e1QhGm080ISgU63vHXLywIH-isk0757Z3Xw`

#### DISCORD_WEBHOOK
**Value:** `https://discord.com/api/webhooks/1497660499670863966/GjrVGaCdXCBgvYrnSM-pIepKHqSIA_HgyiIb7NM8rn4i9L5xHM7QlZJpqD60PAtQiOa-`

#### TURBOVIPLAY_API_KEY
**Value:** `xizpCCPcnb`

#### VOESX_API_KEY
**Value:** `AF1YD2ExCqHrCSUjnwp9fPabywVRU1hwgFf8aKBzkx9gAV7S9Phdn9cS7ZAuWcN1`

#### STREAMTAPE_LOGIN
**Value:** `ad687ba4675c26af3bd4`

#### STREAMTAPE_API_KEY
**Value:** `WgMD3kVBWMsb66q`

#### IMGBB_API_KEY
**Value:** `9f48b991edad5d980312c5f187c7ba7f`

---

## Verification

After adding all secrets, you should see:
- ✅ SETTINGS_JSON
- ✅ CHANNELS_JSON
- ✅ SUPABASE_URL
- ✅ SUPABASE_API_KEY
- ✅ DISCORD_WEBHOOK
- ✅ TURBOVIPLAY_API_KEY
- ✅ VOESX_API_KEY
- ✅ STREAMTAPE_LOGIN
- ✅ STREAMTAPE_API_KEY
- ✅ IMGBB_API_KEY

---

## Next Steps

1. Go to **Actions** tab
2. Click **Continuous Recording (24/7) with Web UI**
3. Click **Run workflow**
4. Wait for it to start
5. Check the logs for your Cloudflare Tunnel URL!

---

## Security Note

⚠️ **NEVER commit settings.json or channels.json to git!**

These files contain sensitive data:
- Chaturbate cookies
- API keys
- Webhook URLs
- Supabase credentials

They are already in `.gitignore` to prevent accidental commits.

---

**Last Updated:** 2026-05-03

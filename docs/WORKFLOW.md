# GoOnDVR Upload & Cleanup Workflow

This document explains the automated workflow for recording, uploading, and cleaning up videos on your EC2 instance.

## Workflow Overview

```
┌─────────────────┐
│  Stream Online  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Start Recording│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Save to Disk   │
│  (.ts or .mp4)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Recording Done  │
│ (Stream Ends)   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│ Finalize Recording      │
│ (Optional: Remux/       │
│  Transcode to MP4)      │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Upload to GoFile.io     │
│ (Parallel with          │
│  next recording)        │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Store Link in Databases │
│ • Local JSON (GitHub)   │
│ • Supabase (if enabled) │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Delete Local File       │
│ (Free up disk space)    │
└─────────────────────────┘
```

## Features

### 1. **Automatic Upload to GoFile**
- After each recording completes, it's automatically uploaded to GoFile.io
- Upload happens in the background while new recordings continue
- Retry mechanism with exponential backoff (4 attempts)
- Large file support with streaming upload (no memory issues)

### 2. **Dual Database System**
- **Local JSON Database** (`database/` folder)
  - Organized by username and date
  - Contains metadata: filename, GoFile link, filesize, duration
  - GitHub Actions compatible for automated backups
  
- **Supabase Database** (optional)
  - Cloud PostgreSQL database
  - Stores upload records with timestamps
  - Accessible via API for web dashboards

### 3. **Automatic Cleanup**
- Local files are deleted **only after successful upload**
- If upload fails, local file is kept for retry
- Docker cache cleanup to prevent disk bloat
- Periodic cleanup via cron job

### 4. **Disk Space Management**
- Monitors disk usage every 5 minutes
- Sends notifications when disk space is low
- Automatic cleanup prevents disk full errors

## Configuration

### Enable GoFile Upload

Edit your configuration or use the web UI:

```bash
# Via command line flag
./goondvr --enable-gofile-upload

# Or in web UI: Settings → Enable GoFile Upload
```

### Configure Supabase (Optional)

```bash
./goondvr \
  --enable-gofile-upload \
  --supabase-url "https://your-project.supabase.co" \
  --supabase-api-key "your-anon-key"
```

Or set via web UI: Settings → Supabase Configuration

### Setup Automatic Cleanup

On your EC2 instance:

```bash
# Make scripts executable
chmod +x scripts/cleanup_ec2.sh
chmod +x scripts/setup_auto_cleanup.sh

# Setup automatic cleanup (runs every 6 hours)
./scripts/setup_auto_cleanup.sh

# Or run manually
./scripts/cleanup_ec2.sh
```

## Manual Cleanup

If you need to free up space immediately:

```bash
# SSH into EC2
ssh -i aws-secrets/aws-key.pem ubuntu@YOUR_EC2_IP

# Run cleanup script
cd /path/to/goondvr
./scripts/cleanup_ec2.sh
```

Or clean specific items:

```bash
# Remove completed videos
sudo rm -rf goondvr/videos/completed/*.mp4

# Clean Docker cache
docker system prune -a -f

# Remove unused Docker volumes
docker volume prune -f
```

## Database Structure

### Local JSON Database

```
database/
├── username1/
│   ├── 2026-04-29/
│   │   └── recordings.json
│   └── 2026-04-30/
│       └── recordings.json
└── username2/
    └── 2026-04-29/
        └── recordings.json
```

**recordings.json format:**
```json
{
  "date": "2026-04-29",
  "username": "streamer_name",
  "site": "chaturbate",
  "recordings": [
    {
      "id": "streamer_1776938258_1626",
      "username": "streamer_name",
      "site": "chaturbate",
      "filename": "streamer_2026-04-29_10-45-27.mp4",
      "gofile_link": "https://gofile.io/d/ABC123",
      "uploaded_at": "2026-04-29T11:20:15Z",
      "filesize_bytes": 3449374287,
      "duration_seconds": 2400.5,
      "status": "uploaded"
    }
  ],
  "summary": {
    "total_recordings": 1,
    "total_size_bytes": 3449374287,
    "last_updated": "2026-04-29T11:20:15Z"
  }
}
```

### Supabase Database

Table: `gofile_uploads`

| Column | Type | Description |
|--------|------|-------------|
| id | integer | Auto-increment primary key |
| streamer_name | text | Username of the streamer |
| gofile_link | text | GoFile download link |
| upload_date | timestamp | When the upload completed |

## Monitoring

### Check Disk Usage

```bash
# On EC2
df -h /

# Check videos directory size
du -sh goondvr/videos/

# Check Docker usage
docker system df
```

### View Cleanup Logs

```bash
# If using cron job
tail -f /var/log/goondvr-cleanup.log

# View application logs
docker logs goondvr
```

### Check Upload Status

```bash
# View recent uploads in database
cat database/username/$(date +%Y-%m-%d)/recordings.json | jq .

# Check Supabase (if configured)
# Use Supabase dashboard or API
```

## Troubleshooting

### Upload Fails

**Symptom:** Files remain in `videos/completed/` directory

**Solutions:**
1. Check internet connectivity
2. Verify GoFile.io is accessible
3. Check logs: `docker logs goondvr`
4. Files will be kept for manual upload or retry

### Disk Full

**Symptom:** Recording stops, "no space left on device" error

**Solutions:**
```bash
# Immediate cleanup
./scripts/cleanup_ec2.sh

# Check what's using space
du -sh goondvr/* | sort -h

# Emergency: delete all completed videos
sudo rm -rf goondvr/videos/completed/*
docker system prune -a -f
```

### Database Not Created

**Symptom:** No `database/` folder or empty folders

**Cause:** Database folders are only created **after successful upload**

**Solution:** 
- Ensure `--enable-gofile-upload` is enabled
- Check that uploads are succeeding
- Failed uploads won't create database entries (by design)

### Supabase Connection Fails

**Symptom:** "failed to store upload record in Supabase" in logs

**Solutions:**
1. Verify Supabase URL and API key
2. Check table exists: `gofile_uploads`
3. Run migration: `supabase/migrations/001_create_gofile_uploads_table.sql`
4. Check Supabase dashboard for errors

## Best Practices

1. **Enable GoFile Upload** - Prevents disk full issues
2. **Setup Automatic Cleanup** - Runs every 6 hours via cron
3. **Monitor Disk Usage** - Check dashboard regularly
4. **Backup Database Folder** - Contains all upload links
5. **Use Supabase** - For centralized access to upload records
6. **Set Disk Alerts** - Configure in Settings (default: 80% warning, 90% critical)

## Performance Notes

- **Upload Speed:** Depends on EC2 bandwidth and file size
- **Parallel Processing:** Upload happens while next recording continues
- **No Downtime:** Recordings never stop due to upload/cleanup
- **Memory Efficient:** Streaming upload, no full file in memory
- **Retry Logic:** 4 attempts with exponential backoff (5s, 10s, 20s)

## Security

- GoFile links are public but unguessable (random IDs)
- Supabase uses API keys (keep them secret)
- Local database is in `database/` folder (backup regularly)
- AWS key file should have restricted permissions (600)

## Cost Considerations

- **GoFile.io:** Free tier available, check limits
- **Supabase:** Free tier: 500MB database, 1GB file storage
- **EC2:** Data transfer costs apply for uploads
- **Storage:** Local storage freed after upload (minimal cost)

## Future Enhancements

Potential improvements:
- [ ] Upload to multiple cloud providers (S3, Backblaze, etc.)
- [ ] Automatic video thumbnails
- [ ] Web dashboard to browse uploads
- [ ] Download from GoFile back to local
- [ ] Automatic expiry management
- [ ] Upload queue with priority
- [ ] Bandwidth throttling for uploads

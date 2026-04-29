# 📤 Upload Feature Documentation

## Overview
GoondVR now automatically uploads completed recordings to multiple video hosting services (GoFile.io, TurboViPlay.com, VOE.sx, Streamtape.com) and stores the links in Supabase database.

---

## ✨ Features

### 1. **Multi-Host Upload on Recording Completion**
When a recording finishes (stream ends, file size/duration limit reached), the system automatically:
1. ✅ Converts/remuxes the video (if configured)
2. ✅ Uploads to multiple hosts in parallel:
   - **GoFile.io** - Always enabled (no API key required)
   - **TurboViPlay.com** - Optional (requires API key)
   - **VOE.sx** - Optional (requires API key)
   - **Streamtape.com** - Optional (requires login + API key)
3. ✅ Stores all successful links in Supabase database in a single record
4. ✅ Deletes the local file to save disk space

### 2. **Manual Upload Button**
For files that are already in the completed folder:
- Click the **"Videos"** button in the header
- Click **"Upload Completed"** button
- Confirms and uploads all files in the completed directory
- Uploads to all configured hosts in parallel
- Stores all successful links in Supabase
- Deletes local files after successful upload

### 3. **View Uploaded Videos**
- Click **"Videos"** button to see all uploaded videos
- Videos are grouped by streamer name
- Multiple links shown per video (one for each host)
- Click **"Open Video"** to view on the hosting service
- Click copy icon to copy the link
- Search by streamer name

---

## 🔧 Configuration

### Required Settings (in settings.json):

```json
{
  "enable_gofile_upload": true,
  "turboviplay_api_key": "your-turboviplay-api-key",
  "voesx_api_key": "your-voesx-api-key",
  "streamtape_login": "your-streamtape-login",
  "streamtape_api_key": "your-streamtape-api-key",
  "supabase_url": "https://your-project.supabase.co",
  "supabase_api_key": "your-api-key"
}
}
```

**Note:** 
- GoFile.io is always enabled (no API key required)
- TurboViPlay, VOE.sx, and Streamtape are optional - leave API keys blank to skip those hosts
- Streamtape requires both login and API key (found in Account Settings → API Details)
- At least one upload host will always work (GoFile)

### How to Enable:

1. **Via Settings UI:**
   - Click "Settings" button
   - Check "Enable GoFile Upload"
   - Enter TurboViPlay API Key (optional)
   - Enter VOE.sx API Key (optional)
   - Enter Streamtape Login (optional)
   - Enter Streamtape API Key (optional)
   - Enter Supabase URL and API Key
   - Click "Apply"

2. **Via settings.json:**
   - Edit `settings.json` file
   - Set `enable_gofile_upload: true`
   - Add your API keys (optional for TurboViPlay and Videasy)
   - Add your Supabase credentials
   - Restart the application

3. **Via Command Line:**
   ```bash
   ./goondvr --enable-gofile-upload \
             --turboviplay-api-key="your-key" \
             --videasy-api-key="your-key" \
             --supabase-url="https://your-project.supabase.co" \
             --supabase-api-key="your-key"
   ```

---

## 📊 Upload Flow

### Automatic Upload (When Recording Ends):

```
Recording Ends
    ↓
Close File
    ↓
FFmpeg Conversion (if enabled)
    ↓
Upload to Multiple Hosts (Parallel):
    ├─→ GoFile.io
    ├─→ TurboViPlay.com (if API key configured)
    └─→ Videasy.net (if API key configured)
    ↓
Store All Successful Links in Supabase
    ↓
Delete Local File
    ↓
Update Disk Usage
```

### Manual Upload (For Existing Files):

```
Click "Upload Completed" Button
    ↓
Scan completed/ directory
    ↓
For each video file:
    ↓
Upload to Multiple Hosts (Parallel):
    ├─→ GoFile.io
    ├─→ TurboViPlay.com (if API key configured)
    └─→ Videasy.net (if API key configured)
    ↓
Store All Successful Links in Supabase
    ↓
Delete Local File
    ↓
Show Results (X/3 hosts successful)
```

---

## 🎯 When Upload Happens

### Automatic Upload Triggers:
1. **Stream Ends Naturally** - When streamer goes offline
2. **File Size Limit** - When max filesize is reached
3. **Duration Limit** - When max duration is reached
4. **Manual Pause** - When you pause a channel (NEW!)

### Manual Upload:
- Click "Upload Completed" button in Videos dialog
- Uploads all files in `videos/completed/` directory

---

## 📁 File Locations

### Before Upload:
```
videos/
  ├── username_2026-04-29_18-40-26.mp4  (recording)
  └── completed/
      └── username_2026-04-29_18-30-00.mp4  (finished)
```

### After Upload:
```
videos/
  └── username_2026-04-29_18-40-26.mp4  (still recording)

Supabase Database:
  ✓ username_2026-04-29_18-30-00.mp4 → GoFile: https://gofile.io/d/xxxxx
  ✓ username_2026-04-29_18-30-00.mp4 → TurboViPlay: https://turboviplay.com/embed-xxxxx.html
  ✓ username_2026-04-29_18-30-00.mp4 → Videasy: https://videasy.net/xxxxx
```

---

## 🔍 Viewing Uploaded Videos

### In the UI:
1. Click **"Videos"** button in header
2. See all uploaded videos grouped by streamer
3. Click **"Open Video"** to view
4. Click copy icon to copy link
5. Use search box to filter by streamer name

### Via API:
```bash
# Get all videos
curl http://localhost:8080/api/videos

# Get videos by username
curl http://localhost:8080/api/videos/username

# Get database stats
curl http://localhost:8080/api/database/stats
```

---

## 🛠️ Troubleshooting

### Upload Not Working?

**Check 1: Is GoFile upload enabled?**
```json
"enable_gofile_upload": true
```

**Check 2: Is Supabase configured?**
```json
"supabase_url": "https://your-project.supabase.co",
"supabase_api_key": "your-api-key"
```

**Check 3: Check the logs**
Look for messages like:
- `uploading to GoFile...`
- `upload successful: https://...`
- `upload record stored in Supabase`
- `local file deleted`

**Check 4: Files in completed folder?**
If files are in `videos/completed/`, use the manual upload button.

### Upload Failed?

If upload fails:
- ❌ File is NOT deleted
- ❌ Link is NOT stored in database
- ✅ File remains in completed folder
- ✅ You can retry with manual upload button

### Can't See Videos in UI?

**Check 1: Supabase configured?**
The Videos dialog requires Supabase to be configured.

**Check 2: Any videos uploaded?**
If no videos have been uploaded yet, the list will be empty.

**Check 3: Check browser console**
Open browser DevTools (F12) and check for errors.

---

## 📊 Database Structure

### Supabase Table: `gofile_uploads`

```sql
CREATE TABLE gofile_uploads (
    id SERIAL PRIMARY KEY,
    streamer_name TEXT NOT NULL,
    gofile_link TEXT NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Example Record:
```json
{
  "id": 1,
  "streamer_name": "lilkimchii",
  "gofile_link": "https://gofile.io/d/abc123",
  "upload_date": "2026-04-29T18:41:15.000Z"
}
```

---

## 🎨 UI Features

### Videos Dialog:
- **Glassmorphic design** - Matches the modern UI
- **Grouped by streamer** - Easy to find videos
- **Search functionality** - Filter by streamer name
- **Quick actions** - Open video or copy link
- **Upload button** - Upload completed files
- **Video count** - Shows total number of videos

### Video Cards:
- **Upload date/time** - When the video was uploaded
- **Open Video button** - Opens GoFile link in new tab
- **Copy button** - Copies link to clipboard
- **Hover effects** - Beautiful animations

---

## 🚀 Best Practices

### 1. **Monitor Disk Space**
- Enable GoFile upload to automatically free up space
- Set appropriate file size/duration limits
- Check the stats bar for disk usage

### 2. **Regular Uploads**
- Let recordings complete naturally for automatic upload
- Use manual upload button for existing files
- Check Videos dialog to verify uploads

### 3. **Backup Links**
- Links are stored in Supabase (cloud database)
- Export links periodically for backup
- Use API to get all links programmatically

### 4. **Network Considerations**
- Uploads happen in background (doesn't block recording)
- Large files may take time to upload
- Failed uploads keep local file for retry

---

## 📈 Performance

### Upload Speed:
- Depends on your internet upload speed
- Typical: 5-10 MB/s
- Large files (1GB+) may take several minutes

### Background Processing:
- ✅ Uploads happen in background
- ✅ Recording continues during upload
- ✅ Multiple uploads can happen simultaneously
- ✅ No impact on recording quality

### Disk Space Savings:
- Files are deleted after successful upload
- Saves 100% of recording size
- Only keeps currently recording files

---

## 🔐 Security

### GoFile:
- Public file hosting service
- Links are shareable
- No authentication required to view

### Supabase:
- Secure cloud database
- API key required for access
- Row Level Security enabled

### Recommendations:
- Keep API keys secure
- Don't share API keys publicly
- Use environment variables for production

---

## 📝 API Endpoints

### Get All Videos:
```
GET /api/videos
```

### Get Videos by Username:
```
GET /api/videos/:username
```

### Get Database Stats:
```
GET /api/database/stats
```

### Upload Completed Files:
```
POST /api/upload/completed
```

---

## 🎯 Future Enhancements

Potential improvements:
- [ ] Upload progress indicator
- [ ] Retry failed uploads automatically
- [ ] Custom upload destinations (S3, etc.)
- [ ] Video thumbnails in UI
- [ ] Download videos from UI
- [ ] Bulk delete from database
- [ ] Export links to CSV
- [ ] Upload queue management

---

**Last Updated:** 2026-04-29  
**Version:** 2.1.0

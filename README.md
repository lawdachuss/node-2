# GoOnDVR - Automated Recording System

## ✅ Status: FULLY OPERATIONAL

**EC2:** 54.210.37.19  
**Web UI:** http://54.210.37.19:8080  
**Disk:** 4.1 GB / 28 GB (15%)

---

## 🎨 NEW: Next-Level UI!

The web interface has been completely redesigned with:
- ✨ Modern glassmorphism design
- 🎭 Smooth animations & transitions
- 🌈 Beautiful gradient accents
- 📱 Enhanced mobile responsiveness
- 🔔 Toast notifications
- 🎯 Better visual hierarchy

**[See UI Documentation →](docs/UI_QUICK_START.md)**

---

## 🚀 Quick Start

### Connect to EC2:
```powershell
.\scripts\connect_ec2.ps1
```

### Check Status:
```powershell
.\scripts\run_workflow.ps1 status
```

### View Logs:
```powershell
.\scripts\run_workflow.ps1 logs
```

### Run Cleanup:
```powershell
.\scripts\run_workflow.ps1 cleanup
```

---

## 📋 Available Commands

| Command | Description |
|---------|-------------|
| `status` | Show system status |
| `logs` | View recent logs |
| `cleanup` | Free up disk space |
| `restart` | Restart container |
| `disk` | Check disk usage |
| `database` | View recordings |
| `channels` | List channels |
| `settings` | Show settings |

**Usage:** `.\scripts\run_workflow.ps1 [command]`

---

## 🎯 What It Does

1. ✅ Monitors streams 24/7
2. ✅ Records at max quality
3. ✅ Uploads to GoFile.io
4. ✅ Stores in databases
5. ✅ Sends Discord notifications
6. ✅ Deletes local files
7. ✅ Cleans up every 6 hours

**100% Automated - No manual work needed!**

---

## 🔗 Links

- **Dashboard:** http://54.210.37.19:8080
- **Supabase:** https://iktbuxgnnuebuoqaywev.supabase.co
- **Discord:** Check webhook channel

---

## 📁 Files

```
scripts/
├── connect_ec2.ps1      # SSH into EC2
├── run_workflow.ps1     # Run commands
├── cleanup_ec2.sh       # Cleanup script (on EC2)
└── setup_auto_cleanup.sh # Cron setup (on EC2)

settings.json            # Configuration backup
```

---

## 🔧 Configuration

**Current Settings:**
- Recording: Remux mode (no quality loss)
- Container: MP4
- Quality: Maximum available
- Upload: Automatic to GoFile
- Cleanup: Every 6 hours

**To update:** Edit `settings.json` and upload to EC2

---

## 🚨 Troubleshooting

### Container not running:
```powershell
.\scripts\run_workflow.ps1 restart
```

### Disk full:
```powershell
.\scripts\run_workflow.ps1 cleanup
```

### Check logs:
```powershell
.\scripts\run_workflow.ps1 logs
```

---

**Last Updated:** 2026-04-29  
**Mode:** FULLY AUTOMATED

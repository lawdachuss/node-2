# CI/CD Quick Start Guide

Get your automatic deployment pipeline running in 5 minutes! 🚀

## ⚡ Quick Setup (5 Minutes)

### 1. Install GitHub CLI

**Windows:**
```powershell
winget install --id GitHub.cli
```

**macOS:**
```bash
brew install gh
```

**Linux:**
```bash
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh
```

### 2. Authenticate

```bash
gh auth login
```

### 3. Run Setup Script

**Windows:**
```powershell
.\scripts\setup-github-secrets.ps1
```

**Linux/macOS:**
```bash
chmod +x scripts/setup-github-secrets.sh
./scripts/setup-github-secrets.sh
```

### 4. Push to GitHub

```bash
git add .
git commit -m "Setup CI/CD pipeline"
git push origin main
```

### 5. Watch It Deploy! 🎉

Go to: `https://github.com/YOUR_USERNAME/YOUR_REPO/actions`

## 🎯 What Happens Now?

Every time you push to `main` or `master`:

1. ✅ Code is automatically uploaded to EC2
2. ✅ Configuration is backed up
3. ✅ Docker images are rebuilt
4. ✅ Containers are restarted
5. ✅ Deployment is verified

**No manual SSH needed!** 🎊

## 🔧 Manual Deployment

Trigger deployment without pushing:

```bash
gh workflow run deploy-to-ec2.yml
```

## 📊 Monitor Deployments

**View all deployments:**
```bash
gh run list --workflow=deploy-to-ec2.yml
```

**View latest deployment logs:**
```bash
gh run view --log
```

**Watch deployment in real-time:**
```bash
gh run watch
```

## ⚠️ Troubleshooting

### SSH Connection Failed

**Check your SSH key works:**
```bash
ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
```

If it fails, you need the correct key. See [CI_CD_SETUP.md](CI_CD_SETUP.md) for solutions.

### Deployment Failed

**View logs:**
```bash
gh run view --log
```

**SSH into EC2 and check:**
```bash
ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
cd /home/ubuntu/goondvr
sudo docker-compose logs
```

### Secret Not Found

**Re-run setup script:**
```bash
.\scripts\setup-github-secrets.ps1  # Windows
./scripts/setup-github-secrets.sh   # Linux/macOS
```

## 🎮 Common Commands

```bash
# View workflow runs
gh run list

# View latest run
gh run view

# Watch current run
gh run watch

# Trigger manual deployment
gh workflow run deploy-to-ec2.yml

# View workflow file
gh workflow view deploy-to-ec2.yml

# Cancel a running workflow
gh run cancel <run-id>
```

## 📚 Full Documentation

For detailed information, see [CI_CD_SETUP.md](CI_CD_SETUP.md)

## ✅ Verification

After setup, verify everything works:

- [ ] GitHub CLI installed and authenticated
- [ ] Secret added successfully
- [ ] First push triggered deployment
- [ ] Deployment completed successfully
- [ ] Web UI accessible at http://54.210.37.19:8080
- [ ] Configuration preserved after deployment

## 🆘 Need Help?

1. Check [CI_CD_SETUP.md](CI_CD_SETUP.md) for detailed troubleshooting
2. View GitHub Actions logs
3. SSH into EC2 and check Docker logs
4. Verify security group allows ports 22 and 8080

---

**That's it!** Your CI/CD pipeline is ready. Just push code and it deploys automatically! 🚀

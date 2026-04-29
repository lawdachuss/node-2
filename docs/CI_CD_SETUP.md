# CI/CD Pipeline Setup Guide

This guide explains how to set up automatic deployment to EC2 when you push code to GitHub.

## 🎯 Overview

The CI/CD pipeline automatically:
1. ✅ Builds and tests your code on every push
2. ✅ Deploys to EC2 when you push to `main` or `master` branch
3. ✅ Backs up configuration before deployment
4. ✅ Restores configuration after deployment
5. ✅ Verifies the deployment was successful
6. ✅ Shows deployment status and logs

## 📋 Prerequisites

Before setting up CI/CD, ensure you have:

1. **GitHub Repository** - Your code is pushed to GitHub
2. **EC2 Instance** - Running and accessible
3. **SSH Key** - The private key to access your EC2 instance
4. **Docker** - Installed on your EC2 instance
5. **Docker Compose** - Installed on your EC2 instance

## 🚀 Quick Setup

### Step 1: Install GitHub CLI (if not already installed)

**Windows (PowerShell):**
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

### Step 2: Authenticate with GitHub

```bash
gh auth login
```

Follow the prompts to authenticate.

### Step 3: Add Your SSH Key as a GitHub Secret

**Option A: Using the Setup Script (Recommended)**

**Windows:**
```powershell
.\scripts\setup-github-secrets.ps1
```

**Linux/macOS:**
```bash
chmod +x scripts/setup-github-secrets.sh
./scripts/setup-github-secrets.sh
```

**Option B: Manual Setup**

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secret:
   - **Name:** `EC2_SSH_KEY`
   - **Value:** Copy the entire content of `aws-secrets/aws-key.pem`
5. Click **Add secret**

### Step 4: Update EC2 Host (if needed)

If your EC2 IP address is different from `54.210.37.19`, update it in:

`.github/workflows/deploy-to-ec2.yml`:
```yaml
env:
  EC2_HOST: YOUR_EC2_IP_HERE  # Change this
  EC2_USER: ubuntu
  APP_DIR: /home/ubuntu/goondvr
```

### Step 5: Push to GitHub

```bash
git add .
git commit -m "Setup CI/CD pipeline"
git push origin main
```

The deployment will start automatically! 🎉

## 📊 Monitoring Deployments

### View Deployment Status

1. Go to your GitHub repository
2. Click on the **Actions** tab
3. You'll see all workflow runs

### View Deployment Logs

Click on any workflow run to see detailed logs of:
- Code checkout
- SSH connection
- File upload
- Docker build
- Container startup
- Deployment verification

## 🔧 Workflows

### 1. Deploy to EC2 (`deploy-to-ec2.yml`)

**Triggers:**
- Push to `main` or `master` branch
- Manual trigger via GitHub Actions UI

**What it does:**
1. Checks out your code
2. Creates a deployment package (excludes unnecessary files)
3. Uploads to EC2 via SCP
4. Backs up existing configuration
5. Extracts new code
6. Restores configuration
7. Rebuilds Docker images
8. Restarts containers
9. Verifies deployment

### 2. Build and Test (`build-and-test.yml`)

**Triggers:**
- Pull requests to `main` or `master`
- Push to any branch except `main` or `master`

**What it does:**
1. Builds the Go application
2. Runs `go vet` for code analysis
3. Checks code formatting with `go fmt`
4. Builds Docker image
5. Validates Docker Compose configuration

## 🎮 Manual Deployment

You can manually trigger a deployment:

**Using GitHub CLI:**
```bash
gh workflow run deploy-to-ec2.yml
```

**Using GitHub UI:**
1. Go to **Actions** tab
2. Select **Deploy to EC2** workflow
3. Click **Run workflow**
4. Select branch and click **Run workflow**

## 🔍 Troubleshooting

### SSH Connection Failed

**Error:** `Permission denied (publickey)`

**Solution:**
1. Verify your SSH key is correct:
   ```bash
   ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
   ```
2. If it doesn't work locally, the key is incorrect
3. Update the `EC2_SSH_KEY` secret with the correct key

### Docker Build Failed

**Error:** `docker: command not found`

**Solution:**
Install Docker on your EC2 instance:
```bash
ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu
```

Log out and log back in for group changes to take effect.

### Deployment Verification Failed

**Error:** `Web UI returned status: 000`

**Solution:**
1. Check if containers are running:
   ```bash
   ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
   cd /home/ubuntu/goondvr
   sudo docker-compose ps
   ```
2. Check logs:
   ```bash
   sudo docker-compose logs
   ```
3. Ensure security group allows port 8080

### Configuration Lost After Deployment

**Issue:** Settings or channels are reset after deployment

**Solution:**
The workflow automatically backs up and restores:
- `conf/` directory
- `database/` directory

If this fails, manually backup before deployment:
```bash
ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
cd /home/ubuntu/goondvr
tar -czf backup-$(date +%Y%m%d-%H%M%S).tar.gz conf/ database/
```

## 🔐 Security Best Practices

1. **Never commit SSH keys** - They're in `.gitignore`
2. **Use GitHub Secrets** - Never hardcode credentials
3. **Rotate SSH keys** - Periodically update your EC2 key pair
4. **Limit SSH access** - Use security groups to restrict IP addresses
5. **Review workflow logs** - Check for any exposed sensitive data

## 📝 Customization

### Change Deployment Directory

Edit `.github/workflows/deploy-to-ec2.yml`:
```yaml
env:
  APP_DIR: /path/to/your/app  # Change this
```

### Add Pre-deployment Tests

Add steps before deployment in `deploy-to-ec2.yml`:
```yaml
- name: Run Tests
  run: |
    go test ./...
```

### Add Post-deployment Notifications

Add notification steps (e.g., Slack, Discord):
```yaml
- name: Notify Discord
  if: success()
  run: |
    curl -X POST "${{ secrets.DISCORD_WEBHOOK }}" \
      -H "Content-Type: application/json" \
      -d '{"content":"✅ Deployment successful!"}'
```

### Exclude Additional Files

Edit the `tar` command in `deploy-to-ec2.yml`:
```yaml
tar -czf deploy.tar.gz \
  --exclude='.git' \
  --exclude='videos/*' \
  --exclude='your-folder/*' \
  .
```

## 🎯 Best Practices

1. **Test locally first** - Always test changes locally before pushing
2. **Use feature branches** - Don't push directly to main
3. **Review PR checks** - Ensure build and test pass before merging
4. **Monitor deployments** - Check GitHub Actions after each push
5. **Keep backups** - The workflow backs up config, but keep manual backups too
6. **Use semantic versioning** - Tag releases for easy rollback

## 🔄 Rollback

If a deployment fails, you can rollback:

### Option 1: Revert Git Commit
```bash
git revert HEAD
git push origin main
```

### Option 2: Manual Rollback on EC2
```bash
ssh -i aws-secrets/aws-key.pem ubuntu@54.210.37.19
cd /home/ubuntu/goondvr

# Restore from backup
if [ -d "backup" ]; then
  cp -r backup/conf/* conf/
  cp -r backup/database/* database/
  sudo docker-compose restart
fi
```

### Option 3: Deploy Previous Version
```bash
# Find the commit hash of the working version
git log --oneline

# Create a new branch from that commit
git checkout -b rollback <commit-hash>

# Push to main
git checkout main
git reset --hard rollback
git push origin main --force
```

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Documentation](https://docs.docker.com/)
- [AWS EC2 Documentation](https://docs.aws.amazon.com/ec2/)
- [GitHub CLI Documentation](https://cli.github.com/manual/)

## 🆘 Getting Help

If you encounter issues:

1. **Check workflow logs** - GitHub Actions tab
2. **Check EC2 logs** - SSH into EC2 and run `docker-compose logs`
3. **Verify SSH access** - Test SSH connection manually
4. **Check security groups** - Ensure ports 22 and 8080 are open
5. **Review this guide** - Follow troubleshooting steps

## ✅ Verification Checklist

After setup, verify:

- [ ] GitHub CLI is installed and authenticated
- [ ] `EC2_SSH_KEY` secret is added to GitHub
- [ ] SSH connection to EC2 works
- [ ] Docker and Docker Compose are installed on EC2
- [ ] EC2 security group allows ports 22 and 8080
- [ ] Workflow files are committed and pushed
- [ ] First deployment completed successfully
- [ ] Web UI is accessible at `http://YOUR_EC2_IP:8080`
- [ ] Configuration was preserved after deployment

---

**Last Updated:** 2026-04-29  
**Status:** Ready for Production

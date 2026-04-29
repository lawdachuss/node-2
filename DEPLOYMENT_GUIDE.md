# EC2 Deployment Guide

## Current Issue: SSH Key Authentication Failed

The SSH key at `aws-secrets/aws-key.pem` doesn't have access to the EC2 instance at `54.210.37.19`.

## Solutions

### Option 1: Use the Correct SSH Key (Recommended)

If you have the correct key file somewhere else:

1. **Locate your correct EC2 key pair file** (the one you used when launching the instance)
2. **Copy it to the project:**
   ```powershell
   Copy-Item "C:\path\to\your\correct-key.pem" -Destination "aws-secrets\aws-key.pem"
   ```
3. **Test the connection:**
   ```powershell
   .\scripts\test_ssh_connection.ps1
   ```

### Option 2: Add Your Current Key to EC2

If you want to use the current key file, you need to add it to the EC2 instance:

#### Using AWS Systems Manager Session Manager (No SSH needed):

1. **Go to AWS Console** → EC2 → Instances
2. **Select your instance** (54.210.37.19)
3. **Click "Connect"** → Choose "Session Manager" tab
4. **Click "Connect"** (opens a browser-based terminal)
5. **Run these commands:**
   ```bash
   # Switch to ubuntu user
   sudo su - ubuntu
   
   # Create .ssh directory if it doesn't exist
   mkdir -p ~/.ssh
   chmod 700 ~/.ssh
   
   # Edit authorized_keys
   nano ~/.ssh/authorized_keys
   ```
6. **Add your public key:**
   - You need to generate a public key from your private key
   - On your local machine, run:
     ```powershell
     ssh-keygen -y -f aws-secrets/aws-key.pem
     ```
   - Copy the output (starts with `ssh-rsa`)
   - Paste it into the `authorized_keys` file in the Session Manager terminal
   - Press `Ctrl+X`, then `Y`, then `Enter` to save
7. **Set correct permissions:**
   ```bash
   chmod 600 ~/.ssh/authorized_keys
   ```
8. **Test connection** from your local machine:
   ```powershell
   .\scripts\test_ssh_connection.ps1
   ```

### Option 3: Create New Key Pair and Launch New Instance

If you don't have access to the correct key:

1. **Go to AWS Console** → EC2 → Key Pairs
2. **Create new key pair:**
   - Click "Create key pair"
   - Name: `goondvr-key`
   - Type: RSA
   - Format: .pem
   - Click "Create key pair" (downloads the file)
3. **Save the key:**
   ```powershell
   Move-Item "$env:USERPROFILE\Downloads\goondvr-key.pem" -Destination "aws-secrets\aws-key.pem"
   ```
4. **Launch new EC2 instance** with this key pair
5. **Update the IP address** in scripts:
   ```powershell
   # Edit scripts and replace 54.210.37.19 with your new instance IP
   ```

### Option 4: Use Existing Session to Deploy

If you can access the EC2 instance through AWS Console (Session Manager):

1. **Connect via Session Manager** (AWS Console → EC2 → Connect → Session Manager)
2. **Run these commands directly in the browser terminal:**
   ```bash
   # Switch to ubuntu user
   sudo su - ubuntu
   cd ~
   
   # Install Docker if not installed
   sudo apt-get update
   sudo apt-get install -y docker.io docker-compose
   sudo systemctl start docker
   sudo systemctl enable docker
   sudo usermod -aG docker ubuntu
   
   # Create project directory
   mkdir -p goondvr
   cd goondvr
   ```
3. **Upload files manually:**
   - Use AWS S3 bucket as intermediary:
     ```powershell
     # On your local machine
     aws s3 cp . s3://your-bucket/goondvr/ --recursive --exclude ".git/*" --exclude "videos/*"
     ```
   - Then in EC2 Session Manager:
     ```bash
     aws s3 cp s3://your-bucket/goondvr/ . --recursive
     ```
   
   OR use SCP with a working key if you have one

4. **Deploy:**
   ```bash
   cd ~/goondvr
   sudo docker-compose build
   sudo docker-compose up -d
   ```

## Verify Which Key is Needed

To find out which key pair your EC2 instance is using:

1. **Go to AWS Console** → EC2 → Instances
2. **Select your instance** (54.210.37.19)
3. **Look at "Key pair name"** in the instance details
4. **You need the .pem file** that matches this key pair name

## After Fixing SSH Access

Once you have SSH access working:

```powershell
# Test connection
.\scripts\test_ssh_connection.ps1

# Deploy application
.\scripts\deploy_to_ec2.ps1

# Check status
.\scripts\run_workflow.ps1 status

# View logs
.\scripts\run_workflow.ps1 logs
```

## Need Help?

- **AWS Support:** Check your EC2 instance details in AWS Console
- **Key Pair Name:** Visible in EC2 instance details
- **Session Manager:** Alternative way to access EC2 without SSH
- **Security Groups:** Ensure port 22 is open for SSH (if using SSH)

## Quick Commands Reference

```powershell
# Test SSH connection
.\scripts\test_ssh_connection.ps1

# Connect to EC2
.\scripts\connect_ec2.ps1

# Deploy application
.\scripts\deploy_to_ec2.ps1

# Check status
.\scripts\run_workflow.ps1 status

# View logs
.\scripts\run_workflow.ps1 logs

# Run cleanup
.\scripts\run_workflow.ps1 cleanup
```

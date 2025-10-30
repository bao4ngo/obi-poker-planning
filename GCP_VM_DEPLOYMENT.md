# Deploy to Google Cloud VM - Complete Guide

## Overview
- **Frontend**: Next.js on VM (Port 3000)
- **Backend**: Go API on VM (Port 8080)
- **Database**: Neon PostgreSQL (External)
- **Web Server**: Nginx as reverse proxy

## Prerequisites
- ✅ Google Cloud VM running Ubuntu 24.04 LTS
- ✅ VM has external IP address
- ✅ Firewall allows HTTP (80) and HTTPS (443)
- ✅ SSH access to VM
- ✅ Neon database ready

---

## Part 1: Initial VM Setup

### Step 1: Connect to Your VM
```bash
# From your local machine
gcloud compute ssh your-vm-name --zone=your-zone

# Or use SSH directly
ssh username@your-vm-external-ip
```

### Step 2: Update System and Install Dependencies
```bash
# Update system
sudo apt update
sudo apt upgrade -y

# Install essential tools
sudo apt install -y git curl wget build-essential

# Install Node.js 22 LTS (for frontend)
# Download and install nvm:
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
# in lieu of restarting the shell
\. "$HOME/.nvm/nvm.sh"
# Download and install Node.js:
nvm install 22
# Verify the Node.js version:
node -v # Should print "v22.21.1".
# Verify npm version:
npm -v # Should print "10.9.4".

# Install Go 1.25 (for backend)
cd ~
wget -c https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
rm go1.25.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version  # Should show go1.21.13

# Install Nginx (reverse proxy and web server)
sudo apt install -y nginx

# Check Nginx status
sudo systemctl status nginx
sudo systemctl enable nginx

# Install PM2 (process manager for Node.js)
sudo npm install -g pm2

# Verify PM2 installation
pm2 --version
```

---

## Part 2: Deploy Backend (Go API)

### Step 0: Set Up GitHub SSH Access (First Time Only)

**Generate SSH key on your VM:**
```bash
# Check if you already have an SSH key
ls -la ~/.ssh/id_*.pub

# If no key exists, generate a new one
ssh-keygen -t ed25519 -C "your-email@example.com"
# Press Enter to accept default location (~/.ssh/id_ed25519)
# Press Enter twice to skip passphrase (or set one for security)

# Start SSH agent
eval "$(ssh-agent -s)"

# Add your SSH key to the agent
ssh-add ~/.ssh/id_ed25519

# Display your public key
cat ~/.ssh/id_ed25519.pub
# Copy the entire output (starts with "ssh-ed25519")
```

**Add SSH key to GitHub:**
1. Copy the SSH public key from the output above
2. Go to GitHub: https://github.com/settings/keys
3. Click **"New SSH key"**
4. Give it a title (e.g., "GCP VM Ubuntu")
5. Paste your public key
6. Click **"Add SSH key"**

**Test SSH connection:**
```bash
# Test connection to GitHub
ssh -T git@github.com
# You should see: "Hi bao4ngo! You've successfully authenticated..."
```

### Step 1: Clone Your Repository
```bash
# Create app directory
sudo mkdir -p /var/www
cd /var/www

# Clone your repo using SSH
sudo git clone git@github.com:bao4ngo/obi-poker-planning.git

# Set ownership to your user
sudo chown -R $USER:$USER obi-poker-planning

# Configure git for this repo
cd obi-poker-planning
git config user.name "Your Name"
git config user.email "your-email@example.com"

# Navigate to backend
cd back_end
```

**Alternative: Clone with HTTPS (simpler but requires token for private repos)**
```bash
# If you prefer HTTPS instead of SSH
sudo git clone https://github.com/bao4ngo/obi-poker-planning.git
```

### Step 2: Configure Environment Variables
```bash
# Create .env file
nano .env
```

Add this content:
```env
DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=npg_Qi9lKObJM5LB
DB_NAME=neondb
DB_SSLMODE=require
ALLOWED_ORIGINS=http://your-vm-ip,http://your-domain.com
PORT=8080
```

Save with `Ctrl+O`, `Enter`, then `Ctrl+X`

### Step 3: Build and Run Backend
```bash
# Download dependencies
go mod download
go mod tidy

# Build the application
go build -o poker-api .

# Make it executable
chmod +x poker-api

# Test run with environment variables (press Ctrl+C to stop after verifying it works)
DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech \
DB_PORT=5432 \
DB_USER=neondb_owner \
DB_PASSWORD=npg_Qi9lKObJM5LB \
DB_NAME=neondb \
DB_SSLMODE=require \
PORT=8080 \
ALLOWED_ORIGINS=http://localhost:3000 \
./poker-api

# You should see:
# - "Successfully connected to PostgreSQL database"
# - "Server starting on :8080"

# If you see connection errors to localhost:5432, it means
# environment variables are not set - use the command above
```

### Step 4: Create Systemd Service for Backend
```bash
# Create service file
sudo nano /etc/systemd/system/poker-api.service
```

Add this content (press `Ctrl+O` to save, `Enter`, then `Ctrl+X` to exit):
```ini
[Unit]
Description=Poker Planning API
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=simple
User=YOUR_USERNAME_HERE
WorkingDirectory=/var/www/obi-poker-planning/back_end
ExecStart=/var/www/obi-poker-planning/back_end/poker-api
Restart=always
RestartSec=5
StandardOutput=append:/var/log/poker-api.log
StandardError=append:/var/log/poker-api-error.log

# Environment variables
Environment="DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech"
Environment="DB_PORT=5432"
Environment="DB_USER=neondb_owner"
Environment="DB_PASSWORD=npg_Qi9lKObJM5LB"
Environment="DB_NAME=neondb"
Environment="DB_SSLMODE=require"
Environment="ALLOWED_ORIGINS=http://YOUR_VM_IP_HERE"
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

**IMPORTANT: Replace these values:**
- `YOUR_USERNAME_HERE` - Run `whoami` to get your username
- `YOUR_VM_IP_HERE` - Run `curl ifconfig.me` to get your VM external IP

```bash
# Create log files with proper permissions
sudo touch /var/log/poker-api.log /var/log/poker-api-error.log
sudo chown $USER:$USER /var/log/poker-api.log /var/log/poker-api-error.log

# Reload systemd to recognize new service
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable poker-api

# Start the service
sudo systemctl start poker-api

# Check status (should show "active (running)")
sudo systemctl status poker-api

# View real-time logs
sudo journalctl -u poker-api -f
# Press Ctrl+C to exit logs
```

---

## Part 3: Deploy Frontend (Next.js)

### Step 1: Navigate to Frontend Directory
```bash
cd /var/www/obi-poker-planning/front_end
```

### Step 2: Configure Environment Variables
```bash
# Create .env.production
nano .env.production
```

Add this content:
```env
NEXT_PUBLIC_API_URL=http://your-vm-ip:8080
```

Or if using domain:
```env
NEXT_PUBLIC_API_URL=https://api.your-domain.com
```

Save with `Ctrl+O`, `Enter`, then `Ctrl+X`

### Step 3: Build Frontend
```bash
# Install dependencies (this may take a few minutes)
npm install

# Build for production
npm run build

# The build output will be in .next folder

# Test run to make sure it works (press Ctrl+C to stop)
npm start

# You should see: "Ready - started server on 0.0.0.0:3000"
```

### Step 4: Set Up PM2 for Frontend
```bash
# Start frontend with PM2
pm2 start npm --name "poker-frontend" -- start

# Save PM2 process list (so it persists after reboot)
pm2 save

# Generate and run startup script for Ubuntu 24.04
pm2 startup systemd
# Copy the command it outputs and run it (starts with 'sudo env...')

# Check PM2 status
pm2 status

# View logs
pm2 logs poker-frontend

# Other useful PM2 commands:
# pm2 stop poker-frontend    # Stop the app
# pm2 restart poker-frontend # Restart the app
# pm2 delete poker-frontend  # Remove from PM2
```

---

## Part 4: Configure Nginx Reverse Proxy

### Step 1: Create Nginx Configuration
```bash
# Create new site configuration
sudo nano /etc/nginx/sites-available/poker-planning
```

Add this content (press `Ctrl+O` to save, `Enter`, then `Ctrl+X` to exit):

**Option A: If using domain name**
```nginx
# Backend API
server {
    listen 80;
    server_name api.your-domain.com;

    # Increase buffer sizes for WebSocket
    client_body_buffer_size 128k;
    client_max_body_size 10M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeout settings
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket support
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket timeout
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
}

# Frontend
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Option B: If using only IP address (simpler setup)**
```nginx
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name _;

    # Frontend on root
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Backend API on /api
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket on /ws
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
}
```

### Step 2: Enable Configuration
```bash
# Remove default Nginx site (optional but recommended)
sudo rm /etc/nginx/sites-enabled/default

# Create symbolic link to enable your site
sudo ln -s /etc/nginx/sites-available/poker-planning /etc/nginx/sites-enabled/

# Test Nginx configuration for syntax errors
sudo nginx -t
# Should output: "syntax is ok" and "test is successful"

# If test passes, reload Nginx
sudo systemctl reload nginx

# Restart Nginx to apply changes
sudo systemctl restart nginx

# Check Nginx status
sudo systemctl status nginx

# Ensure Nginx starts on boot
sudo systemctl enable nginx
```

---

## Part 5: Configure Google Cloud Firewall

### Option A: Using gcloud CLI (from your local Windows machine)

**Note**: Make sure you have Google Cloud SDK installed on Windows. Download from: https://cloud.google.com/sdk/docs/install

```powershell
# Open PowerShell and authenticate
gcloud auth login

# Set your project
gcloud config set project YOUR_PROJECT_ID

# Allow HTTP traffic
gcloud compute firewall-rules create allow-http ^
    --allow tcp:80 ^
    --source-ranges 0.0.0.0/0 ^
    --description "Allow HTTP traffic" ^
    --direction INGRESS

# Allow HTTPS traffic
gcloud compute firewall-rules create allow-https ^
    --allow tcp:443 ^
    --source-ranges 0.0.0.0/0 ^
    --description "Allow HTTPS traffic" ^
    --direction INGRESS

# Optional: Allow direct access to app ports (for testing)
gcloud compute firewall-rules create allow-app-ports ^
    --allow tcp:3000,tcp:8080 ^
    --source-ranges 0.0.0.0/0 ^
    --description "Allow application ports"

# List firewall rules to verify
gcloud compute firewall-rules list
```

### Option B: Using GCP Console (Web UI - Easiest)
1. Go to **Google Cloud Console**: https://console.cloud.google.com
2. Navigate to **VPC Network** → **Firewall** (or search "Firewall" in the search bar)
3. Click **CREATE FIREWALL RULE**
4. Configure the rule:
   - **Name**: `allow-http-https`
   - **Description**: `Allow HTTP and HTTPS traffic`
   - **Logs**: Off (or On if you want to monitor)
   - **Network**: `default` (or your network)
   - **Priority**: `1000`
   - **Direction of traffic**: `Ingress`
   - **Action on match**: `Allow`
   - **Targets**: `All instances in the network`
   - **Source filter**: `IPv4 ranges`
   - **Source IPv4 ranges**: `0.0.0.0/0`
   - **Protocols and ports**: 
     - Check `Specified protocols and ports`
     - Select `TCP` and enter: `80,443`
5. Click **CREATE**

**Verify firewall rules:**
- Go to **VPC Network** → **Firewall**
- You should see your rule listed with:
  - Ingress ✅
  - Allow ✅
  - tcp:80,443 ✅
  - 0.0.0.0/0 ✅

---

## Part 6: Testing and Verification

### Step 1: Check Services Status
```bash
# Check backend
sudo systemctl status poker-api
curl http://localhost:8080/api/sessions

# Check frontend
pm2 status
curl http://localhost:3000

# Check Nginx
sudo systemctl status nginx
```

### Step 2: Test from Your Browser
```bash
# Get your VM external IP address
curl -4 ifconfig.me

# Or use GCP command
gcloud compute instances list
```

Access your application:
- **Frontend**: `http://YOUR_VM_IP` (e.g., http://34.123.45.67)
- **Backend API Test**: `http://YOUR_VM_IP/api/sessions` or `http://YOUR_VM_IP:8080/api/sessions`

**Expected results:**
- Frontend: Should see your Poker Planning app
- Backend: Should return `[]` (empty array) or existing sessions JSON

### Step 3: Check Logs
```bash
# Backend logs (systemd service)
sudo journalctl -u poker-api -f          # Follow logs in real-time
sudo journalctl -u poker-api -n 50       # Last 50 lines
sudo tail -f /var/log/poker-api.log      # Application stdout
sudo tail -f /var/log/poker-api-error.log # Application stderr

# Frontend logs (PM2)
pm2 logs poker-frontend                  # Follow logs
pm2 logs poker-frontend --lines 100      # Last 100 lines

# Nginx logs
sudo tail -f /var/log/nginx/error.log    # Nginx errors
sudo tail -f /var/log/nginx/access.log   # HTTP requests

# System logs
dmesg | tail                             # Kernel messages
sudo journalctl -xe                      # Recent system logs
```

---

## Part 7: Enable HTTPS with Let's Encrypt (Recommended)

### Prerequisites
- Domain name pointing to your VM IP
- Ports 80 and 443 open in firewall

### Install Certbot for Ubuntu 24.04
```bash
# Install Certbot and Nginx plugin
sudo apt install -y certbot python3-certbot-nginx

# Verify installation
certbot --version
```

### Get SSL Certificate

**For domain with both www and non-www:**
```bash
# Replace with your actual domain
sudo certbot --nginx -d your-domain.com -d www.your-domain.com -d api.your-domain.com

# Follow the prompts:
# 1. Enter your email address
# 2. Agree to terms of service (Y)
# 3. Share email with EFF (Y/N - your choice)
# 4. Choose redirect option (2 = redirect HTTP to HTTPS - recommended)
```

**For single domain:**
```bash
sudo certbot --nginx -d your-domain.com
```

### Verify SSL Certificate
```bash
# Check certificate status
sudo certbot certificates

# Test renewal (dry run)
sudo certbot renew --dry-run
```

### Auto-Renewal Setup
Certbot automatically sets up renewal. Verify with:
```bash
# Check systemd timer
sudo systemctl status certbot.timer

# Check renewal configuration
sudo cat /etc/cron.d/certbot
```

### Manual Renewal (if needed)
```bash
# Renew all certificates
sudo certbot renew

# Renew specific domain
sudo certbot renew --cert-name your-domain.com

# Reload Nginx after renewal
sudo systemctl reload nginx
```

### Test HTTPS
After getting certificate, visit:
- `https://your-domain.com` (should show secure lock icon)
- `https://api.your-domain.com/api/sessions`
- `http://your-domain.com` (should redirect to https)

---

## Part 8: Maintenance and Updates

### Update Application Code
```bash
# SSH into your VM
gcloud compute ssh your-vm-name --zone=your-zone

# Navigate to repository
cd /var/www/obi-poker-planning

# Pull latest changes from GitHub
git pull origin main

# Update Backend
cd back_end
go mod tidy
go build -o poker-api .
sudo systemctl restart poker-api
sudo systemctl status poker-api

# Update Frontend
cd ../front_end
npm install                    # Only if package.json changed
npm run build
pm2 restart poker-frontend
pm2 status

# Verify everything is running
sudo systemctl status poker-api
pm2 status
sudo systemctl status nginx
```

### Monitor Services
```bash
# Check all services status
sudo systemctl status poker-api
pm2 status
sudo systemctl status nginx

# Real-time monitoring
# Backend logs
sudo journalctl -u poker-api -f

# Frontend logs
pm2 logs poker-frontend

# System resources
top           # Press 'q' to quit
htop          # Better top (install: sudo apt install htop)
free -h       # Memory usage
df -h         # Disk usage
netstat -tulpn | grep LISTEN  # Check listening ports
```

### Restart Services
```bash
# Restart backend
sudo systemctl restart poker-api
sudo systemctl status poker-api

# Restart frontend
pm2 restart poker-frontend
pm2 status

# Restart Nginx
sudo systemctl restart nginx
sudo systemctl status nginx

# Restart all services at once
sudo systemctl restart poker-api && pm2 restart poker-frontend && sudo systemctl restart nginx

# Reboot entire VM (if necessary)
sudo reboot
# Wait a few minutes, then reconnect via SSH
```

---

## Troubleshooting

### Backend Won't Start
```bash
# Check detailed logs
sudo journalctl -u poker-api -n 100 --no-pager
sudo cat /var/log/poker-api-error.log

# Common error: "dial tcp 127.0.0.1:5432: connect: connection refused"
# This means it's trying to connect to localhost instead of Neon
# Solution: Check environment variables in service file

# Verify environment variables are set correctly
sudo systemctl show poker-api --property=Environment
# Should show DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech

# If variables are missing or wrong, edit service file:
sudo nano /etc/systemd/system/poker-api.service
# Make sure all Environment= lines are correct (see Step 4 above)

# After editing service file:
sudo systemctl daemon-reload
sudo systemctl restart poker-api
sudo systemctl status poker-api

# Check if port 8080 is already in use
sudo lsof -i :8080
sudo netstat -tlnp | grep 8080
# If something is using it, kill with: sudo kill -9 <PID>

# Test database connection manually with environment variables
cd /var/www/obi-poker-planning/back_end
DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech \
DB_PORT=5432 \
DB_USER=neondb_owner \
DB_PASSWORD=npg_Qi9lKObJM5LB \
DB_NAME=neondb \
DB_SSLMODE=require \
./poker-api
# Press Ctrl+C after checking output

# Rebuild binary if needed
cd /var/www/obi-poker-planning/back_end
go build -o poker-api .
chmod +x poker-api
sudo systemctl restart poker-api
```

### Frontend Won't Start
```bash
# Check PM2 logs in detail
pm2 logs poker-frontend --lines 200

# Check if port 3000 is in use
sudo lsof -i :3000
sudo netstat -tlnp | grep 3000

# Delete and recreate PM2 process
pm2 delete poker-frontend
cd /var/www/obi-poker-planning/front_end
npm run build
pm2 start npm --name "poker-frontend" -- start
pm2 save

# Check Node.js version
node --version  # Should be v20.x.x

# Clear Next.js cache and rebuild
rm -rf .next
npm run build
pm2 restart poker-frontend
```

### CORS Errors
```bash
# Get your VM's external IP
curl -4 ifconfig.me

# Update ALLOWED_ORIGINS in systemd service
sudo nano /etc/systemd/system/poker-api.service

# Change this line to include your IP/domain:
# Environment="ALLOWED_ORIGINS=http://YOUR_VM_IP,https://your-domain.com"

# After saving, reload and restart
sudo systemctl daemon-reload
sudo systemctl restart poker-api

# Verify the environment variable is set
sudo systemctl show poker-api --property=Environment | grep ALLOWED_ORIGINS

# Check logs for CORS-related errors
sudo journalctl -u poker-api -n 50 | grep -i cors
```

### Can't Access from Browser
```bash
# 1. Check if services are running
sudo systemctl status poker-api
pm2 status
sudo systemctl status nginx

# 2. Test locally on VM
curl http://localhost:8080/api/sessions  # Backend
curl http://localhost:3000               # Frontend
curl http://localhost                    # Nginx

# 3. Check Ubuntu firewall (UFW)
sudo ufw status
# If active, allow ports:
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp  # SSH
sudo ufw reload

# 4. Check GCP firewall rules
# Go to GCP Console → VPC Network → Firewall
# Ensure you have rules allowing tcp:80,443 from 0.0.0.0/0

# 5. Check listening ports
sudo netstat -tlnp | grep -E ':(80|443|3000|8080)'

# 6. Test from your local machine (Windows PowerShell)
# curl http://YOUR_VM_IP
# curl http://YOUR_VM_IP/api/sessions

# 7. Check Nginx configuration
sudo nginx -t
sudo cat /etc/nginx/sites-enabled/poker-planning
```

### Git/GitHub SSH Issues
```bash
# Permission denied (publickey)
# 1. Check if SSH key exists
ls -la ~/.ssh/id_*.pub

# 2. Check if SSH agent has the key
ssh-add -l

# 3. If key not loaded, add it
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519

# 4. Test GitHub connection
ssh -T git@github.com
# Should see: "Hi bao4ngo! You've successfully authenticated"

# 5. If using sudo with git, configure sudo to use your SSH
sudo -E git pull  # -E preserves environment including SSH_AUTH_SOCK

# 6. Or change ownership and use git without sudo
cd /var/www/obi-poker-planning
sudo chown -R $USER:$USER .
git pull  # No sudo needed

# 7. Check SSH key on GitHub
# Go to https://github.com/settings/keys
# Make sure your key is added

# 8. Debug SSH connection
ssh -vT git@github.com  # Verbose output for debugging
```

---

## Architecture Summary

```
Internet
    ↓
Google Cloud VM (External IP)
    ↓
Nginx (Port 80/443)
    ├─→ Frontend (localhost:3000) - PM2
    └─→ Backend API (localhost:8080) - Systemd
            ↓
        Neon PostgreSQL (External)
```

---

## Quick Reference

### Common Commands
```bash
# Check all services
sudo systemctl status poker-api
pm2 status
sudo systemctl status nginx

# View logs
sudo journalctl -u poker-api -f  # Backend
pm2 logs poker-frontend          # Frontend
sudo tail -f /var/log/nginx/error.log  # Nginx

# Restart all
sudo systemctl restart poker-api
pm2 restart poker-frontend
sudo systemctl restart nginx

# Pull updates and redeploy
cd /var/www/obi-poker-planning
git pull
cd back_end && go build -o poker-api . && sudo systemctl restart poker-api
cd ../front_end && npm run build && pm2 restart poker-frontend
```

### Important Paths
- App Directory: `/var/www/obi-poker-planning`
- Backend Binary: `/var/www/obi-poker-planning/back_end/poker-api`
- Backend Service: `/etc/systemd/system/poker-api.service`
- Nginx Config: `/etc/nginx/sites-available/poker-planning`
- Backend Logs: `/var/log/poker-api.log`
- Frontend: Managed by PM2

---

## Next Steps
1. Set up automated backups
2. Configure monitoring (e.g., Google Cloud Monitoring)
3. Set up CI/CD pipeline
4. Configure domain name and SSL
5. Set up log rotation

Good luck with your deployment! 🚀

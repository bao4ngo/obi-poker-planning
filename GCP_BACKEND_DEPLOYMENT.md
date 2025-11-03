# Deploy Go Backend API to GCP VM

This guide shows how to deploy only the Go backend API to a Google Cloud Platform VM, while keeping:
- **Database**: Neon (PostgreSQL)
- **Frontend**: Vercel
- **Backend**: GCP VM (this deployment)

## Prerequisites

- Google Cloud SDK installed (`gcloud` command available)
- A GCP project created
- Neon database already set up
- Vercel frontend already deployed
- Go backend code in `back_end/` directory

## Step 1: Set Up GCP Project

```bash
# Login to Google Cloud
gcloud auth login

# Set your project ID (replace with your actual project ID)
export PROJECT_ID="your-project-id"
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable compute.googleapis.com
```

## Step 2: Create a VM Instance

```bash
# Create a VM instance with Ubuntu
gcloud compute instances create obi-poker-backend \
  --zone=us-central1-a \
  --machine-type=e2-micro \
  --image-family=ubuntu-2204-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=10GB \
  --tags=http-server,https-server

# Create firewall rule to allow traffic on your backend port (e.g., 8080)
gcloud compute firewall-rules create allow-backend-8080 \
  --allow=tcp:8080 \
  --target-tags=http-server \
  --description="Allow traffic on port 8080 for backend API"
```

## Step 3: Connect to Your VM

```bash
# SSH into the VM
gcloud compute ssh obi-poker-backend --zone=us-central1-a
```

## Step 4: Install Go on the VM

Once connected via SSH, run these commands:

```bash
# Update system packages
sudo apt-get update
sudo apt-get upgrade -y

# Install Go (version 1.21 or later)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

## Step 5: Set Up Your Backend Application

```bash
# Create app directory
mkdir -p ~/obi-poker-backend
cd ~/obi-poker-backend

# Install git if not already installed
sudo apt-get install -y git

# Clone your repository (or you can upload files manually)
git clone https://github.com/bao4ngo/obi-poker-planning.git
cd obi-poker-planning/back_end

# Or manually upload files (from your local machine, not in SSH):
# gcloud compute scp --recurse ./back_end obi-poker-backend:~/obi-poker-backend/ --zone=us-central1-a
```

## Step 6: Configure Environment Variables

Create a `.env` file with your Neon database connection string:

```bash
# Still in the VM SSH session
cd ~/obi-poker-planning/back_end

# Create environment file
cat > .env << 'EOF'
DATABASE_URL=postgresql://username:password@your-neon-host.neon.tech/your-database?sslmode=require
PORT=8080
FRONTEND_URL=https://your-vercel-app.vercel.app
EOF

# Or export directly (not persistent across reboots)
export DATABASE_URL="postgresql://username:password@your-neon-host.neon.tech/your-database?sslmode=require"
export PORT=8080
export FRONTEND_URL="https://your-vercel-app.vercel.app"
```

**Important**: Replace the values with your actual:
- Neon database connection string
- Vercel frontend URL

## Step 7: Build and Run the Backend

```bash
# Install dependencies
go mod download
go mod tidy

# Build the application
go build -o obi-poker-api main.go

# Test run (foreground)
./obi-poker-api
```

Test from another terminal:
```bash
curl http://localhost:8080/health
```

## Step 8: Set Up as a Systemd Service (Run in Background)

Create a systemd service to keep your app running:

```bash
# Create service file
sudo nano /etc/systemd/system/obi-poker-backend.service
```

Add this content:

```ini
[Unit]
Description=Obi Poker Planning Backend API
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/obi-poker-planning/back_end
Environment="DATABASE_URL=postgresql://username:password@your-neon-host.neon.tech/your-database?sslmode=require"
Environment="PORT=8080"
Environment="FRONTEND_URL=https://your-vercel-app.vercel.app"
ExecStart=/home/YOUR_USERNAME/obi-poker-planning/back_end/obi-poker-api
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Replace**:
- `YOUR_USERNAME` with your VM username (run `whoami` to check)
- Database URL, frontend URL with your actual values

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable obi-poker-backend
sudo systemctl start obi-poker-backend

# Check status
sudo systemctl status obi-poker-backend

# View logs
sudo journalctl -u obi-poker-backend -f
```

## Step 9: Get Your VM's External IP

```bash
# From your local machine (not in SSH)
gcloud compute instances describe obi-poker-backend \
  --zone=us-central1-a \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)'
```

Your backend will be accessible at: `http://YOUR_VM_IP:8080`

## Step 10: Update Vercel Frontend Configuration

Update your Vercel frontend environment variables to point to the new backend:

1. Go to your Vercel dashboard
2. Select your project
3. Go to Settings → Environment Variables
4. Update the API URL:
   - Name: `NEXT_PUBLIC_API_URL`
   - Value: `http://YOUR_VM_IP:8080`
5. Redeploy your frontend

## Step 11: Set Up HTTPS with Public IP (No Domain Required)

Since Let's Encrypt requires a domain name, here are your options for HTTPS with just a public IP:

### Option A: Self-Signed Certificate (Development/Testing)

**Pros:** Free, works immediately  
**Cons:** Browser will show security warning, not recommended for production

### Option B: Use a Free DNS Service (Recommended)

**Pros:** Real SSL certificate, no browser warnings  
**Cons:** Uses a subdomain like `your-app.nip.io`

### Option C: Keep HTTP on Port 8080

**Pros:** Simple, works immediately  
**Cons:** No encryption (Vercel to backend traffic is unencrypted)

---

## Option A: Self-Signed Certificate for HTTPS

### Prerequisites

Firewall rules allowing HTTPS traffic

### Prerequisites

Firewall rules allowing HTTPS traffic

### Step 11.1: Update Firewall Rules

```bash
# From your local machine, ensure HTTPS is allowed
gcloud compute firewall-rules create allow-https \
  --allow=tcp:443 \
  --target-tags=https-server \
  --description="Allow HTTPS traffic on port 443"

# Verify your VM has the https-server tag (should already have it from Step 2)
gcloud compute instances describe obi-poker-backend --zone=asia-southeast1-b \
  --format='get(tags.items[])'
```

### Step 11.2: Install Nginx

```bash
# SSH into your VM
gcloud compute ssh obi-poker-backend --zone=asia-southeast1-b

# Install Nginx
sudo apt-get update
sudo apt-get install -y nginx
```

### Step 11.3: Generate Self-Signed Certificate

```bash
# Create SSL directory
sudo mkdir -p /etc/nginx/ssl

# Generate self-signed certificate (valid for 365 days)
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/nginx-selfsigned.key \
  -out /etc/nginx/ssl/nginx-selfsigned.crt \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=$(curl -s ifconfig.me)"

# Create Diffie-Hellman group (this may take a few minutes)
sudo openssl dhparam -out /etc/nginx/ssl/dhparam.pem 2048
```

### Step 11.4: Configure Nginx with Self-Signed Certificate

```bash
# Create Nginx configuration
sudo nano /etc/nginx/sites-available/obi-poker-backend
```

Add this configuration:

```nginx
# HTTP server - redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name _;

    return 301 https://$host$request_uri;
}

# HTTPS server with self-signed certificate
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name _;

    # Self-signed SSL certificate
    ssl_certificate /etc/nginx/ssl/nginx-selfsigned.crt;
    ssl_certificate_key /etc/nginx/ssl/nginx-selfsigned.key;
    ssl_dhparam /etc/nginx/ssl/dhparam.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_session_timeout 10m;
    ssl_session_cache shared:SSL:10m;

    # Proxy settings for backend API
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        
        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        
        # Standard proxy headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $server_name;
        
        proxy_cache_bypass $http_upgrade;
        proxy_buffering off;
        proxy_read_timeout 86400;
    }
}
```

### Step 11.5: Enable and Test Nginx

```bash
# Enable the site
sudo ln -s /etc/nginx/sites-available/obi-poker-backend /etc/nginx/sites-enabled/

# Remove default site
sudo rm /etc/nginx/sites-enabled/default

# Test Nginx configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx

# Check status
sudo systemctl status nginx
```

### Step 11.6: Get Your Public IP

```bash
# Get your VM's external IP
gcloud compute instances describe obi-poker-backend \
  --zone=asia-southeast1-b \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)'

# Or from within the VM:
curl ifconfig.me
```

### Step 11.7: Update Vercel Frontend Configuration

Your backend is now accessible at: `https://YOUR_VM_IP`

Update your Vercel environment variables:
1. Go to Vercel dashboard → Your project → Settings → Environment Variables
2. Update: `NEXT_PUBLIC_API_URL` = `https://YOUR_VM_IP`
3. Redeploy your frontend

**Important:** Your browser and Vercel will show a security warning because the certificate is self-signed. You'll need to:
- Click "Advanced" → "Proceed to site" in your browser
- In production, consider using Option B below for a real certificate

### Step 11.8: Test HTTPS Connection

```bash
# Test from local machine (will show certificate warning)
curl -k https://YOUR_VM_IP/api/sessions

# Test with certificate details
openssl s_client -connect YOUR_VM_IP:443
```

---

## Option B: Use nip.io for Free DNS + Real SSL Certificate (Better Option!)

nip.io is a free wildcard DNS service that maps any IP address to a hostname.

### How it works:
- If your IP is `34.101.123.45`
- You can use: `34.101.123.45.nip.io` as your domain
- It automatically resolves to your IP

### Step B1: Install Certbot

```bash
# SSH into your VM
gcloud compute ssh obi-poker-backend --zone=asia-southeast1-b

# Install Nginx and Certbot
sudo apt-get update
sudo apt-get install -y nginx certbot python3-certbot-nginx
```

### Step B2: Configure Nginx for nip.io Domain

```bash
# Get your public IP first
export PUBLIC_IP=$(curl -s ifconfig.me)
echo "Your nip.io domain: ${PUBLIC_IP}.nip.io"

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/obi-poker-backend
```

Add this configuration (replace `YOUR_IP` with your actual IP):

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name YOUR_IP.nip.io;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_buffering off;
        proxy_read_timeout 86400;
    }
}
```

### Step B3: Enable Nginx

```bash
# Enable the site
sudo ln -s /etc/nginx/sites-available/obi-poker-backend /etc/nginx/sites-enabled/
sudo rm /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl restart nginx
```

### Step B4: Get SSL Certificate from Let's Encrypt

```bash
# Replace with your actual IP
export PUBLIC_IP=$(curl -s ifconfig.me)

# Get SSL certificate
sudo certbot --nginx -d ${PUBLIC_IP}.nip.io --non-interactive --agree-tos --email your-email@example.com
```

**Note:** Let's Encrypt may rate-limit nip.io domains. If this fails, use Option A (self-signed) or Option C (HTTP only).

### Step B5: Update Vercel

Update your Vercel environment variable:
- `NEXT_PUBLIC_API_URL` = `https://YOUR_IP.nip.io`

---

## Option C: Keep HTTP Only (Simplest)

If HTTPS complications are too much, you can keep using HTTP:

**Your backend URL:** `http://YOUR_VM_IP:8080`

**Security consideration:** 
- Traffic between Vercel and your backend is unencrypted
- For production, consider deploying your backend to a platform with built-in HTTPS (like Cloud Run, Railway, or Render)

---

## Recommended Approach

For a production app without a domain, I recommend:

1. **Short term:** Use **Option A** (self-signed certificate) with HTTPS on port 443
2. **Long term:** Get a cheap domain ($1-10/year) and use real SSL certificates
   - Namecheap, Google Domains, or Cloudflare offer cheap domains
   - Then follow the original Step 11 guide for proper SSL

Or consider deploying to platforms with built-in HTTPS:
- **Google Cloud Run** - Automatic HTTPS, serverless
- **Railway** - Free tier, automatic HTTPS
- **Render** - Free tier, automatic HTTPS

---

## Useful Commands

```bash
# View backend logs
sudo journalctl -u obi-poker-backend -f

# Restart backend
sudo systemctl restart obi-poker-backend

# Stop backend
sudo systemctl stop obi-poker-backend

# Check backend status
sudo systemctl status obi-poker-backend

# Check Nginx status (if using HTTPS)
sudo systemctl status nginx

# View Nginx error logs
sudo tail -f /var/log/nginx/error.log

# SSH into VM
gcloud compute ssh obi-poker-backend --zone=asia-southeast1-b

# Update backend code
cd ~/obi-poker-planning
git pull
cd back_end
go build -o obi-poker-api main.go
sudo systemctl restart obi-poker-backend
```

## Cost Estimation

- **e2-micro VM**: ~$7-10/month (eligible for free tier)
- **Network egress**: Variable based on usage
- **Neon database**: Based on your plan
- **Vercel**: Based on your plan

## Security Best Practices

1. **Use environment variables** for sensitive data (never commit .env)
2. **Set up firewall rules** to only allow necessary ports
3. **Enable CORS** properly in your Go backend for your Vercel domain
4. **Use HTTPS** in production (via Nginx + Certbot or Cloudflare)
5. **Regular updates**: Keep your VM and dependencies updated
6. **Monitoring**: Set up Google Cloud Monitoring for alerts

## Troubleshooting

### Backend won't start
```bash
# Check logs
sudo journalctl -u obi-poker-backend -n 50

# Check if port is in use
sudo lsof -i :8080

# Test database connection
psql "postgresql://username:password@your-neon-host.neon.tech/your-database?sslmode=require"
```

### Cannot connect from frontend
- Verify firewall rules allow port 8080
- Check CORS settings in your Go backend
- Verify the VM external IP is correct
- Check if backend is running: `sudo systemctl status obi-poker-backend`

### High latency
- Consider using a GCP region closer to your users
- Use a CDN or edge functions
- Optimize database queries

## Next Steps

1. ✅ Backend deployed on GCP VM
2. ✅ Frontend connects to GCP backend
3. ✅ Database on Neon
4. Consider: Set up monitoring and alerting
5. Consider: Implement CI/CD for automatic deployments
6. Consider: Use Cloud Load Balancer for multiple VM instances

---

**Need help?** Check the logs with `sudo journalctl -u obi-poker-backend -f`

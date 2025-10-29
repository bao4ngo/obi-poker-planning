# Quick Render.com Deployment Guide

## Prerequisites
✅ Neon database ready
✅ Vercel frontend deployed
✅ Code with pgx driver (SCRAM-SHA-256 compatible)

## Deploy Backend to Render.com (5 minutes)

### Step 1: Push to GitHub (if not already)
```powershell
git add .
git commit -m "Prepare for Render deployment"
git push origin main
```

### Step 2: Create Render Account & Deploy
1. Go to https://render.com/
2. Sign up/Login with your GitHub account
3. Click **"New +"** → **"Web Service"**
4. Connect your GitHub repository: `bao4ngo/obi-poker-planning`
5. Configure:
   - **Name**: `poker-planning-api`
   - **Region**: `Singapore`
   - **Branch**: `main`
   - **Root Directory**: `back_end`
   - **Runtime**: `Docker`
   - **Plan**: `Free`

### Step 3: Set Environment Variables
Add these in Render dashboard:

```
DB_HOST=ep-calm-voice-a1pxz353-pooler.ap-southeast-1.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=npg_Qi9lKObJM5LB
DB_NAME=neondb
DB_SSLMODE=require
ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
PORT=8080
```

### Step 4: Deploy
- Click **"Create Web Service"**
- Render will automatically:
  - Build your Docker image
  - Deploy to Singapore region
  - Give you a URL like: `https://poker-planning-api.onrender.com`

### Step 5: Update Vercel Frontend
Go to Vercel → Your project → Settings → Environment Variables:
```
NEXT_PUBLIC_API_URL=https://poker-planning-api.onrender.com
```
Redeploy frontend.

## Verify Deployment
```powershell
# Test API
curl https://poker-planning-api.onrender.com/api/sessions

# Check logs in Render dashboard
```

## Benefits vs Fly.io
✅ No CLI needed - everything via web UI
✅ No network connectivity issues
✅ Free tier available
✅ Auto-deploys on git push
✅ Easy environment variable management

## Notes
- First request may be slow (free tier spins down after inactivity)
- Upgrade to paid plan for always-on service
- Render handles SSL certificates automatically

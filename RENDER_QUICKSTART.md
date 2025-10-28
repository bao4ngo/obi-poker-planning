# Quick Deploy to Render.com

The fastest way to deploy your backend to Render.com

## 🚀 Quick Steps

### 1. Sign Up
Go to [render.com](https://render.com) and create a free account

### 2. Push to Git (if not already)
```powershell
git init
git add .
git commit -m "Initial commit"
git remote add origin YOUR_GIT_URL
git push -u origin main
```

### 3. Deploy with Blueprint
1. Click **"New +"** → **"Blueprint"**
2. Connect your Git repository
3. Select your repository
4. Render detects `render.yaml`
5. Click **"Apply"**
6. Wait 5-10 minutes for deployment

### 4. Setup Database Schema
```powershell
# Get connection string from Render dashboard
# Go to PostgreSQL service → Connect → External URL

# Run schema
psql "YOUR_CONNECTION_STRING" -f back_end/database/schema.sql
```

### 5. Update Frontend
In Vercel, set environment variable:
```
NEXT_PUBLIC_API_URL=https://poker-planning-api.onrender.com
```

### 6. Update CORS
In Render dashboard:
- Go to web service → Environment
- Update `ALLOWED_ORIGINS`:
```
https://your-app.vercel.app,http://localhost:3000
```

## ✅ Done!

Your backend will be at: `https://poker-planning-api.onrender.com`

## 📖 Need More Details?

See [RENDER_DEPLOYMENT.md](./RENDER_DEPLOYMENT.md) for complete guide.

## ⚠️ Important

Free tier services **spin down after 15 minutes** of inactivity. First request will be slow (30-60s).

**Solution**: Upgrade to Starter plan ($7/month) for always-on service.

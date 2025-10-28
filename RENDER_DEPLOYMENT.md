# Deploy to Render.com Guide

This guide will help you deploy the Poker Planning application backend and database to Render.com, with the frontend on Vercel.

## Why Render.com?

- ✅ **Easy to use** - Simple dashboard, no CLI required
- ✅ **Free tier** - PostgreSQL database + web service included
- ✅ **Auto-deploy** - Deploys from Git automatically
- ✅ **Built-in PostgreSQL** - Managed database included
- ✅ **No credit card** required for free tier

## Prerequisites

1. **Render.com Account**: Sign up at [render.com](https://render.com)
2. **Git Repository**: Your code should be in a Git repository (GitHub, GitLab, or Bitbucket)
   - If you don't have it in Git yet, initialize it:
   ```powershell
   git init
   git add .
   git commit -m "Initial commit"
   ```
3. **Vercel Account**: For frontend deployment (already done ✅)

## Deployment Options

### Option 1: Blueprint (Automated - Recommended)

Using the `render.yaml` file for automatic setup.

### Option 2: Manual Setup

Step-by-step manual configuration via dashboard.

---

## Option 1: Deploy with Blueprint (Automated)

This is the easiest method - Render will read the `render.yaml` file and set everything up automatically.

### Step 1: Push to Git

```powershell
# If not already in Git
git init
git add .
git commit -m "Add Render configuration"

# Push to your Git provider (GitHub, GitLab, etc.)
git remote add origin https://github.com/yourusername/your-repo.git
git push -u origin main
```

### Step 2: Connect Render to Your Repository

1. Go to [render.com](https://render.com) and sign in
2. Click **"New +"** → **"Blueprint"**
3. Connect your Git repository (GitHub, GitLab, or Bitbucket)
4. Select your repository
5. Give your blueprint a name: `poker-planning`
6. Render will detect the `render.yaml` file
7. Click **"Apply"**

### Step 3: Update CORS Origins

After deployment:
1. Go to your web service in Render dashboard
2. Click **"Environment"** tab
3. Find `ALLOWED_ORIGINS` and update it with your Vercel URL:
   ```
   https://your-app.vercel.app,http://localhost:3000
   ```
4. Click **"Save Changes"** (service will redeploy automatically)

### Step 4: Initialize Database Schema

1. In Render dashboard, go to your **postgres service** (`poker-planning-db`)
2. Click **"Connect"** → Copy the **External Database URL**
3. On your local machine, run:
   ```powershell
   # Install psql if you haven't (PostgreSQL client)
   # Or use pgAdmin with the connection string
   
   # Using psql:
   $env:DATABASE_URL = "postgres://user:pass@host/dbname"
   psql $env:DATABASE_URL -f back_end/database/schema.sql
   ```

   Or use Render's built-in shell:
   - In your postgres service, click **"Shell"**
   - Run:
   ```bash
   psql $DATABASE_URL
   ```
   - Then paste the contents of `back_end/database/schema.sql`

---

## Option 2: Manual Setup (Step-by-Step)

If you prefer manual setup or the blueprint doesn't work:

### Step 1: Create PostgreSQL Database

1. Sign in to [render.com](https://render.com)
2. Click **"New +"** → **"PostgreSQL"**
3. Configure:
   - **Name**: `poker-planning-db`
   - **Database**: `poker_planning`
   - **User**: `poker_planning_user` (optional)
   - **Region**: Choose closest to you (e.g., Oregon)
   - **Plan**: **Free**
4. Click **"Create Database"**
5. Wait for database to be provisioned (~2 minutes)

### Step 2: Initialize Database Schema

1. In your database dashboard, click **"Connect"** button
2. Copy the **External Connection String**
3. Use one of these methods:

   **Method A: Using psql (if installed)**
   ```powershell
   psql "postgres://username:password@host:5432/poker_planning" -f back_end/database/schema.sql
   ```

   **Method B: Using Render Shell**
   - Click **"Shell"** tab in your database
   - Run: `\c poker_planning`
   - Copy and paste contents of `back_end/database/schema.sql`

   **Method C: Using pgAdmin**
   - Open pgAdmin
   - Create new server with the external connection details
   - Run the SQL from `back_end/database/schema.sql`

### Step 3: Create Web Service (Backend)

1. Click **"New +"** → **"Web Service"**
2. Connect your Git repository
3. Configure:
   - **Name**: `poker-planning-api`
   - **Region**: Same as database (e.g., Oregon)
   - **Branch**: `main`
   - **Root Directory**: `back_end`
   - **Environment**: **Docker**
   - **Dockerfile Path**: `back_end/Dockerfile`
   - **Plan**: **Free**
4. Click **"Advanced"** to set environment variables

### Step 4: Configure Environment Variables

Add these environment variables:

| Key | Value | Notes |
|-----|-------|-------|
| `PORT` | `10000` | Render uses port 10000 |
| `DB_HOST` | (from database internal URL) | e.g., `poker-planning-db.region.render.com` |
| `DB_PORT` | `5432` | PostgreSQL default |
| `DB_NAME` | `poker_planning` | Your database name |
| `DB_USER` | (from database) | e.g., `poker_planning_user` |
| `DB_PASSWORD` | (from database) | Copy from database connection info |
| `ALLOWED_ORIGINS` | `https://your-app.vercel.app,http://localhost:3000` | Update with your Vercel URL |

**To get database credentials:**
1. Go to your PostgreSQL service
2. Click **"Connect"** 
3. Look at the **Internal Connection String**: 
   ```
   postgres://user:password@host:5432/dbname
   ```
   - Extract: user, password, host, dbname

### Step 5: Deploy

1. Click **"Create Web Service"**
2. Render will:
   - Pull your code from Git
   - Build the Docker image
   - Deploy the service
   - Give you a URL like: `https://poker-planning-api.onrender.com`

---

## Step 6: Update Frontend (Vercel)

Update your frontend to use the new backend URL:

1. Go to [vercel.com](https://vercel.com/dashboard)
2. Select your project
3. Go to **Settings** → **Environment Variables**
4. Update `NEXT_PUBLIC_API_URL`:
   ```
   https://poker-planning-api.onrender.com
   ```
5. Go to **Deployments** tab
6. Click **"Redeploy"** on the latest deployment

---

## Verify Deployment

### Test Backend Health

```powershell
# Test API endpoint
curl https://poker-planning-api.onrender.com/api/sessions
```

Should return an empty array `[]` or existing sessions.

### Test Frontend

1. Visit your Vercel URL: `https://your-app.vercel.app`
2. Create a new session
3. Join the session
4. Test voting functionality

### Test WebSocket

WebSocket should automatically work with `wss://` (secure WebSocket) since Render provides HTTPS.

---

## Important: Render Free Tier Limitations

⚠️ **Free tier services spin down after 15 minutes of inactivity**

This means:
- First request after inactivity will be slow (30-60 seconds)
- Subsequent requests will be fast
- WebSocket connections may drop during inactivity

**Solutions:**
1. **Use paid tier** ($7/month) - Keeps service always running
2. **Add uptime monitor** - Ping your service every 10 minutes to keep it awake
3. **Show loading message** - Inform users about initial load time

---

## Monitoring and Maintenance

### View Logs

1. Go to your web service in Render
2. Click **"Logs"** tab
3. Real-time logs will appear

### View Metrics

1. Click **"Metrics"** tab
2. See CPU, memory, and request metrics

### Database Management

#### Using Render Shell
1. Go to PostgreSQL service
2. Click **"Shell"** tab
3. Run SQL commands directly

#### Using External Tools (pgAdmin, DBeaver)
1. Get **External Connection String**
2. Configure your tool with these details

#### Backup Database
1. Go to PostgreSQL service
2. Click **"Backups"** tab
3. Backups are automatic on paid plans
4. For free tier, manual backup:
   ```powershell
   pg_dump "postgres://connection-string" > backup.sql
   ```

---

## Update/Redeploy

### Automatic Deployment (Default)
Render automatically deploys when you push to your Git repository:

```powershell
git add .
git commit -m "Update feature"
git push origin main
```

Render will detect the push and redeploy automatically.

### Manual Deployment
1. Go to your web service
2. Click **"Manual Deploy"** → **"Deploy latest commit"**

---

## Troubleshooting

### Service Not Starting

**Check logs:**
1. Go to web service
2. Click "Logs" tab
3. Look for errors

**Common issues:**
- Missing environment variables
- Database connection failed
- Port configuration wrong (should be 10000)

### Database Connection Failed

**Verify credentials:**
```powershell
# Test connection
psql "postgres://user:pass@host:5432/dbname" -c "SELECT 1"
```

**Use Internal URL:**
- Use the **Internal Connection String** for `DB_HOST`
- Format: `poker-planning-db-xxx.postgres.render.com`

### CORS Errors

**Update ALLOWED_ORIGINS:**
1. Go to web service
2. Environment tab
3. Update `ALLOWED_ORIGINS` with your Vercel URL
4. Include protocol: `https://your-app.vercel.app`

### WebSocket Not Connecting

**Check protocol:**
- Production should use `wss://` (secure)
- Your code already handles this in `api.ts`

**Verify backend URL:**
- Make sure `NEXT_PUBLIC_API_URL` in Vercel uses `https://`

---

## Cost Breakdown

### Render Free Tier
- ✅ **PostgreSQL**: 1GB storage
- ✅ **Web Service**: 512MB RAM, shared CPU
- ✅ **Bandwidth**: 100GB/month
- ⚠️ **Limitation**: Spins down after 15 min inactivity

### Render Paid Plans
- **Starter ($7/month)**: Always-on, no spin down
- **Standard ($25/month)**: More resources, backups

### Total Cost
- **Free**: $0/month (Render free + Vercel free)
- **Hobby**: $7/month (Render starter + Vercel free)

---

## Alternative: Keep Service Alive (Free Tier)

If you want to use free tier without spin-down, use a service like **UptimeRobot**:

1. Sign up at [uptimerobot.com](https://uptimerobot.com) (free)
2. Add monitor:
   - **Type**: HTTP(s)
   - **URL**: `https://poker-planning-api.onrender.com/api/sessions`
   - **Interval**: 5 minutes
3. This pings your service regularly, keeping it awake

⚠️ **Note**: This may violate Render's ToS for free tier. Better to upgrade to paid.

---

## Comparison: Render vs Fly.io

| Feature | Render | Fly.io |
|---------|--------|--------|
| Ease of use | ⭐⭐⭐⭐⭐ Very easy | ⭐⭐⭐ Requires CLI |
| Free tier | ✅ PostgreSQL + Web | ✅ 3 VMs + PostgreSQL |
| Spin down | ⚠️ Yes (15 min) | ❌ No |
| Dashboard | ✅ Great UI | ⭐⭐⭐ Basic |
| Git integration | ✅ Built-in | ⭐⭐ Via GitHub Actions |
| China access | ❓ May work | ❌ Blocked |

---

## Next Steps

1. ✅ Deploy backend to Render
2. ✅ Initialize database schema
3. ✅ Update Vercel environment variable
4. ✅ Test complete flow
5. ⭐ (Optional) Upgrade to paid tier for always-on service
6. ⭐ (Optional) Add custom domain

---

## Support

- **Render Docs**: [render.com/docs](https://render.com/docs)
- **Render Community**: [community.render.com](https://community.render.com)
- **Status**: [status.render.com](https://status.render.com)

---

**Congratulations! Your app is now deployed! 🎉**

Your architecture:
```
Users → Vercel (Frontend) → Render (Backend) → Render (PostgreSQL)
```

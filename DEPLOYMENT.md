# Deployment Guide

This guide will help you deploy the Poker Planning application to the cloud:
- **Backend**: Fly.io (with PostgreSQL database)
- **Frontend**: Vercel

## Prerequisites

Before you begin, make sure you have:

1. **Fly.io Account**: Sign up at [fly.io](https://fly.io)
2. **Vercel Account**: Sign up at [vercel.com](https://vercel.com)
3. **Fly CLI**: Install from [fly.io/docs/hands-on/install-flyctl/](https://fly.io/docs/hands-on/install-flyctl/)
4. **Vercel CLI** (optional): `npm install -g vercel`
5. **Git**: Ensure your project is in a Git repository

## Part 1: Deploy Backend to Fly.io

### Step 1: Login to Fly.io

```bash
fly auth login
```

### Step 2: Create PostgreSQL Database

```bash
# Create a new Postgres cluster (choose a name for your database app)
fly postgres create --name poker-planning-db --region sea

# Note down the connection details provided (username, password, hostname, port, database name)
```

The output will show something like:
```
Postgres cluster poker-planning-db created
  Username:    postgres
  Password:    <password>
  Hostname:    poker-planning-db.internal
  Proxy port:  5432
  Postgres port: 5433
  Connection string: postgres://postgres:<password>@poker-planning-db.internal:5432/poker_planning
```

### Step 3: Initialize the Database Schema

```bash
# Connect to your Postgres database
fly postgres connect -a poker-planning-db

# Once connected, create the database and run the schema
CREATE DATABASE poker_planning;
\c poker_planning

# Copy and paste the contents of back_end/database/schema.sql
# Or exit and use:
\q

# From your local machine, you can also use:
fly postgres connect -a poker-planning-db < back_end/database/schema.sql
```

Alternatively, you can use Fly.io proxy to connect from your local machine:

```bash
# Start proxy in a terminal
fly proxy 5433:5432 -a poker-planning-db

# In another terminal, connect using psql
psql postgres://postgres:<password>@localhost:5433/poker_planning -f back_end/database/schema.sql
```

### Step 4: Deploy Backend Application

```bash
cd back_end

# Launch the app (first time)
fly launch

# You'll be prompted with questions:
# - App name: Choose a unique name (e.g., obi-poker-planning)
# - Region: Choose closest to you (e.g., sea for Seattle)
# - Would you like to set up a Postgresql database now? NO (we already created one)
# - Would you like to set up an Upstash Redis database now? NO
# - Would you like to deploy now? NO (we need to configure first)
```

### Step 5: Configure Environment Variables

```bash
# Set database connection details
fly secrets set DB_HOST=obi-poker-planning-db.internal -a obi-poker-planning
fly secrets set DB_PORT=5432 -a obi-poker-planning
fly secrets set DB_USER=postgres -a obi-poker-planning
fly secrets set DB_PASSWORD=<your-password> -a obi-poker-planning
fly secrets set DB_NAME=poker_planning -a obi-poker-planning

# Set CORS to allow your frontend (we'll update this after deploying frontend)
fly secrets set ALLOWED_ORIGINS=http://localhost:3000 -a obi-poker-planning
```

### Step 6: Attach Database to App

```bash
# Attach the database to your app
fly postgres attach --app obi-poker-planning obi-poker-planning-db
```

### Step 7: Deploy the Backend

```bash
# Deploy the application
fly deploy

# Check the status
fly status

# View logs
fly logs

# Your backend will be available at: https://obi-poker-planning.fly.dev
```

## Part 2: Deploy Frontend to Vercel

### Step 1: Prepare Frontend

Make sure your frontend code is pushed to a Git repository (GitHub, GitLab, or Bitbucket).

### Step 2: Deploy to Vercel (Option A: Via Dashboard)

1. Go to [vercel.com/new](https://vercel.com/new)
2. Import your Git repository
3. Select the `front_end` folder as the root directory
4. Configure Build Settings:
   - **Framework Preset**: Next.js
   - **Root Directory**: `front_end`
   - **Build Command**: `npm run build`
   - **Output Directory**: `.next`
5. Add Environment Variable:
   - **Name**: `NEXT_PUBLIC_API_URL`
   - **Value**: `https://obi-poker-planning.fly.dev` (replace with your Fly.io URL)
6. Click "Deploy"

### Step 3: Deploy to Vercel (Option B: Via CLI)

```bash
cd front_end

# Login to Vercel
vercel login

# Set environment variable
vercel env add NEXT_PUBLIC_API_URL

# When prompted, enter: https://obi-poker-planning.fly.dev

# Deploy to production
vercel --prod
```

### Step 4: Update Backend CORS Settings

After deploying to Vercel, update the backend to allow your Vercel domain:

```bash
# Update CORS to include your Vercel URL
fly secrets set ALLOWED_ORIGINS=https://your-app.vercel.app,http://localhost:3000 -a obi-poker-planning

# Restart the app to apply changes
fly apps restart obi-poker-planning
```

## Part 3: Testing Your Deployment

1. Visit your Vercel URL: `https://your-app.vercel.app`
2. Create a new session
3. Join the session and test the WebSocket connection
4. Verify that votes are being saved and persisted

## Monitoring and Maintenance

### View Backend Logs

```bash
fly logs -a obi-poker-planning
```

### View Backend Metrics

```bash
fly dashboard -a obi-poker-planning
```

### Access Database

```bash
# Connect to PostgreSQL
fly postgres connect -a poker-planning-db

# Or use proxy for local tools
fly proxy 5433:5432 -a poker-planning-db
```

### Update Backend

```bash
cd back_end
fly deploy -a obi-poker-planning
```

### Update Frontend

Push changes to your Git repository, and Vercel will automatically deploy the updates.

## Custom Domain (Optional)

### For Backend (Fly.io)

```bash
# Add a custom domain
fly certs add api.yourdomain.com -a obi-poker-planning

# Follow the DNS instructions provided
```

### For Frontend (Vercel)

1. Go to your project settings in Vercel
2. Navigate to "Domains"
3. Add your custom domain
4. Follow the DNS configuration instructions

Don't forget to update the CORS settings after adding custom domains!

## Troubleshooting

### Backend Issues

**Problem**: Database connection fails
```bash
# Check database status
fly status -a poker-planning-db

# Verify secrets are set
fly secrets list -a obi-poker-planning

# Check logs for detailed errors
fly logs -a obi-poker-planning
```

**Problem**: CORS errors
```bash
# Verify ALLOWED_ORIGINS includes your frontend URL
fly secrets list -a obi-poker-planning

# Update if needed
fly secrets set ALLOWED_ORIGINS=https://your-app.vercel.app -a obi-poker-planning
```

### Frontend Issues

**Problem**: Cannot connect to backend
- Verify `NEXT_PUBLIC_API_URL` is set correctly in Vercel environment variables
- Check that the backend URL is accessible: `curl https://obi-poker-planning.fly.dev/api/sessions`
- Ensure CORS is configured correctly on the backend

**Problem**: WebSocket connection fails
- Verify the backend supports WSS (it should with Fly.io's TLS termination)
- Check browser console for detailed error messages
- Ensure the WebSocket URL transformation is correct (https → wss)

## Cost Estimation

### Fly.io (Backend + Database)
- **Free tier**: 3 shared-cpu-1x VMs with 256MB RAM
- **PostgreSQL**: Free for development (1GB storage)
- **Estimated cost**: $0-5/month for small traffic

### Vercel (Frontend)
- **Hobby tier**: Free for personal projects
- **Unlimited bandwidth** for personal use
- **Estimated cost**: $0/month for hobby projects

## Security Recommendations

1. **Use strong database passwords**: Change the default password
2. **Enable HTTPS only**: Already configured in fly.toml
3. **Limit CORS origins**: Only allow your specific frontend domain
4. **Use environment variables**: Never commit secrets to Git
5. **Regular updates**: Keep dependencies up to date

## Scaling Considerations

### Backend (Fly.io)

```bash
# Scale to multiple regions
fly scale count 2 -a obi-poker-planning

# Scale resources
fly scale vm shared-cpu-2x --memory 512 -a obi-poker-planning

# Add auto-scaling
fly autoscale set min=1 max=3 -a obi-poker-planning
```

### Database (Fly.io)

```bash
# Scale database resources
fly scale vm shared-cpu-1x --memory 512 -a poker-planning-db

# Add replicas for high availability
fly postgres create --name poker-planning-db-replica --region iad
```

### Frontend (Vercel)

Vercel automatically scales based on traffic. No manual configuration needed.

## Backup Strategy

### Database Backups

```bash
# Manual backup
fly postgres connect -a poker-planning-db
pg_dump poker_planning > backup.sql

# Or use Fly.io's automated backups (available on paid plans)
```

## Support and Resources

- **Fly.io Docs**: [fly.io/docs](https://fly.io/docs)
- **Vercel Docs**: [vercel.com/docs](https://vercel.com/docs)
- **PostgreSQL Docs**: [postgresql.org/docs](https://www.postgresql.org/docs/)

---

## Quick Reference Commands

```bash
# Backend (Fly.io)
fly deploy -a obi-poker-planning           # Deploy backend
fly logs -a obi-poker-planning             # View logs
fly status -a obi-poker-planning           # Check status
fly ssh console -a obi-poker-planning      # SSH into container

# Database (Fly.io)
fly postgres connect -a poker-planning-db  # Connect to DB
fly proxy 5433:5432 -a poker-planning-db   # Local proxy

# Frontend (Vercel)
vercel --prod                              # Deploy to production
vercel logs                                # View logs
vercel domains                             # Manage domains
```

---

**Congratulations!** 🎉 Your Poker Planning application is now deployed to the cloud!

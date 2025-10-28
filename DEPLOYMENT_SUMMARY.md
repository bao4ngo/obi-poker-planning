# Cloud Deployment Configuration Summary

This document provides a quick reference for all cloud deployment configurations.

## 📦 Files Created for Deployment

### Backend (Fly.io)
- ✅ `back_end/Dockerfile` - Multi-stage Docker build
- ✅ `back_end/fly.toml` - Fly.io app configuration
- ✅ `back_end/.dockerignore` - Docker build exclusions
- ✅ `back_end/.env.example` - Environment variables template
- ✅ `back_end/.gitignore` - Git exclusions
- ✅ `back_end/main.go` - Updated with environment variable support

### Frontend (Vercel)
- ✅ `front_end/vercel.json` - Vercel deployment config
- ✅ `front_end/next.config.js` - Updated with environment support
- ✅ `front_end/.env.example` - Environment variables template
- ✅ `front_end/src/lib/api.ts` - Updated WebSocket URL handling (http/https → ws/wss)
- ✅ `front_end/Dockerfile` - Optional Docker build (if needed)

### Documentation
- ✅ `DEPLOYMENT.md` - Complete deployment guide
- ✅ `QUICKSTART.md` - Local development quick start
- ✅ `README.md` - Updated with deployment info

## 🔧 Environment Variables

### Backend Environment Variables

Required on Fly.io:

```bash
DB_HOST=poker-planning-db.internal
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<your-secure-password>
DB_NAME=poker_planning
ALLOWED_ORIGINS=https://your-app.vercel.app,http://localhost:3000
```

### Frontend Environment Variables

Required on Vercel:

```bash
NEXT_PUBLIC_API_URL=https://your-app-name.fly.dev
```

## 🚀 Deployment Steps Summary

### Step 1: Deploy Backend to Fly.io
```bash
# Login
fly auth login

# Create and setup PostgreSQL
fly postgres create --name poker-planning-db

# Run schema
fly postgres connect -a poker-planning-db < back_end/database/schema.sql

# Deploy backend
cd back_end
fly launch
fly secrets set DB_HOST=poker-planning-db.internal DB_PORT=5432 DB_USER=postgres DB_PASSWORD=<password> DB_NAME=poker_planning
fly deploy
```

### Step 2: Deploy Frontend to Vercel
```bash
# Via CLI
cd front_end
vercel login
vercel env add NEXT_PUBLIC_API_URL
vercel --prod

# Or via Dashboard: https://vercel.com/new
```

### Step 3: Update CORS
```bash
fly secrets set ALLOWED_ORIGINS=https://your-app.vercel.app -a poker-planning-api
```

## 🔍 Key Configuration Changes

### 1. Backend Main.go Changes
- ✅ Environment variable support via `os.Getenv()`
- ✅ Dynamic CORS origins (comma-separated)
- ✅ Configurable port via `PORT` environment variable
- ✅ Database configuration from environment

### 2. Frontend API Changes
- ✅ WebSocket URL transformation: `http→ws`, `https→wss`
- ✅ Environment-based API URL via `NEXT_PUBLIC_API_URL`
- ✅ Production-ready build configuration

### 3. Database Setup
- ✅ PostgreSQL on Fly.io with automated backups
- ✅ Internal networking for secure communication
- ✅ Connection pooling and SSL support

## 🌐 Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                    End Users                        │
└────────────────────┬────────────────────────────────┘
                     │
        ┌────────────┴───────────┐
        │                        │
        ▼                        ▼
┌──────────────┐         ┌──────────────┐
│   Vercel     │         │   Fly.io     │
│  (Frontend)  │◄────────┤  (Backend)   │
│   Next.js    │  HTTPS  │   Golang     │
│              │  WSS    │              │
└──────────────┘         └──────┬───────┘
                                │
                                ▼
                         ┌──────────────┐
                         │   Fly.io     │
                         │  PostgreSQL  │
                         │  (Database)  │
                         └──────────────┘
```

## 🔐 Security Considerations

### Production Checklist
- [ ] Use strong database passwords
- [ ] Enable HTTPS only (configured in fly.toml)
- [ ] Restrict CORS to specific domains
- [ ] Use environment variables for secrets
- [ ] Enable database backups
- [ ] Set up monitoring and alerts
- [ ] Implement rate limiting (optional)
- [ ] Review and update dependencies regularly

## 📊 Cost Estimation

### Free Tier Usage
- **Fly.io**: 3 shared VMs + 1GB PostgreSQL storage = **$0/month**
- **Vercel**: Hobby plan with unlimited bandwidth = **$0/month**
- **Total**: **$0/month** for hobby/personal projects

### Production Scale (Low Traffic)
- **Fly.io**: 2 VMs + 10GB PostgreSQL = **~$10/month**
- **Vercel**: Pro plan = **$20/month**
- **Total**: **~$30/month**

## 🛠️ Useful Commands

### Fly.io Commands
```bash
# Deploy
fly deploy -a poker-planning-api

# Logs
fly logs -a poker-planning-api

# SSH
fly ssh console -a poker-planning-api

# Scale
fly scale count 2 -a poker-planning-api

# Secrets
fly secrets list -a poker-planning-api
fly secrets set KEY=value -a poker-planning-api
```

### Vercel Commands
```bash
# Deploy
vercel --prod

# Logs
vercel logs

# Environment variables
vercel env ls
vercel env add NEXT_PUBLIC_API_URL
```

### Database Commands
```bash
# Connect
fly postgres connect -a poker-planning-db

# Proxy for local access
fly proxy 5433:5432 -a poker-planning-db

# Backup
pg_dump poker_planning > backup.sql
```

## 🔄 Update Workflow

### Updating Backend
1. Make changes locally
2. Test locally
3. Commit to Git
4. Run `fly deploy -a poker-planning-api`
5. Verify with `fly logs`

### Updating Frontend
1. Make changes locally
2. Test locally
3. Commit to Git
4. Push to repository
5. Vercel auto-deploys from Git

## 📈 Monitoring

### Health Checks
- Backend: `https://your-app.fly.dev/api/sessions`
- Frontend: `https://your-app.vercel.app`

### Logs
- Backend: `fly logs -a poker-planning-api`
- Frontend: Vercel Dashboard → Logs
- Database: `fly logs -a poker-planning-db`

## 🆘 Troubleshooting

### Backend Not Starting
```bash
fly logs -a poker-planning-api
fly status -a poker-planning-api
fly secrets list -a poker-planning-api
```

### Database Connection Issues
```bash
fly status -a poker-planning-db
fly postgres connect -a poker-planning-db
```

### CORS Errors
```bash
fly secrets set ALLOWED_ORIGINS=https://your-domain.com -a poker-planning-api
fly apps restart poker-planning-api
```

### WebSocket Not Connecting
- Ensure backend URL uses `https://` (automatically converts to `wss://`)
- Check browser console for errors
- Verify CORS settings include WebSocket headers

## 📚 Additional Resources

- [Fly.io Documentation](https://fly.io/docs)
- [Vercel Documentation](https://vercel.com/docs)
- [Next.js Deployment](https://nextjs.org/docs/deployment)
- [PostgreSQL on Fly.io](https://fly.io/docs/postgres/)

---

**Last Updated**: October 27, 2025
**Status**: Ready for deployment ✅

# Deploy Go Backend to Google Cloud Run

This guide shows how to deploy your Go backend API to Google Cloud Run with:
- ✅ Automatic HTTPS
- ✅ Auto-scaling
- ✅ Pay-per-use pricing
- ✅ No server management
- ✅ Custom domain support (optional)

**Keep:**
- **Database**: Neon (PostgreSQL)
- **Frontend**: Vercel
- **Backend**: Google Cloud Run (this deployment)

## Prerequisites

- Google Cloud SDK installed (`gcloud` command available)
- A GCP project created
- Neon database already set up
- Docker installed (optional, Cloud Run can build for you)
- Go backend code in `back_end/` directory

---

## Step 1: Set Up GCP Project

```bash
# Login to Google Cloud
gcloud auth login

# Set your project ID
export PROJECT_ID="baocanhcut-personal"
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable artifactregistry.googleapis.com

# Set default region (choose one close to your users)
gcloud config set run/region asia-southeast1
```

---

## Step 2: Verify Your Dockerfile

Check if you have a Dockerfile in `back_end/`:

```bash
# From your local machine
ls back_end/Dockerfile
```

If you don't have one or need to update it, create/update `back_end/Dockerfile`:

```dockerfile
# Use official golang image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o obi-poker-api main.go

# Use minimal alpine image for final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/obi-poker-api .

# Expose port (Cloud Run uses PORT env variable)
EXPOSE 8080

# Run the application
CMD ["./obi-poker-api"]
```

---

## Step 3: Update Your Backend for Cloud Run

Cloud Run expects your app to listen on the PORT environment variable.

Check your `main.go` - it should already handle this:

```go
port := getEnv("PORT", "8080")
log.Printf("Server starting on :%s", port)
if err := http.ListenAndServe(":"+port, handler); err != nil {
    log.Fatal("Server failed to start:", err)
}
```

This is already in your code, so you're good! ✅

---

## Step 4: Deploy to Cloud Run

### Option A: Deploy from Source (Easiest - Cloud Run builds for you)

```bash
# Navigate to your backend directory
cd back_end

# Deploy to Cloud Run (Cloud Run will build the Docker image)
gcloud run deploy obi-poker-backend \
  --source . \
  --platform managed \
  --region asia-southeast1 \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars="DB_HOST=ep-calm-voice-a1pxz353.ap-southeast1.aws.neon.tech" \
  --set-env-vars="DB_PORT=5432" \
  --set-env-vars="DB_USER=neondb_owner" \
  --set-env-vars="DB_PASSWORD=npg_Qi9lKObJM5LB" \
  --set-env-vars="DB_NAME=neondb" \
  --set-env-vars="DB_SSLMODE=require" \
  --set-env-vars="ALLOWED_ORIGINS=https://your-vercel-app.vercel.app,http://localhost:3000"
```

**Replace `your-vercel-app.vercel.app` with your actual Vercel URL**

### Option B: Build Docker Image Locally and Deploy

```bash
# Navigate to your backend directory
cd back_end

# Build the Docker image
docker build -t gcr.io/$PROJECT_ID/obi-poker-backend:latest .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/obi-poker-backend:latest

# Deploy to Cloud Run
gcloud run deploy obi-poker-backend \
  --image gcr.io/$PROJECT_ID/obi-poker-backend:latest \
  --platform managed \
  --region asia-southeast1 \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars="DB_HOST=ep-calm-voice-a1pxz353.ap-southeast1.aws.neon.tech,DB_PORT=5432,DB_USER=neondb_owner,DB_PASSWORD=npg_Qi9lKObJM5LB,DB_NAME=neondb,DB_SSLMODE=require,ALLOWED_ORIGINS=https://your-vercel-app.vercel.app"
```

---

## Step 5: Get Your Cloud Run URL

After deployment completes, you'll see output like:

```
Service [obi-poker-backend] revision [obi-poker-backend-00001-xyz] has been deployed and is serving 100 percent of traffic.
Service URL: https://obi-poker-backend-1234567890-as.a.run.app
```

**Copy this URL!** This is your backend URL with automatic HTTPS ✅

Or get it anytime with:

```bash
gcloud run services describe obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --format 'value(status.url)'
```

---

## Step 6: Test Your Deployment

```bash
# Get your Cloud Run URL
export BACKEND_URL=$(gcloud run services describe obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --format 'value(status.url)')

echo "Testing backend at: $BACKEND_URL"

# Test health endpoint (if you have one)
curl $BACKEND_URL/health

# Test sessions endpoint
curl $BACKEND_URL/api/sessions
```

---

## Step 7: Update Vercel Frontend

Update your Vercel environment variables:

1. Go to your Vercel dashboard
2. Select your project
3. Go to Settings → Environment Variables
4. Update or add:
   - Name: `NEXT_PUBLIC_API_URL`
   - Value: `https://obi-poker-backend-1234567890-as.a.run.app` (your Cloud Run URL)
5. Redeploy your frontend

---

## Step 8: Update CORS Settings (Important!)

Make sure your backend allows requests from your Vercel domain.

Update your Cloud Run service with the correct ALLOWED_ORIGINS:

```bash
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --set-env-vars="ALLOWED_ORIGINS=https://your-vercel-app.vercel.app,http://localhost:3000"
```

---

## Step 9: (Optional) Set Up Custom Domain

If you have a domain, you can map it to Cloud Run:

```bash
# Map your domain to Cloud Run
gcloud run domain-mappings create \
  --service obi-poker-backend \
  --domain api.yourdomain.com \
  --region asia-southeast1
```

Then add the DNS records shown in the output to your domain provider.

---

## Managing Environment Variables

### View Current Environment Variables

```bash
gcloud run services describe obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --format='value(spec.template.spec.containers[0].env)'
```

### Update Environment Variables

```bash
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --update-env-vars="DB_HOST=new-host,DB_PORT=5432"
```

### Use Secret Manager (Recommended for Sensitive Data)

```bash
# Create a secret
echo -n "npg_Qi9lKObJM5LB" | gcloud secrets create db-password --data-file=-

# Grant Cloud Run access to the secret
gcloud secrets add-iam-policy-binding db-password \
  --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Update Cloud Run to use the secret
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --set-secrets="DB_PASSWORD=db-password:latest"
```

---

## Updating Your Backend

When you make code changes:

```bash
# Option 1: Deploy from source (easiest)
cd back_end
gcloud run deploy obi-poker-backend \
  --source . \
  --platform managed \
  --region asia-southeast1

# Option 2: Build and deploy Docker image
docker build -t gcr.io/$PROJECT_ID/obi-poker-backend:latest .
docker push gcr.io/$PROJECT_ID/obi-poker-backend:latest
gcloud run deploy obi-poker-backend \
  --image gcr.io/$PROJECT_ID/obi-poker-backend:latest \
  --platform managed \
  --region asia-southeast1
```

---

## Viewing Logs

```bash
# View recent logs
gcloud run services logs read obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --limit 50

# Stream logs in real-time
gcloud run services logs tail obi-poker-backend \
  --platform managed \
  --region asia-southeast1

# View logs in Cloud Console
# Go to: https://console.cloud.google.com/run
```

---

## Scaling Configuration

Cloud Run auto-scales by default. You can configure limits:

```bash
# Set minimum and maximum instances
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --min-instances 0 \
  --max-instances 10

# Set concurrency (requests per instance)
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --concurrency 80

# Set CPU and memory
gcloud run services update obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --cpu 1 \
  --memory 512Mi
```

**Free Tier Settings (Recommended for Testing):**
- `--min-instances 0` - Scale to zero when not used
- `--max-instances 1` - Limit to 1 instance
- `--cpu 1` - 1 vCPU
- `--memory 512Mi` - 512MB RAM

---

## Cost Estimation

**Cloud Run Pricing (Free Tier):**
- 2 million requests per month - FREE
- 360,000 GB-seconds of memory - FREE
- 180,000 vCPU-seconds - FREE

**After free tier:**
- ~$0.00002400 per request
- ~$0.00000250 per GB-second
- ~$0.00001000 per vCPU-second

**Estimated cost for small app:** $0-5/month

Compare to VM (e2-micro): $7-10/month always running

---

## Useful Commands

```bash
# List all Cloud Run services
gcloud run services list --platform managed

# Get service details
gcloud run services describe obi-poker-backend \
  --platform managed \
  --region asia-southeast1

# Get service URL
gcloud run services describe obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --format 'value(status.url)'

# Delete service
gcloud run services delete obi-poker-backend \
  --platform managed \
  --region asia-southeast1

# View revisions
gcloud run revisions list \
  --service obi-poker-backend \
  --platform managed \
  --region asia-southeast1

# Rollback to previous revision
gcloud run services update-traffic obi-poker-backend \
  --to-revisions REVISION_NAME=100 \
  --platform managed \
  --region asia-southeast1
```

---

## Troubleshooting

### Service won't start

```bash
# Check logs
gcloud run services logs read obi-poker-backend \
  --platform managed \
  --region asia-southeast1 \
  --limit 100

# Common issues:
# 1. Port mismatch - Make sure your app listens on PORT env variable
# 2. Database connection - Check DB_HOST, DB_PASSWORD, etc.
# 3. Missing dependencies - Check your Dockerfile
```

### Database connection errors

```bash
# Test database connection from Cloud Shell
psql "postgresql://neondb_owner:npg_Qi9lKObJM5LB@ep-calm-voice-a1pxz353.ap-southeast1.aws.neon.tech/neondb?sslmode=require"

# Make sure DB_SSLMODE=require is set
gcloud run services update obi-poker-backend \
  --update-env-vars="DB_SSLMODE=require"
```

### CORS errors

```bash
# Update ALLOWED_ORIGINS to include your Vercel domain
gcloud run services update obi-poker-backend \
  --update-env-vars="ALLOWED_ORIGINS=https://your-vercel-app.vercel.app"
```

---

## CI/CD with GitHub Actions (Optional)

Create `.github/workflows/deploy-cloud-run.yml`:

```yaml
name: Deploy to Cloud Run

on:
  push:
    branches:
      - main
    paths:
      - 'back_end/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - id: 'auth'
      uses: 'google-github-actions/auth@v1'
      with:
        credentials_json: '${{ secrets.GCP_SA_KEY }}'
    
    - name: 'Deploy to Cloud Run'
      uses: 'google-github-actions/deploy-cloudrun@v1'
      with:
        service: 'obi-poker-backend'
        region: 'asia-southeast1'
        source: './back_end'
        env_vars: |
          DB_HOST=ep-calm-voice-a1pxz353.ap-southeast1.aws.neon.tech
          DB_PORT=5432
          DB_USER=neondb_owner
          DB_NAME=neondb
          DB_SSLMODE=require
        secrets: |
          DB_PASSWORD=db-password:latest
```

---

## Advantages of Cloud Run vs VM

| Feature | Cloud Run | VM (e2-micro) |
|---------|-----------|---------------|
| HTTPS | ✅ Automatic | ❌ Manual setup |
| Scaling | ✅ Automatic | ❌ Manual |
| Cost | ✅ Pay per use | ❌ Always on |
| Maintenance | ✅ None | ❌ OS updates |
| Cold starts | ⚠️ Yes (~1s) | ✅ No |
| Always available | ✅ Yes | ✅ Yes |

---

## Next Steps

1. ✅ Backend deployed on Cloud Run with HTTPS
2. ✅ Update Vercel with Cloud Run URL
3. ✅ Test your application
4. Consider: Set up monitoring and alerting
5. Consider: Add CI/CD with GitHub Actions
6. Consider: Use Cloud Run with custom domain

---

**Your backend is now live with automatic HTTPS!** 🎉

No more certificate issues, no server management, just deploy and go!

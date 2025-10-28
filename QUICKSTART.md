# Quick Start Guide

This guide will help you get the application running locally in minutes.

## Prerequisites Check

Make sure you have installed:
- [ ] Go 1.21+ (`go version`)
- [ ] Node.js 18+ (`node --version`)
- [ ] PostgreSQL 12+ (`psql --version`)
- [ ] npm (`npm --version`)

## Quick Setup (5 minutes)

### 1. Clone and Navigate
```bash
cd first_api_go
```

### 2. Start PostgreSQL
Make sure PostgreSQL is running:
- **Windows**: Start PostgreSQL service from Services
- **Mac**: `brew services start postgresql`
- **Linux**: `sudo systemctl start postgresql`

### 3. Setup Backend
```bash
cd back_end

# Install dependencies
go mod download

# Setup database (creates DB and tables)
go run cmd/setup/main.go

# Start backend server
go run main.go
```

Keep this terminal open. Backend will run on http://localhost:8080

### 4. Setup Frontend (New Terminal)
```bash
cd front_end

# Install dependencies
npm install

# Start frontend server
npm run dev
```

Keep this terminal open. Frontend will run on http://localhost:3000

### 5. Access Application
Open your browser and go to: **http://localhost:3000**

## Verify Setup

### Test Backend
```bash
curl http://localhost:8080/api/sessions
```
Should return: `{"sessions":[]}`

### Test Frontend
1. Open http://localhost:3000
2. Enter your name: "Test Host"
3. Enter session name: "Test Session"
4. Click "Create Session"
5. You should see the session page with voting cards

## Troubleshooting

### Database Connection Failed
```bash
# Check PostgreSQL is running
# Windows
Get-Service postgresql*

# Mac/Linux
ps aux | grep postgres

# Test connection
psql -U postgres -d poker_planning -c "SELECT version();"
```

### Port Already in Use
```bash
# Backend (8080)
# Windows
netstat -ano | findstr :8080
# Kill process: taskkill /PID <PID> /F

# Mac/Linux
lsof -ti:8080 | xargs kill

# Frontend (3000)
# Windows
netstat -ano | findstr :3000
# Kill process: taskkill /PID <PID> /F

# Mac/Linux
lsof -ti:3000 | xargs kill
```

### Module Not Found (Backend)
```bash
cd back_end
go mod tidy
go mod download
```

### Module Not Found (Frontend)
```bash
cd front_end
rm -rf node_modules package-lock.json
npm install
```

## Common Commands

### Backend
```bash
# Run server
go run main.go

# Build executable
go build -o poker-api

# Run tests
go test ./...

# Reset database
go run cmd/reset/main.go
```

### Frontend
```bash
# Development
npm run dev

# Production build
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

## Environment Configuration

### Backend (.env)
Create `back_end/.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=poker_planning
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.local)
Create `front_end/.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Next Steps

1. ✅ Local development is working
2. 📖 Read [DEPLOYMENT.md](./DEPLOYMENT.md) to deploy to cloud
3. 🎯 Customize the application for your needs
4. 🔒 Update database password in production

## Need Help?

- Check [README.md](./README.md) for full documentation
- Check [DEPLOYMENT.md](./DEPLOYMENT.md) for cloud deployment
- Review [POSTGRESQL_SETUP.md](./POSTGRESQL_SETUP.md) for database setup details

---

**Happy Planning! 🎯**

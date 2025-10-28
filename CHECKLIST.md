# Pre-Deployment Checklist

Use this checklist before deploying to production.

## 📋 Backend Checklist

### Code Review
- [ ] All environment variables are read from `os.Getenv()`
- [ ] No hardcoded credentials in code
- [ ] CORS origins are configurable
- [ ] Database connection uses connection pooling
- [ ] Error handling is comprehensive
- [ ] Logging is implemented for debugging

### Configuration Files
- [ ] `Dockerfile` is present and tested
- [ ] `fly.toml` is configured with correct region
- [ ] `.dockerignore` excludes unnecessary files
- [ ] `.gitignore` excludes sensitive files
- [ ] `.env.example` is up to date

### Database
- [ ] Schema is tested and working
- [ ] Migration scripts are ready
- [ ] Indexes are optimized
- [ ] Backup strategy is planned

### Testing
- [ ] All API endpoints tested
- [ ] WebSocket connections tested
- [ ] Database queries tested
- [ ] Error scenarios tested

## 📋 Frontend Checklist

### Code Review
- [ ] Environment variables use `NEXT_PUBLIC_` prefix
- [ ] API URL is configurable
- [ ] WebSocket URL handles http/https → ws/wss conversion
- [ ] Error handling is user-friendly
- [ ] Loading states are implemented

### Configuration Files
- [ ] `vercel.json` is present
- [ ] `next.config.js` is configured
- [ ] `.env.example` is up to date
- [ ] `.gitignore` excludes sensitive files

### Build & Performance
- [ ] Production build completes without errors
- [ ] No console errors in production build
- [ ] Images are optimized
- [ ] Bundle size is reasonable

### Testing
- [ ] All pages load correctly
- [ ] Forms submit properly
- [ ] WebSocket connections work
- [ ] Responsive design on mobile/tablet/desktop
- [ ] Cross-browser testing done

## 📋 Deployment Process Checklist

### Pre-Deployment
- [ ] Code is committed to Git
- [ ] README.md is updated
- [ ] DEPLOYMENT.md is reviewed
- [ ] All dependencies are up to date
- [ ] Version numbers are updated (if applicable)

### Fly.io Setup
- [ ] Fly.io account created
- [ ] Fly CLI installed and logged in
- [ ] PostgreSQL database created
- [ ] Database schema applied
- [ ] Environment secrets configured
- [ ] CORS origins include production URL

### Vercel Setup
- [ ] Vercel account created
- [ ] Repository connected to Vercel
- [ ] Environment variables set
- [ ] Build settings configured
- [ ] Domain configured (if custom domain)

### Post-Deployment
- [ ] Backend health check passes
- [ ] Frontend loads without errors
- [ ] WebSocket connection works
- [ ] Can create a session
- [ ] Can join a session
- [ ] Can vote and reveal votes
- [ ] Data persists after refresh
- [ ] Error handling works as expected

## 📋 Security Checklist

### Credentials & Secrets
- [ ] Strong database password set
- [ ] Environment variables used for all secrets
- [ ] No secrets committed to Git
- [ ] `.env` files are in `.gitignore`

### Network Security
- [ ] HTTPS enforced (fly.toml: force_https = true)
- [ ] CORS restricted to specific origins
- [ ] Database not publicly accessible
- [ ] WebSocket uses WSS in production

### Application Security
- [ ] Input validation implemented
- [ ] SQL injection prevention (using parameterized queries)
- [ ] XSS prevention (React's built-in protection)
- [ ] CSRF protection (if needed)
- [ ] Rate limiting considered (optional)

## 📋 Monitoring & Maintenance Checklist

### Monitoring Setup
- [ ] Fly.io metrics dashboard checked
- [ ] Vercel analytics configured (optional)
- [ ] Error tracking setup (optional: Sentry, etc.)
- [ ] Uptime monitoring configured (optional)

### Backup Strategy
- [ ] Database backup plan in place
- [ ] Backup restoration tested
- [ ] Code is in version control

### Documentation
- [ ] Deployment process documented
- [ ] Environment variables documented
- [ ] Troubleshooting guide available
- [ ] Team trained on deployment process

## 📋 Cost Management Checklist

### Resource Optimization
- [ ] Using free tier where possible
- [ ] Auto-scaling configured appropriately
- [ ] Database size is reasonable
- [ ] Unused resources cleaned up

### Monitoring Costs
- [ ] Fly.io billing dashboard checked
- [ ] Vercel usage limits understood
- [ ] Alerts set for unusual usage
- [ ] Budget limits configured

## 📋 Final Pre-Launch Checklist

### Functionality
- [ ] Create session works
- [ ] Join session works
- [ ] Add planning items works
- [ ] Vote submission works
- [ ] Vote reveal works
- [ ] Final estimate setting works
- [ ] User connection status accurate
- [ ] Sessions persist across restarts

### User Experience
- [ ] Loading states show properly
- [ ] Error messages are clear
- [ ] Success feedback is visible
- [ ] Navigation is intuitive
- [ ] Mobile experience is good

### Performance
- [ ] Page load time < 3 seconds
- [ ] WebSocket reconnection works
- [ ] No memory leaks
- [ ] Database queries are fast

### Compliance (if applicable)
- [ ] Privacy policy reviewed
- [ ] Terms of service reviewed
- [ ] Data retention policy defined
- [ ] GDPR compliance (if applicable)

## 🎉 Launch Day

- [ ] Notify team of deployment
- [ ] Share production URL
- [ ] Monitor logs for first hour
- [ ] Test all critical functionality
- [ ] Be ready for quick rollback if needed
- [ ] Celebrate! 🎊

## 📞 Emergency Contacts

Document these for your team:

- **Fly.io Status**: https://status.fly.io
- **Vercel Status**: https://www.vercel-status.com
- **Database Admin**: [Your contact]
- **DevOps Lead**: [Your contact]

---

**Tip**: Print this checklist and check items as you go. Keep a copy for future deployments!

**Status**: Ready to deploy when all items are checked ✅

# Agile Poker Planning Application

A full-stack web application for Agile Poker Planning sessions with real-time collaboration.

## 🎯 Features

- **Create Sessions**: Host can create planning sessions with custom names
- **Join Sessions**: Team members can join using their names
- **Real-time Updates**: WebSocket-based real-time synchronization
- **Planning Items**: Add and manage multiple items to estimate
- **Voting System**: Standard Fibonacci sequence poker cards (0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, ?)
- **Vote Reveal**: Host can reveal votes when everyone has voted
- **Final Estimates**: Set and track final estimates for each item
- **Participant Tracking**: See who's connected and who has voted
- **Persistent Storage**: PostgreSQL database for session persistence
- **Username Validation**: Case-insensitive unique usernames per session

## 🏗️ Architecture

### Backend (Golang)
- RESTful API with Gorilla Mux
- WebSocket support with Gorilla WebSocket
- PostgreSQL database with lib/pq driver
- Session persistence and recovery
- CORS enabled for frontend integration

### Frontend (Next.js)
- TypeScript for type safety
- Tailwind CSS for styling
- Real-time WebSocket client
- Responsive design
- Environment-based configuration

## 🚀 Deployment

This application can be deployed to the cloud:
- **Frontend**: Vercel
- **Backend**: Render.com or Fly.io
- **Database**: Render.com PostgreSQL or Fly.io PostgreSQL

**📖 Deployment Guides:**
- **[RENDER_QUICKSTART.md](./RENDER_QUICKSTART.md)** - 🚀 Quick deploy to Render.com (Easiest!)
- **[RENDER_DEPLOYMENT.md](./RENDER_DEPLOYMENT.md)** - Complete Render.com guide
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Fly.io deployment guide

## �📁 Project Structure

```
first_api_go/
├── back_end/              # Golang backend
│   ├── main.go            # Server entry point
│   ├── go.mod             # Go module definition
│   ├── Dockerfile         # Docker configuration for Fly.io
│   ├── fly.toml           # Fly.io deployment config
│   ├── .env.example       # Environment variables template
│   ├── models/
│   │   └── models.go      # Data models
│   ├── handlers/
│   │   ├── session.go     # REST API handlers
│   │   └── websocket.go   # WebSocket handlers
│   ├── db/
│   │   ├── db.go          # Database connection
│   │   └── queries.go     # Database queries
│   ├── database/
│   │   ├── schema.sql     # Database schema
│   │   └── drop.sql       # Database reset script
│   └── cmd/
│       ├── setup/         # Database setup tool
│       └── reset/         # Database reset tool
│
└── front_end/             # Next.js frontend
    ├── src/
    │   ├── pages/         # Next.js pages
    │   ├── lib/           # API client
    │   ├── types/         # TypeScript types
    │   └── styles/        # Global styles
    ├── package.json
    ├── next.config.js     # Next.js configuration
    ├── vercel.json        # Vercel deployment config
    └── .env.example       # Environment variables template
```

## � Documentation

- **[READY_TO_DEPLOY.md](./READY_TO_DEPLOY.md)** - 🚀 Start here for cloud deployment
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Complete deployment guide (Fly.io + Vercel)
- **[QUICKSTART.md](./QUICKSTART.md)** - Quick start for local development
- **[CHECKLIST.md](./CHECKLIST.md)** - Pre-deployment checklist
- **[DEPLOYMENT_SUMMARY.md](./DEPLOYMENT_SUMMARY.md)** - Quick reference for configs
- **[POSTGRESQL_SETUP.md](./POSTGRESQL_SETUP.md)** - Database setup details

## �🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 18.x or higher
- npm or yarn
- **PostgreSQL 12 or higher**

### Local Development Setup

#### Backend Setup

1. Navigate to the backend directory:
```bash
cd back_end
```

2. Install Go dependencies:
```bash
go mod download
```

3. **Set up environment variables:**
```bash
# Copy the example file
cp .env.example .env

# Edit .env if needed (default values work for local PostgreSQL)
```

4. **Set up the database (Cross-platform):**
```bash
go run cmd/setup/main.go
```

This will create the PostgreSQL database and all tables.

5. Run the server:
```bash
go run main.go
```

The backend will start on `http://localhost:8080`

#### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd front_end
```

2. Install dependencies:
```bash
npm install
```

3. **Set up environment variables (optional for local development):**
```bash
# Copy the example file
cp .env.example .env.local

# Default API_URL (http://localhost:8080) will be used if not set
```

4. Run the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 🎮 How to Use

### For Hosts

1. Go to `http://localhost:3000`
2. Enter your name
3. Enter a session name
4. Click "Create Session"
5. Share the session link with your team
6. Add planning items
7. Select an item to start voting
8. Reveal votes when everyone has voted
9. Set the final estimate
10. Move to the next item

### For Participants
1. Get the session link from the host (or enter session ID)
2. Enter your name
3. Join the session
4. Wait for the host to select an item
5. Vote using the poker cards
6. View results when revealed
7. Discuss and agree on the final estimate

## 🔌 API Endpoints

### REST API

- `POST /api/sessions` - Create a new session
- `GET /api/sessions` - Get all sessions
- `GET /api/sessions/{sessionId}` - Get session details
- `POST /api/sessions/{sessionId}/items` - Add a planning item
- `POST /api/sessions/{sessionId}/current-item` - Set current item

### WebSocket

- `WS /ws/{sessionId}` - Connect to session for real-time updates

## 🛠️ Technologies Used

### Backend
- **Go**: Programming language
- **Gorilla Mux**: HTTP router
- **Gorilla WebSocket**: WebSocket implementation
- **RS CORS**: CORS middleware
- **Google UUID**: UUID generation

### Frontend
- **Next.js**: React framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Utility-first CSS
- **WebSocket API**: Real-time communication

## 📝 Card Values

The application uses standard Fibonacci sequence for estimation:
- **0**: No effort
- **1, 2, 3, 5, 8, 13, 21, 34, 55, 89**: Story points
- **?**: Unknown/Need more information

## 🔒 Security Notes

- The current implementation is for development purposes
- WebSocket allows all origins (should be restricted in production)
- No authentication implemented (add as needed)
- Sessions are stored in memory (use a database for production)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is open source and available under the MIT License.

## 🐛 Known Issues

- Sessions are lost on server restart (in-memory storage)
- No persistence layer
- No authentication/authorization

## 🚧 Future Enhancements

- [ ] User authentication
- [ ] Persistent storage (database)
- [ ] Session history
- [ ] Export results
- [ ] Custom card values
- [ ] Timer for voting rounds
- [ ] Chat functionality
- [ ] Multiple estimation methods
- [ ] Admin dashboard
- [ ] Session analytics

## 📞 Support

For issues, questions, or contributions, please open an issue in the repository.

---

Happy Planning! 🎯🃏

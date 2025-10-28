# Poker Planning Frontend

A Next.js-based frontend for Agile Poker Planning with real-time WebSocket communication.

## Features

- Create new poker planning sessions
- Join existing sessions
- Real-time voting with WebSocket
- Host controls for revealing and resetting votes
- Set final estimates
- Participant tracking
- Responsive design with Tailwind CSS

## Prerequisites

- Node.js 18.x or higher
- npm or yarn

## Installation

1. Navigate to the frontend directory:
```bash
cd front_end
```

2. Install dependencies:
```bash
npm install
```

## Configuration

The application expects the backend API to be running on `http://localhost:8080` by default. You can change this by creating a `.env.local` file:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Running the Application

Development mode:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`

Build for production:
```bash
npm run build
npm start
```

## Usage

### Creating a Session

1. Enter your name on the home page
2. Enter a session name
3. Click "Create Session"
4. Share the session link with participants

### Joining a Session

1. Enter your name on the home page
2. Click "Join Existing Session"
3. Enter the session ID when prompted

### Host Features

- Add planning items
- Select the current item for voting
- Reveal votes
- Reset votes
- Set final estimates

### Participant Features

- Vote on the current item
- View other participants
- See when votes are revealed
- View final estimates

## Tech Stack

- **Framework**: Next.js 14
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Real-time Communication**: WebSocket API
- **State Management**: React Hooks

## Project Structure

```
front_end/
├── src/
│   ├── pages/
│   │   ├── _app.tsx           # App wrapper
│   │   ├── index.tsx           # Home page (create/join session)
│   │   └── session/
│   │       └── [sessionId].tsx # Session page
│   ├── lib/
│   │   └── api.ts              # API client
│   ├── types/
│   │   └── index.ts            # TypeScript types
│   └── styles/
│       └── globals.css         # Global styles
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── next.config.js
```

## WebSocket Events

The application handles the following WebSocket events:

- `welcome` - Initial connection
- `user_joined` - New user joined
- `user_left` - User disconnected
- `item_added` - New item added
- `vote_submitted` - Vote cast
- `votes_revealed` - Votes revealed by host
- `votes_reset` - Votes reset by host
- `current_item_changed` - Active item changed
- `final_estimate_set` - Final estimate set

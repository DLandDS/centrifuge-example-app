# Centrifuge Chat Application

A real-time topic-based messaging application built with:
- **Centrifuge** - Real-time messaging server
- **Go Backend** - JWT authentication and API endpoints  
- **SvelteKit Frontend** - Modern web interface

## Architecture

```
├── centrifuge/     # Docker Compose setup for Centrifuge server
├── backend/        # Go server with JWT authentication
├── frontend/       # SvelteKit web application
└── README.md
```

## Quick Start

### 1. Start Centrifuge Server

```bash
cd centrifuge
docker-compose up -d
```

The Centrifuge server will be available at:
- WebSocket: `ws://localhost:8000/connection/websocket`
- Admin Panel: http://localhost:8000/admin (admin/admin)

### 2. Start Go Backend

```bash
cd backend
go run main.go
```

The API server will be available at http://localhost:3001

### 3. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

The web application will be available at http://localhost:5173

## Features

- **Real-time messaging** via WebSocket connection to Centrifuge
- **Topic-based channels** (general, tech, random)
- **JWT authentication** with secure token generation
- **Persistent login** with localStorage
- **Responsive UI** with modern design

## Demo Usage

1. Open http://localhost:5173 in your browser
2. Login with any username/password (demo mode)
3. Select a topic from the sidebar
4. Start chatting in real-time!

You can open multiple browser tabs/windows to test multi-user functionality.

## API Endpoints

- `POST /api/login` - User authentication
- `GET /api/user` - Get user information (protected)
- `POST /api/centrifuge-token` - Refresh Centrifuge token (protected)
- `GET /api/health` - Health check

## Configuration

### Centrifuge
- Port: 8000
- Secret Key: `my-secret-key`
- Namespaces: `chat:*` with history enabled

### Backend
- Port: 3001  
- JWT Secret: `jwt-secret-key`
- Centrifuge Secret: `my-secret-key` (matches Centrifuge)

### Frontend
- Port: 5173 (dev mode)
- Connects to backend at http://localhost:3001
- Connects to Centrifuge at ws://localhost:8000

## Development

Each component can be developed independently:

- **Centrifuge**: Modify `centrifuge/config.json` for server settings
- **Backend**: Standard Go development with live reload
- **Frontend**: SvelteKit with hot reload and TypeScript support

## Production Notes

For production deployment:
1. Use environment variables for secrets
2. Configure proper CORS settings
3. Use HTTPS/WSS connections
4. Set up proper authentication database
5. Configure reverse proxy (nginx/traefik)
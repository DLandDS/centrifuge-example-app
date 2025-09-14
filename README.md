# Centrifuge Topic-Based Messaging Application

A real-time topic-based messaging application built with Go backend (using Centrifuge) and Svelte frontend. Features JWT authentication and direct WebSocket connections for real-time messaging.

## Features

- 🚀 Real-time messaging using Centrifuge WebSocket server
- 🔐 JWT-based authentication
- 📨 Topic-based messaging (channels)
- 💬 Multiple chat rooms/topics
- 🎨 Modern Svelte frontend
- 🐳 Docker support

## Architecture

- **Backend**: Go with Centrifuge server for real-time messaging
- **Frontend**: Svelte with direct WebSocket connection to Centrifuge
- **Authentication**: JWT tokens for secure communication
- **Communication**: WebSocket connection from frontend directly to Centrifuge

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Make (optional, for using Makefile commands)

### Option 1: Using Make (Recommended)

```bash
# Setup everything from scratch
make setup

# Run in development mode
make dev
```

### Option 2: Manual Setup

```bash
# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd frontend && npm install

# Build frontend
cd frontend && npm run build
cd ..

# Run the application
go run main.go
```

### Option 3: Using Docker

```bash
docker build -t centrifuge-app .
docker run -p 8080:8080 centrifuge-app
```

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. Enter any username (password is optional for this demo)
3. Click "Login" to authenticate and connect to Centrifuge
4. Select a topic/channel from the available topics
5. Start sending messages in real-time!

## Available Topics

- `#general` - General discussion
- `#tech` - Technology discussions
- `#random` - Random conversations
- `#announcements` - Important announcements

## API Endpoints

### Authentication

- `POST /api/login` - Login and get JWT token
  ```json
  {
    "username": "your_username",
    "password": "optional_password"
  }
  ```

### Messaging

- `POST /api/topics/{topic}/messages` - Send message to a topic (requires JWT)
  ```json
  {
    "content": "Your message content"
  }
  ```

### WebSocket

- `WS /connection/websocket` - Centrifuge WebSocket endpoint (requires JWT token)

## Development

### Project Structure

```
.
├── main.go                 # Go backend server
├── go.mod                  # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── App.svelte     # Main Svelte component
│   │   └── main.js        # Frontend entry point
│   ├── public/
│   │   └── index.html     # HTML template
│   ├── package.json       # Frontend dependencies
│   └── rollup.config.js   # Build configuration
├── Makefile               # Build automation
├── Dockerfile             # Container configuration
└── README.md              # This file
```

### Building

```bash
# Build everything
make build

# Build only frontend
make build-frontend

# Clean build artifacts
make clean
```

### Environment Variables

You can customize the application using environment variables:

- `JWT_SECRET` - Secret key for JWT tokens (default: "your-secret-key")
- `PORT` - Server port (default: 8080)

## How It Works

1. **Authentication**: Users login with username/password and receive a JWT token
2. **Connection**: Frontend connects to Centrifuge WebSocket endpoint with JWT token
3. **Authorization**: Centrifuge validates JWT token on connection
4. **Subscription**: Users can subscribe to different topics/channels
5. **Messaging**: Messages are published to Centrifuge topics and broadcast to all subscribers
6. **Real-time**: All connected clients receive messages instantly via WebSocket

## Security Features

- JWT token validation on WebSocket connection
- Token-based API authentication
- CORS configuration for cross-origin requests
- Input validation and sanitization

## Technologies Used

### Backend
- [Centrifuge](https://github.com/centrifugal/centrifuge) - Real-time messaging server
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [CORS](https://github.com/rs/cors) - CORS middleware

### Frontend
- [Svelte](https://svelte.dev/) - Frontend framework
- [Centrifuge-js](https://github.com/centrifugal/centrifuge-js) - WebSocket client
- [Rollup](https://rollupjs.org/) - Module bundler

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
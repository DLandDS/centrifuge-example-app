# Go Backend Server

This is the Go backend server that provides JWT authentication for the Centrifuge messaging application.

## Features

- JWT-based authentication
- Centrifuge token generation
- CORS support for frontend
- RESTful API endpoints

## API Endpoints

### Public Endpoints

- `POST /api/login` - User login
- `GET /api/health` - Health check

### Protected Endpoints (require JWT token)

- `GET /api/user` - Get user information
- `POST /api/centrifuge-token` - Refresh Centrifuge token

## Quick Start

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

The server will start on port 3001.

## Environment

- **Port**: 3001
- **JWT Secret**: jwt-secret-key
- **Centrifuge Secret**: my-secret-key (matches centrifuge config)

## Authentication Flow

1. User logs in with username/password
2. Backend returns JWT token and Centrifuge token
3. Frontend uses JWT token for API calls
4. Frontend uses Centrifuge token to connect to Centrifuge WebSocket
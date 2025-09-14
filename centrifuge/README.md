# Centrifuge Setup

This directory contains the Docker Compose configuration for running Centrifuge server.

## Quick Start

1. Start Centrifuge server:
   ```bash
   docker-compose up -d
   ```

2. Access admin panel at http://localhost:8000/admin
   - Username: admin
   - Password: admin

3. Stop the server:
   ```bash
   docker-compose down
   ```

## Configuration

- **Port**: 8000
- **Admin Panel**: Enabled at `/admin`
- **WebSocket**: Available at `ws://localhost:8000/connection/websocket`
- **API**: Available at `http://localhost:8000/api`

## Available Namespaces

- `chat:*` - For chat messages with history enabled
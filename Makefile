.PHONY: build run dev install-frontend build-frontend clean

# Install Go dependencies
install:
	go mod tidy

# Install frontend dependencies
install-frontend:
	cd frontend && npm install

# Build frontend
build-frontend:
	cd frontend && npm run build

# Build Go binary
build: install build-frontend
	go build -o bin/centrifuge-app main.go

# Run in development mode
dev:
	go run main.go

# Run production build
run: build
	./bin/centrifuge-app

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf frontend/public/build/
	rm -rf frontend/node_modules/

# Setup everything from scratch
setup: install install-frontend build-frontend
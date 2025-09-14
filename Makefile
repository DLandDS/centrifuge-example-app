.PHONY: build run dev install-frontend build-frontend clean test

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

# Run tests
test:
	go test -v

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

# Docker commands
docker-build:
	docker build -t centrifuge-app .

docker-run:
	docker run -p 8080:8080 centrifuge-app

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down
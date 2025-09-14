package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/centrifugal/centrifuge"
)

const (
	jwtSecret = "jwt-secret-key"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ChatMessage struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Topic     string    `json:"topic"`
}

// Simple JWT-like token generation (for demo purposes)
func generateSimpleToken(userID, username string) string {
	payload := fmt.Sprintf("%s:%s:%d", userID, username, time.Now().Add(24*time.Hour).Unix())
	h := hmac.New(sha256.New, []byte(jwtSecret))
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s.%s", hex.EncodeToString([]byte(payload)), signature)
}

// Validate simple token
func validateToken(tokenString string) (string, string, bool) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return "", "", false
	}
	
	payloadBytes, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", "", false
	}
	
	h := hmac.New(sha256.New, []byte(jwtSecret))
	h.Write(payloadBytes)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	
	if parts[1] != expectedSignature {
		return "", "", false
	}
	
	payload := string(payloadBytes)
	parts = strings.Split(payload, ":")
	if len(parts) != 3 {
		return "", "", false
	}
	
	// Check expiration (simplified)
	return parts[0], parts[1], true
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

// Auth middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		userID, username, valid := validateToken(tokenString)
		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("UserID", userID)
		r.Header.Set("Username", username)

		next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple user validation (in real app, check against database)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// For demo purposes, accept any non-empty credentials
	userID := fmt.Sprintf("user_%s", req.Username)
	user := User{
		ID:       userID,
		Username: req.Username,
		Email:    fmt.Sprintf("%s@example.com", req.Username),
	}

	// Generate token
	jwtToken := generateSimpleToken(userID, req.Username)

	response := LoginResponse{
		Token: jwtToken,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	username := r.Header.Get("Username")

	user := User{
		ID:       userID,
		Username: username,
		Email:    fmt.Sprintf("%s@example.com", username),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Centrifuge event handlers
func authHandler(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
	// Get token from event
	token := e.Token
	if token == "" {
		return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}
	
	_, username, valid := validateToken(token)
	if !valid {
		return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
	}

	return centrifuge.ConnectReply{
		Data: []byte(fmt.Sprintf(`{"user":"%s"}`, username)),
		Subscriptions: map[string]centrifuge.SubscribeOptions{
			"chat:general": {},
			"chat:tech":    {},
			"chat:random":  {},
		},
	}, nil
}

func main() {
	// Create Centrifuge node
	node, err := centrifuge.New(centrifuge.Config{
		LogLevel: centrifuge.LogLevelInfo,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set up event handlers
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		return authHandler(ctx, e)
	})

	// Start Centrifuge node
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}

	// Setup HTTP routes
	http.HandleFunc("/api/login", corsMiddleware(loginHandler))
	http.HandleFunc("/api/user", corsMiddleware(authMiddleware(userInfoHandler)))
	
	// Health check
	http.HandleFunc("/api/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))

	// WebSocket handler for Centrifuge
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize:     1024,
		UseWriteBufferPool: true,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo
		},
	})
	http.Handle("/connection/websocket", wsHandler)

	log.Println("Server starting on :3001")
	log.Println("WebSocket endpoint: ws://localhost:3001/connection/websocket")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
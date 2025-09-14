package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// JWT secret key - in production, this should be from environment variable
var jwtSecret = []byte("your-secret-key")

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Create Centrifuge node
	node, err := centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelInfo,
		LogHandler: handleLog,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set connection handler
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		// Validate JWT token
		token := e.Token
		if token == "" {
			log.Printf("No token provided for connection")
			return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}

		claims := &jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !parsedToken.Valid {
			log.Printf("Invalid token: %v", err)
			return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}

		userID := (*claims)["user_id"].(string)
		username := (*claims)["username"].(string)

		log.Printf("User %s connected successfully", username)

		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID: userID,
				Info:   []byte(fmt.Sprintf(`{"username": "%s"}`, username)),
			},
			Data: []byte(fmt.Sprintf(`{"user_id": "%s", "username": "%s"}`, userID, username)),
		}, nil
	})

	// Set client connect handler  
	node.OnConnect(func(client *centrifuge.Client) {
		// Set up client event handlers
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			// Allow subscription to any topic for this example
			// In production, you might want to add authorization logic here
			log.Printf("User %s subscribed to topic %s", client.UserID(), e.Channel)
			cb(centrifuge.SubscribeReply{}, nil)
		})
	})

	// Start Centrifuge node
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}

	// Create HTTP router
	router := mux.NewRouter()

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// API routes
	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/topics/{topic}/messages", authMiddleware(publishMessageHandler(node))).Methods("POST")

	// Centrifuge WebSocket endpoint
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{})
	router.HandleFunc("/connection/websocket", wsHandler.ServeHTTP)

	// Serve static files (frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/public/")))

	// Apply CORS
	handler := c.Handler(router)

	// Start HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Server starting on :8080")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	node.Shutdown(ctx)
	log.Println("Server stopped")
}

func handleLog(entry centrifuge.LogEntry) {
	log.Printf("[%s] %s", strings.ToUpper(string(entry.Level)), entry.Message)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple authentication - in production, validate against database
	if loginReq.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  loginReq.Username,
		"username": loginReq.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: tokenString,
		User: User{
			ID:       loginReq.Username,
			Username: loginReq.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func publishMessageHandler(node *centrifuge.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topic := vars["topic"]

		var message struct {
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user info from context (set by auth middleware)
		username := r.Context().Value("username").(string)

		// Create message
		msg := Message{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Topic:     topic,
			Username:  username,
			Content:   message.Content,
			Timestamp: time.Now(),
		}

		// Publish to Centrifuge
		data, _ := json.Marshal(msg)
		_, err := node.Publish(topic, data)
		if err != nil {
			http.Error(w, "Error publishing message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for OPTIONS requests
		if r.Method == "OPTIONS" {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", (*claims)["user_id"])
		ctx = context.WithValue(ctx, "username", (*claims)["username"])
		next(w, r.WithContext(ctx))
	}
}
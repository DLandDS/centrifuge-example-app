package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/centrifugal/centrifuge"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func TestMain(m *testing.M) {
	// Set test environment variables
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("PORT", "8080")
	
	// Run tests
	code := m.Run()
	
	// Clean up
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("PORT")
	
	os.Exit(code)
}

func createTestServer() *httptest.Server {
	// Create Centrifuge node for testing
	node, _ := centrifuge.New(centrifuge.Config{
		LogLevel: centrifuge.LogLevelError, // Reduce log noise in tests
	})
	
	// Start node
	node.Run()
	
	// Create router
	router := mux.NewRouter()
	
	// Add CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})
	
	// Add routes
	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/topics/{topic}/messages", authMiddleware(publishMessageHandler(node))).Methods("POST")
	
	// Apply CORS
	handler := c.Handler(router)
	
	return httptest.NewServer(handler)
}

func TestLoginHandler(t *testing.T) {
	server := createTestServer()
	defer server.Close()
	
	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{
			name:           "Valid login",
			username:       "testuser",
			password:       "testpass",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty username",
			username:       "",
			password:       "testpass",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid login without password",
			username:       "testuser",
			password:       "",
			expectedStatus: http.StatusOK,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginReq := LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}
			
			jsonData, _ := json.Marshal(loginReq)
			resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
			
			if resp.StatusCode == http.StatusOK {
				var loginResp LoginResponse
				if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if loginResp.Token == "" {
					t.Error("Expected token in response")
				}
				
				if loginResp.User.Username != tt.username {
					t.Errorf("Expected username %s, got %s", tt.username, loginResp.User.Username)
				}
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	server := createTestServer()
	defer server.Close()
	
	// First, login to get a token
	loginReq := LoginRequest{Username: "testuser", Password: "testpass"}
	jsonData, _ := json.Marshal(loginReq)
	resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()
	
	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	token := loginResp.Token
	
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			token:          "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			token:          "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageReq := map[string]string{"content": "Test message"}
			jsonData, _ := json.Marshal(messageReq)
			
			req, _ := http.NewRequest("POST", server.URL+"/api/topics/general/messages", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test default values
	defaultSecret := getEnv("NONEXISTENT_VAR", "default-value")
	if defaultSecret != "default-value" {
		t.Errorf("Expected default value, got %s", defaultSecret)
	}
	
	// Test existing environment variable
	os.Setenv("TEST_VAR", "test-value")
	testValue := getEnv("TEST_VAR", "default-value")
	if testValue != "test-value" {
		t.Errorf("Expected test-value, got %s", testValue)
	}
	os.Unsetenv("TEST_VAR")
}
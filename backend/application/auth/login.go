package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"forum/application"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"` // Can be nickname or email
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := AuthenticateUser(req)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func AuthenticateUser(req LoginRequest) (string, error) {
	var storedPassword string
	db := application.GetDB()
	row := db.QueryRow("SELECT password FROM User WHERE nickname = ? OR email = ?", req.Username, req.Username)
	if err := row.Scan(&storedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Placeholder token generation
	token := "some-jwt-token"
	return token, nil
}

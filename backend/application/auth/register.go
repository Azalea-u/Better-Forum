package auth

import (
	"encoding/json"
	"fmt"
	"forum/application"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := RegisterUser(req); err != nil {
		log.Printf("Registration failed: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func RegisterUser(req RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	db := application.GetDB()

	_, err = db.Exec(
		"INSERT INTO User (nickname, age, gender, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.Nickname, req.Age, req.Gender, req.FirstName, req.LastName, req.Email, passwordHash,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username  string `json:"username"` // Can be nickname or email
	Password  string `json:"password"`
	SessionID string `json:"session_id"` // Unique session ID
	IP        string `json:"ip_address"` // User's IP address
	UserAgent string `json:"user_agent"` // User agent (browser/device)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := AuthenticateUser(req, db)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Update or insert into OnlineStatus with hashed IP and User Agent
	if err := UpdateOnlineStatus(userID, req.SessionID, req.IP, req.UserAgent, true, db); err != nil {
		log.Printf("Failed to update online status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func AuthenticateUser(req LoginRequest, db *sql.DB) (int, error) {
	var storedPassword string
	var userID int
	row := db.QueryRow("SELECT id, password FROM User WHERE nickname = ? OR email = ?", req.Username, req.Username)
	if err := row.Scan(&userID, &storedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
		return 0, errors.New("invalid credentials")
	}

	return userID, nil
}

func UpdateOnlineStatus(userID int, sessionID, ipAddress, userAgent string, isOnline bool, db *sql.DB) error {
	currentTime := time.Now()

	// Hash IP address and User Agent using bcrypt
	hashedIP, err := bcrypt.GenerateFromPassword([]byte(ipAddress), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash IP address: %w", err)
	}

	hashedUserAgent, err := bcrypt.GenerateFromPassword([]byte(userAgent), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user agent: %w", err)
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM OnlineStatus WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check online status: %w", err)
	}

	if exists {
		query := "UPDATE OnlineStatus SET session_id = ?, ip_address = ?, user_agent = ?, is_online = ?, last_active = ?, login_time = ? WHERE user_id = ?"
		if !isOnline {
			query = "UPDATE OnlineStatus SET is_online = ?, logout_time = ?, last_active = ? WHERE user_id = ?"
		}
		_, err = db.Exec(query, sessionID, hashedIP, hashedUserAgent, isOnline, currentTime, currentTime, userID)
	} else {
		_, err = db.Exec("INSERT INTO OnlineStatus (user_id, session_id, ip_address, user_agent, is_online, last_active, login_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
			userID, sessionID, hashedIP, hashedUserAgent, isOnline, currentTime, currentTime)
	}

	if err != nil {
		return fmt.Errorf("failed to update online status: %w", err)
	}

	return nil
}

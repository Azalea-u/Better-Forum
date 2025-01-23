package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type LogoutRequest struct {
	SessionID string `json:"session_id"` // Unique session ID
	UserID    int    `json:"user_id"`    // User ID
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := UpdateOnlineStatus(req.UserID, req.SessionID, "", "", false, db); err != nil {
		log.Printf("Failed to update online status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

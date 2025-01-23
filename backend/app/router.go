package application

import (
	"forum/app/auth"
	"net/http"
)

func Router() *http.ServeMux {
	db := GetDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket)

	// Authentication routes
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r, db)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterHandler(w, r, db)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.LogoutHandler(w, r, db)
	})

	return mux
}

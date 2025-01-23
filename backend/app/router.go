package application

import (
	"database/sql"
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

func middleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies, err := r.Cookie("session_id")
		if err != nil {
			// redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// check if the session id is valid
		q := db.QueryRow("SELECT id FROM OnlineStatus WHERE session_id = ?", cookies.Value)
		var id int
		if err := q.Scan(&id); err != nil {
			// redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

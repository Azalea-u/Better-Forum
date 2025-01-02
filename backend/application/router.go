package application

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket)
	return mux
}
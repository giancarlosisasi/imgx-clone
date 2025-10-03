package app

import (
	"encoding/json"
	"net/http"
)

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/v1/images/upload", app.uploadImageHandler)
	mux.HandleFunc("GET /api/v1/images/{filename}", app.getImageHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"message": "server is running",
	}

	json.NewEncoder(w).Encode(response)
}

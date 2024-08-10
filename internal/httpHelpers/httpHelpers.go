package httpHelpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleErrorResponse(w http.ResponseWriter, message string, status int) {
	response := map[string]string{"error": message}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleSuccessResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

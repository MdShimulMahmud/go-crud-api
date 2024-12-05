package utils

import (
	"encoding/json"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})

}

func SuccessHandler(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
	}
}

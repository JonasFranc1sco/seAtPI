package responses

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, statusCode int, data any) {
	response := SuccessResponse{
		Data: data,
	}

	JSON(w, statusCode, response)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error: message,
	}

	JSON(w, statusCode, response)
}

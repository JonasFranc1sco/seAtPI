package handlers

import (
	"app/app/internal/responses"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"service": "inventory-api",
	}

	responses.Success(w, http.StatusOK, data)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"version": "0.1.0",
	}

	responses.Success(w, http.StatusOK, data)
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"name":        "Inventory-api",
		"description": "API de inventário e conformidade de equipamentos",
	}

	responses.Success(w, http.StatusOK, data)
}

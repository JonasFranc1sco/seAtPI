package handlers

import (
	"app/app/internal/responses"
	"app/app/internal/services"
	"net/http"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

func (h *InventoryHandler) ImportInventoryCSVHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "arquivo CSV não enviado no campo file")
		return
	}
	defer file.Close()

	if header.Size == 0 {
		responses.Error(w, http.StatusBadRequest, "arquivo CSV vazio")
		return
	}

	result, err := h.service.ImportInventoryCSV(file)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(w, http.StatusOK, result)
}

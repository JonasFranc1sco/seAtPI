package handlers

import (
	"app/app/internal/responses"
	"app/app/internal/services"
	"net/http"
)

type AntivirusHandler struct {
	service *services.AntivirusService
}

func NewAntivirusHandler(service *services.AntivirusService) *AntivirusHandler {
	return &AntivirusHandler{
		service: service,
	}
}

func (h *AntivirusHandler) ImportTrendCSVHandler(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.ImportTrendCSV(file)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(w, http.StatusOK, result)
}

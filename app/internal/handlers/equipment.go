package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/app/internal/models"
	"app/app/internal/responses"
	"app/app/internal/services"

	"github.com/gorilla/mux"
)

type EquipmentHandler struct {
	service *services.EquipmentService
}

func (h *EquipmentHandler) CreateEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	var equipment models.Equipment

	err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	createdEquipment, err := h.service.CreateEquipment(equipment)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAssertTag):
			responses.Error(w, http.StatusBadRequest, "tombo inválido")
		case errors.Is(err, services.ErrRequiredField):
			responses.Error(w, http.StatusBadRequest, "campos obrigatórios não informados")
		case errors.Is(err, services.ErrEquipmentAlreadyExists):
			responses.Error(w, http.StatusConflict, "equipamento já cadastrado")
		default:
			responses.Error(w, http.StatusInternalServerError, "erro interno ao cadastrar equipamento")
		}

		return
	}

	responses.Success(w, http.StatusCreated, createdEquipment)
}

func (h *EquipmentHandler) UpdateEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetTag := vars["tombo"]

	var equipment models.Equipment

	err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	updatedEquipment, err := h.service.UpdateEquipment(assetTag, equipment)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAssertTag):
			responses.Error(w, http.StatusBadRequest, "tombo inválido")
		case errors.Is(err, services.ErrRequiredField):
			responses.Error(w, http.StatusBadRequest, "campos obrigatórios não informados")
		case errors.Is(err, services.ErrEquipmentNotFound):
			responses.Error(w, http.StatusNotFound, "equipamento não encontrado")
		default:
			responses.Error(w, http.StatusInternalServerError, "erro interno ao atualizar equipamento")
		}

		return
	}

	responses.Success(w, http.StatusOK, updatedEquipment)
}

func (h *EquipmentHandler) DeleteEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetTag := vars["tombo"]

	err := h.service.DeleteEquipment(assetTag)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAssertTag):
			responses.Error(w, http.StatusBadRequest, "tombo inválido")
		case errors.Is(err, services.ErrEquipmentNotFound):
			responses.Error(w, http.StatusNotFound, "equipamento não encontrado")
		default:
			responses.Error(w, http.StatusInternalServerError, "erro interno ao remover equipamento")
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewEquipmentHandler(service *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{
		service: service,
	}
}

func (h *EquipmentHandler) ListEquipmentsHandler(w http.ResponseWriter, r *http.Request) {
	equipments, err := h.service.ListEquipments()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "erro interno ao listar equipamentos")
		return
	}

	responses.Success(w, http.StatusOK, equipments)
}

func (h *EquipmentHandler) GetEquipmentByAssetTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetTag := vars["tombo"]

	equipment, err := h.service.GetEquipmentByAssetTag(assetTag)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidAssertTag):
			responses.Error(w, http.StatusBadRequest, "tombo inválido")
		case errors.Is(err, services.ErrEquipmentNotFound):
			responses.Error(w, http.StatusNotFound, "equipamento não encontrado")
		default:
			responses.Error(w, http.StatusInternalServerError, "erro interno ao buscar equipamento")
		}

		return
	}

	responses.Success(w, http.StatusOK, equipment)
}

func (h *EquipmentHandler) GetEquipmentBySerialNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serialNumber := vars["serial"]

	equipment, err := h.service.GetEquipmentBySerialNumber(serialNumber)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidSerialNumber):
			responses.Error(w, http.StatusBadRequest, "serial number inválido")
		case errors.Is(err, services.ErrEquipmentNotFound):
			responses.Error(w, http.StatusNotFound, "equipamento não encontrado")
		default:
			responses.Error(w, http.StatusInternalServerError, "erro interno ao buscar equipamento")
		}

		return
	}

	responses.Success(w, http.StatusOK, equipment)
}

func (h *EquipmentHandler) ImportPatrimonyCSVHandler(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.ImportPatrimonyCSV(file)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(w, http.StatusOK, result)
}

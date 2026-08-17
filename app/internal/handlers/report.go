package handlers

import (
	"app/app/internal/responses"
	"app/app/internal/services"
	"encoding/csv"
	"log"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

func (h *ReportHandler) EquipmentsWithoutAntivirusHandler(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetEquipmentWithoutAntivirus()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "erro interno ao gerar relatório de equipamentos sem antivírus")
		return
	}

	responses.Success(w, http.StatusOK, report)
}

func (h *ReportHandler) EquipmentWithAntivirusProblemsHandler(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetEquipmentWithAntivirusProblems()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "erro interno ao gerar relatório do antivírus com problema")
		return
	}

	responses.Success(w, http.StatusOK, report)
}

func (h *ReportHandler) InventoryEquipmentAntivirusUnifiedHandler(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetInventoryEquipmentAntivirusUnified()
	if err != nil {
		log.Printf("Erro no handler: %v", err)
		responses.Error(w, http.StatusInternalServerError, "erro interno ao gerar relatório de unificação de tabelas")
		return
	}

	// Definição cabeçalhos para download do CSV unificado
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"relatorio_unificado.csv\"")

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	header := []string{"Asset Tag", "Serial Number", "Warranty End Date", "Antivirus Recommended Actions", "Antivirus Sensor Connectivity", "Antivirus Last Agent Status Reported", "Antivirus AntiMalware", "Operating System", "IP Address", "Status", "Logged In Users", "Organization Unit"}
	if err := writer.Write(header); err != nil {
		log.Printf("Erro ao escrever cabeçalho: %v", err)
		return
	}

	for _, item := range report.Equipments {
		row := []string{
			item.Equipment.AssetTag,
			item.Equipment.SerialNumber,
			item.Equipment.WarrantyEndDate,
			item.Antivirus.RecommendedActions,
			item.Antivirus.SensorConnectivity,
			item.Antivirus.LastAgentStatusReported,
			item.Antivirus.AntiMalware,
			item.Inventory.OS,
			item.Inventory.IpAddress,
			item.Inventory.Status,
			item.Inventory.LoggedInUsers,
			item.Inventory.OrganizationUnit,
		}

		if err := writer.Write(row); err != nil {
			log.Printf("Erro ao escrever linha: %v", err)
			return
		}
	}

	writer.Flush()
}

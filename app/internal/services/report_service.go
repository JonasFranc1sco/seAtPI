package services

import (
	"app/app/internal/models"
	"app/app/internal/repositories"
	"strings"
)

type ReportService struct {
	repository *repositories.ReportRepository
}

func buildAntivirusIssues(endpoint models.AntivirusEndpoint) []string {
	issues := []string{}

	sensorConnectivity := strings.ToUpper(strings.TrimSpace(endpoint.SensorConnectivity))
	endpointSensor := strings.ToUpper(strings.TrimSpace(endpoint.EndpointSensor))
	antiMalware := strings.ToUpper(strings.TrimSpace(endpoint.AntiMalware))
	recommendedActions := strings.TrimSpace(endpoint.RecommendedActions)

	if sensorConnectivity != "CONNECTED" {
		issues = append(issues, "sensor desconectado")
	}

	if !strings.Contains(endpointSensor, "ENABLED: RUNNING") {
		issues = append(issues, "endpoint sensor não está em execução")
	}

	if strings.Contains(antiMalware, "OUTDATED") {
		issues = append(issues, "anti-malware desatualizado")
	}

	if strings.Contains(antiMalware, "NOT OPTIMIZED") {
		issues = append(issues, "anti-malware não otimizado")
	}

	if strings.Contains(antiMalware, "DISABLED") {
		issues = append(issues, "anti-malware desabilitado")
	}

	if recommendedActions != "" {
		issues = append(issues, "existem ações recomendadas pelo Trend")
	}

	return issues
}

func NewReportService(repository *repositories.ReportRepository) *ReportService {
	return &ReportService{
		repository: repository,
	}
}

func (s *ReportService) GetEquipmentWithoutAntivirus() (models.EquipmentsWithoutAntivirusReport, error) {
	equipments, err := s.repository.FindEquipmentsWithoutAntivirus()
	if err != nil {
		return models.EquipmentsWithoutAntivirusReport{}, err
	}

	report := models.EquipmentsWithoutAntivirusReport{
		Total:      len(equipments),
		Equipments: equipments,
	}

	return report, nil
}

func (s *ReportService) GetEquipmentWithAntivirusProblems() (models.EquipmentsWithAntivirusProblemsReport, error) {
	items, err := s.repository.FindEquipmentsWithAntivirusProblems()
	if err != nil {
		return models.EquipmentsWithAntivirusProblemsReport{}, err
	}

	for index := range items {
		items[index].Issues = buildAntivirusIssues(items[index].Antivirus)
	}

	report := models.EquipmentsWithAntivirusProblemsReport{
		Total:      len(items),
		Equipments: items,
	}

	return report, nil
}

func (s *ReportService) GetInventoryEquipmentAntivirusUnified() (models.UnifiedFields, error) {
	items, err := s.repository.FindInventoryEquipmentAntivirusUnified()
	if err != nil {
		return models.UnifiedFields{}, err
	}

	report := models.UnifiedFields{
		Total: len(items),
		Equipments: items,
	}

	return report, nil
}

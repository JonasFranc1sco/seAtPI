package services

import (
	"app/app/internal/models"
	"app/app/internal/repositories"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type AntivirusService struct {
	repository *repositories.AntivirusRepository
}

func NewAntivirusService(repository *repositories.AntivirusRepository) *AntivirusService {
	return &AntivirusService{
		repository: repository,
	}
}

func (s *AntivirusService) ImportTrendCSV(reader io.Reader) (models.ImportResult, error) {
	reFJ := regexp.MustCompile(`[a-zA-Z]`)
	reZeros := regexp.MustCompile(`^0+`)
	data, err := io.ReadAll(reader)
	if err != nil {
		return models.ImportResult{}, fmt.Errorf("erro ao ler arquivo CSV do Trend: %w", err)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return models.ImportResult{}, errors.New("CSV vazio")
	}

	delimiter := detectCSVDelimiter(data)

	csvReader := csv.NewReader(bytes.NewReader(data))
	csvReader.Comma = delimiter
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1

	headers, err := csvReader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return models.ImportResult{}, errors.New("CSV vazio")
		}

		return models.ImportResult{}, fmt.Errorf("erro ao ler cabeçalho do CSV do Trend: %w", err)
	}

	columns := mapCSVHeaders(headers)

	requiredColumns := []string{
		"ENDPOINT NAME",
		"RECOMMENDED ACTIONS",
		"ENDPOINT SENSOR",
		"OS NAME",
	}

	for _, column := range requiredColumns {
		if _, exists := columns[column]; !exists {
			return models.ImportResult{}, fmt.Errorf("coluna obrigatória não encontrada: %s", column)
		}
	}

	result := models.ImportResult{
		Delimiter: delimiterToString(delimiter),
	}

	lineNumber := 1

	for {
		record, err := csvReader.Read()
		lineNumber++

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:    lineNumber,
				Message: fmt.Sprintf("erro ao ler linha: %v", err),
			})
			continue
		}

		result.TotalRows++

		endpoint := models.AntivirusEndpoint{
			EndpointName:            getCSVValue(record, columns, "ENDPOINT NAME"),
			RecommendedActions:      getCSVValue(record, columns, "RECOMMENDED ACTIONS"),
			EndpointSensor:          getCSVValue(record, columns, "ENDPOINT SENSOR"),
			OSName:                  getCSVValue(record, columns, "OS NAME"),
			SensorConnectivity:      getCSVValue(record, columns, "SENSOR CONNECTIVITY"),
			LastAgentStatusReported: getCSVValue(record, columns, "LAST AGENT STATUS REPORTED"),
			AntiMalware:             getCSVValue(record, columns, "ANTI-MALWARE"),
		}

		endpoint.EndpointName = strings.TrimSpace(strings.ToUpper(endpoint.EndpointName))
		endpoint.EndpointName = reFJ.ReplaceAllString(endpoint.EndpointName, "")
		endpoint.EndpointName = reZeros.ReplaceAllString(endpoint.EndpointName, "")

		endpoint.RecommendedActions = strings.TrimSpace(endpoint.RecommendedActions)
		endpoint.EndpointSensor = strings.TrimSpace(endpoint.EndpointSensor)
		endpoint.OSName = strings.TrimSpace(endpoint.OSName)
		endpoint.SensorConnectivity = strings.TrimSpace(endpoint.SensorConnectivity)
		endpoint.LastAgentStatusReported = strings.TrimSpace(endpoint.LastAgentStatusReported)
		endpoint.AntiMalware = strings.TrimSpace(endpoint.AntiMalware)

		if endpoint.EndpointName == "" {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:    lineNumber,
				Message: "linha ignorada: Endpoint name vazio",
			})
			continue
		}

		_, alreadyExists, err := s.repository.FindByEndpointName(endpoint.EndpointName)
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: endpoint.EndpointName,
				Message:  fmt.Sprintf("erro ao verificar endpoint existente: %v", err),
			})
			continue
		}

		if err := s.repository.Upsert(endpoint); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: endpoint.EndpointName,
				Message:  err.Error(),
			})
			continue
		}

		result.Imported++

		if alreadyExists {
			result.Updated++
		} else {
			result.Created++
		}
	}

	return result, nil
}

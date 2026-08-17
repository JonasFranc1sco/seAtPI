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

type InventoryService struct {
	repository *repositories.InventoryRepository
}

func NewInventoryService(repository *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{
		repository: repository,
	}
}

func (s *InventoryService) ImportInventoryCSV(reader io.Reader) (models.ImportResult, error) {
	reFJ := regexp.MustCompile(`[a-zA-Z]`)
	reZeros := regexp.MustCompile(`^0+`)

	data, err := io.ReadAll(reader)
	if err != nil {
		return models.ImportResult{}, fmt.Errorf("erro ao ler arquivo CSV do inventário: %w", err)
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

		return models.ImportResult{}, fmt.Errorf("erro ao ler cabeçalho do CSV do inventário: %w", err)
	}

	columns := mapCSVHeaders(headers)

	requiredColumns := []string{
		"HOST NAME",
		"STATUS",
		"IP ADDRESS",
		"INVENTORY STATUS",
		"OS",
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

		inventory := models.Inventory{
			HostName:         getCSVValue(record, columns, "HOST NAME"),
			Status:           getCSVValue(record, columns, "STATUS"),
			IpAddress:        getCSVValue(record, columns, "IP ADDRESS"),
			InventoryStatus:  getCSVValue(record, columns, "INVENTORY STATUS"),
			OS:               getCSVValue(record, columns, "OS"),
			DeviceCategory:   getCSVValue(record, columns, "DEVICE CATEGORY"),
			LoggedInUsers:    getCSVValue(record, columns, "LOGGED IN USERS"),
			OrganizationUnit: getCSVValue(record, columns, "ORGANIZATION UNIT"),
		}

		inventory.HostName = strings.TrimSpace(inventory.HostName)
		inventory.HostName = reFJ.ReplaceAllString(inventory.HostName, "")
		inventory.HostName = reZeros.ReplaceAllString(inventory.HostName, "")

		inventory.Status = strings.TrimSpace(inventory.Status)
		inventory.IpAddress = strings.TrimSpace(inventory.IpAddress)
		inventory.InventoryStatus = strings.TrimSpace(inventory.InventoryStatus)
		inventory.OS = strings.TrimSpace(inventory.OS)
		inventory.DeviceCategory = strings.TrimSpace(inventory.DeviceCategory)
		inventory.LoggedInUsers = strings.TrimSpace(inventory.LoggedInUsers)
		inventory.OrganizationUnit = strings.TrimSpace(inventory.OrganizationUnit)

		if inventory.HostName == "" {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:    lineNumber,
				Message: "Linha ignorada, Host Name vazio",
			})
			continue
		}

		_, alreadyExists, err := s.repository.FindByHostName(inventory.HostName)
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: inventory.HostName,
				Message:  fmt.Sprintf("erro ao verificar host name existente: %v", err),
			})
			continue
		}

		if err := s.repository.Upsert(inventory); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: inventory.HostName,
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

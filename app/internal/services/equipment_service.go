package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"app/app/internal/models"
	"app/app/internal/repositories"
)

var ErrEquipmentNotFound = errors.New("equipamento não encontrado")
var ErrInvalidAssertTag = errors.New("tombo inválido")
var ErrInvalidSerialNumber = errors.New("serial number inválido")
var ErrEquipmentAlreadyExists = errors.New("equipamento já cadastrado")
var ErrRequiredField = errors.New("campo obrigatório não informado")

type EquipmentService struct {
	repository *repositories.EquipmentRepository
}

// Detect CSV delimiter type (comma or semicolon)
func detectCSVDelimiter(data []byte) rune {
	firstLine := string(data)

	if index := strings.IndexAny(firstLine, "\r\n"); index >= 0 {
		firstLine = firstLine[:index]
	}
	commaCount := strings.Count(firstLine, ",")
	semicolonCount := strings.Count(firstLine, ";")

	if semicolonCount > commaCount {
		return ';'
	}
	return ','
}

func delimiterToString(delimiter rune) string {
	if delimiter == ';' {
		return ";"
	}

	return ","
}

// Convert CSV header to map
func mapCSVHeaders(headers []string) map[string]int {
	columns := make(map[string]int)

	for index, header := range headers {
		normalizedHeader := strings.TrimSpace(header)
		normalizedHeader = strings.TrimPrefix(normalizedHeader, "\ufeff")
		normalizedHeader = strings.ToUpper(normalizedHeader)
		columns[normalizedHeader] = index
	}

	return columns
}

// Get value based on name columns
func getCSVValue(record []string, columns map[string]int, columnName string) string {
	index, exists := columns[columnName]
	if !exists {
		return ""
	}

	if index >= len(record) {
		return ""
	}

	return record[index]
}

func (s *EquipmentService) CreateEquipment(equipment models.Equipment) (models.Equipment, error) {
	equipment.AssetTag = strings.TrimSpace(strings.ToUpper(equipment.AssetTag))
	equipment.Building = strings.TrimSpace(equipment.Building)
	equipment.Room = strings.TrimSpace(equipment.Room)
	equipment.BudgetUnit = strings.TrimSpace(equipment.BudgetUnit)
	equipment.SerialNumber = strings.TrimSpace(strings.ToUpper(equipment.SerialNumber))
	equipment.Description = strings.TrimSpace(equipment.Description)
	equipment.AssetType = strings.TrimSpace(equipment.AssetType)
	equipment.WarrantyEndDate = strings.TrimSpace(equipment.WarrantyEndDate)

	if equipment.AssetTag == "" {
		return models.Equipment{}, ErrInvalidAssertTag
	}

	if equipment.Building == "" ||
		equipment.Room == "" ||
		equipment.BudgetUnit == "" ||
		equipment.Description == "" ||
		equipment.AssetType == "" {
		return models.Equipment{}, ErrRequiredField
	}

	_, found, err := s.repository.FindByAssetTag(equipment.AssetTag)
	if err != nil {
		return models.Equipment{}, err
	}

	if found {
		return models.Equipment{}, ErrEquipmentAlreadyExists
	}

	createdEquipment, err := s.repository.Create(equipment)
	if err != nil {
		return models.Equipment{}, err
	}

	return createdEquipment, nil
}

func NewEquipmentService(repository *repositories.EquipmentRepository) *EquipmentService {
	return &EquipmentService{
		repository: repository,
	}
}

func (s *EquipmentService) UpdateEquipment(assetTag string, equipment models.Equipment) (models.Equipment, error) {
	normalizedAssetTag := strings.TrimSpace(strings.ToUpper(assetTag))

	if normalizedAssetTag == "" {
		return models.Equipment{}, ErrInvalidAssertTag
	}

	equipment.AssetTag = normalizedAssetTag
	equipment.Building = strings.TrimSpace(equipment.Building)
	equipment.Room = strings.TrimSpace(equipment.Room)
	equipment.BudgetUnit = strings.TrimSpace(equipment.BudgetUnit)
	equipment.SerialNumber = strings.TrimSpace(strings.ToUpper(equipment.SerialNumber))
	equipment.Description = strings.TrimSpace(equipment.Description)
	equipment.AssetType = strings.TrimSpace(equipment.AssetType)
	equipment.WarrantyEndDate = strings.TrimSpace(equipment.WarrantyEndDate)

	if equipment.Building == "" ||
		equipment.Room == "" ||
		equipment.BudgetUnit == "" ||
		equipment.Description == "" ||
		equipment.AssetType == "" {
		return models.Equipment{}, ErrRequiredField
	}

	_, found, err := s.repository.FindByAssetTag(normalizedAssetTag)
	if err != nil {
		return models.Equipment{}, err
	}

	if !found {
		return models.Equipment{}, ErrEquipmentNotFound
	}

	updatedEquipment, updated, err := s.repository.Update(normalizedAssetTag, equipment)
	if err != nil {
		return models.Equipment{}, err
	}

	if !updated {
		return models.Equipment{}, ErrEquipmentNotFound
	}

	return updatedEquipment, nil
}

func (s *EquipmentService) DeleteEquipment(assetTag string) error {
	normalizedAssetTag := strings.TrimSpace(strings.ToUpper(assetTag))

	if normalizedAssetTag == "" {
		return ErrInvalidAssertTag
	}

	deleted, err := s.repository.Delete(normalizedAssetTag)
	if err != nil {
		return err
	}

	if !deleted {
		return ErrEquipmentNotFound
	}

	return nil
}

func (s *EquipmentService) ListEquipments() ([]models.Equipment, error) {
	return s.repository.FindAll()
}

func (s *EquipmentService) GetEquipmentByAssetTag(assetTag string) (models.Equipment, error) {
	normalizedAssetTag := strings.TrimSpace(strings.ToUpper(assetTag))

	if normalizedAssetTag == "" {
		return models.Equipment{}, ErrInvalidAssertTag
	}

	equipment, found, err := s.repository.FindByAssetTag(normalizedAssetTag)
	if err != nil {
		return models.Equipment{}, err
	}

	if !found {
		return models.Equipment{}, ErrEquipmentNotFound
	}

	return equipment, nil
}

func (s *EquipmentService) GetEquipmentBySerialNumber(serialNumber string) (models.Equipment, error) {
	normalizedSerialNumber := strings.TrimSpace(strings.ToUpper(serialNumber))

	if normalizedSerialNumber == "" {
		return models.Equipment{}, ErrInvalidSerialNumber
	}
	equipment, found, err := s.repository.FindBySerialNumber(normalizedSerialNumber)
	if err != nil {
		return models.Equipment{}, err
	}

	if !found {
		return models.Equipment{}, ErrEquipmentNotFound
	}

	return equipment, nil

}

func (s *EquipmentService) ImportPatrimonyCSV(reader io.Reader) (models.ImportResult, error) {
	reFJ := regexp.MustCompile(`[a-zA-Z]`)
	reZeros := regexp.MustCompile(`^0+`)
	data, err := io.ReadAll(reader)
	if err != nil {
		return models.ImportResult{}, fmt.Errorf("erro ao ler arquivo CSV: %w", err)
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

		return models.ImportResult{}, fmt.Errorf("erro ao ler cabeçalho do CSV: %w", err)
	}

	columns := mapCSVHeaders(headers)

	requiredColumns := []string{
		"DEPREDIO",
		"DESALA",
		"UO",
		"TOMBO",
		"DEBEM",
		"DESUBGRUPOBEM",
		"DTGARANTIA",
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

		equipment := models.Equipment{
			Building:        getCSVValue(record, columns, "DEPREDIO"),
			Room:            getCSVValue(record, columns, "DESALA"),
			BudgetUnit:      getCSVValue(record, columns, "UO"),
			AssetTag:        getCSVValue(record, columns, "TOMBO"),
			SerialNumber:    getCSVValue(record, columns, "VOLUME"),
			Description:     getCSVValue(record, columns, "DEBEM"),
			AssetType:       getCSVValue(record, columns, "DESUBGRUPOBEM"),
			WarrantyEndDate: getCSVValue(record, columns, "DTGARANTIA"),
		}

		equipment.AssetTag = strings.TrimSpace(strings.ToUpper(equipment.AssetTag))
		equipment.AssetTag = reFJ.ReplaceAllString(equipment.AssetTag, "")
		equipment.AssetTag = reZeros.ReplaceAllString(equipment.AssetTag, "")

		equipment.SerialNumber = strings.TrimSpace(strings.ToUpper(equipment.SerialNumber))
		equipment.Building = strings.TrimSpace(equipment.Building)
		equipment.Room = strings.TrimSpace(equipment.Room)
		equipment.BudgetUnit = strings.TrimSpace(equipment.BudgetUnit)
		equipment.Description = strings.TrimSpace(equipment.Description)
		equipment.AssetType = strings.TrimSpace(equipment.AssetType)
		equipment.WarrantyEndDate = strings.TrimSpace(equipment.WarrantyEndDate)

		if equipment.AssetTag == "" {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:    lineNumber,
				Message: "linha ingorada: TOMBO vazio",
			})
			continue
		}

		if equipment.Building == "" ||
			equipment.Room == "" ||
			equipment.BudgetUnit == "" ||
			equipment.Description == "" ||
			equipment.AssetType == "" {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: equipment.AssetTag,
				Message:  "linha ignorada: campos obrigatórios vazios",
			})
			continue
		}

		_, alreadyExists, err := s.repository.FindByAssetTag(equipment.AssetTag)
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: equipment.AssetTag,
				Message:  fmt.Sprintf("erro ao verificar equipamento existente: %v", err),
			})
			continue
		}

		if err := s.repository.Upsert(equipment); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, models.ImportError{
				Line:     lineNumber,
				AssetTag: equipment.AssetTag,
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

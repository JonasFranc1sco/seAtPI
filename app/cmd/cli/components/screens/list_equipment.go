package screens

import (
	"app/app/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"charm.land/huh/v2"
)

func ListEquipments(endpoint string) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return fmt.Errorf("erro ao buscar equipamentos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var result struct {
		Data []models.Equipment `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("erro ao decodificar equipamentos: %w", err)
	}

	if len(result.Data) == 0 {
		huh.NewNote().
			Title("Nenhum equipamento encontrado.").
			Description("A API não retornou equipamentos.").
			Run()
		return nil
	}

	options := make([]huh.Option[string], 0, len(result.Data))
	for _, eq := range result.Data {
		labels := fmt.Sprintf("%s - %s [%s]", eq.AssetTag, eq.Description, eq.SerialNumber)
		options = append(options, huh.NewOption(labels, eq.AssetTag))
	}

	var chosen string
	if err := huh.NewSelect[string]().
		Title("Equipamentos disponíveis").
		Options(options...).
		Value(&chosen).
		Run(); err != nil {
		return err
	}

	for _, eq := range result.Data {
		if eq.AssetTag == chosen {
			return huh.NewNote().
				Title(fmt.Sprintf("Equipamento %s", eq.AssetTag)).
				Description(fmt.Sprintf(
					"Descrição: %s\nSerial: %s\nPrédio: %s\nSala: %s\nUO: %s\nGarantia até: %s",
					eq.Description, eq.SerialNumber, eq.AssetType,
					eq.Building, eq.Room, eq.BudgetUnit, eq.WarrantyEndDate)).
				Run()
		}
	}
	return nil
}

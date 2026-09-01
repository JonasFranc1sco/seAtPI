package screens

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"charm.land/huh/v2"
)

func ListEquipments(baseURL string) error {
	resp, err := http.Get(baseURL + "/equipamentos")
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

	var equipments []Equipment
	if err := json.Unmarshal(body, &equipments); err != nil {
		return fmt.Errorf("erro ao decodificar equipamentos: %w", err)
	}

	for _, eq := range equipments {
		huh.NewNote().
			Title(fmt.Sprintf("Equipamento %s", eq.Tombo)).
			Description(fmt.Sprintf(
				"Nome: %s\nSerial: %s\nLocal: %s\nStatus: %s\nDescrição: %s",
				eq.Name, eq.SerialNumber, eq.Local, eq.Status, eq.Description,
			)).Run()
	}

	return nil
}

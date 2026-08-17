package repositories

import (
	"app/app/internal/models"
	"database/sql"
	"fmt"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) FindEquipmentsWithoutAntivirus() ([]models.Equipment, error) {
	query := `
	SELECT
		e.id,
		e.building,
		e.room,
		e.budget_unit,
		e.asset_tag,
		COALESCE(e.serial_number, ''),
		e.description,
		e.asset_type,
		COALESCE(e.warranty_end_date, '')
	FROM equipments e
	LEFT JOIN antivirus_endpoints a
		ON UPPER(TRIM(e.asset_tag)) = UPPER(TRIM(a.endpoint_name))
	WHERE a.endpoint_name IS NULL
	ORDER BY e.building, e.room, e.asset_tag;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar relatório de equipamentos sem antivírus: %w", err)
	}
	defer rows.Close()

	equipments := []models.Equipment{}

	for rows.Next() {
		var equipment models.Equipment

		err := rows.Scan(
			&equipment.ID,
			&equipment.Building,
			&equipment.Room,
			&equipment.BudgetUnit,
			&equipment.AssetTag,
			&equipment.SerialNumber,
			&equipment.Description,
			&equipment.AssetType,
			&equipment.WarrantyEndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler equipamento do relatório: %w", err)
		}

		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer relatório de equipamentos sem antivírus: %w", err)
	}

	return equipments, nil

}

func (r *ReportRepository) FindEquipmentsWithAntivirusProblems() ([]models.EquipmentAntivirusIssue, error) {
	query := `
	SELECT
		e.id,
		e.building,
		e.room,
		e.budget_unit,
		e.asset_tag,
		COALESCE(e.serial_number, ''),
		e.description,
		e.asset_type,
		COALESCE(e.warranty_end_date, ''),

		a.id,
		a.endpoint_name,
		COALESCE(a.recommended_actions, ''),
		COALESCE(a.endpoint_sensor, ''),
		COALESCE(a.os_name, ''),
		COALESCE(a.sensor_connectivity, ''),
		COALESCE(a.last_agent_status_reported, ''),
		COALESCE(a.anti_malware, '')
	FROM equipments e
	INNER JOIN antivirus_endpoints a
		ON UPPER(TRIM(e.asset_tag)) = UPPER(TRIM(a.endpoint_name))
	WHERE
		UPPER(COALESCE(a.sensor_connectivity, '')) <> 'CONNECTED'
		OR UPPER(COALESCE(a.endpoint_sensor, '')) NOT LIKE '%ENABLED: RUNNING%'
		OR UPPER(COALESCE(a.anti_malware, '')) LIKE '%OUTDATED%'
		OR UPPER(COALESCE(a.anti_malware, '')) LIKE '%NOT OPTIMIZED%'
		OR UPPER(COALESCE(a.anti_malware, '')) LIKE '%DISABLED%'
		OR TRIM(COALESCE(a.recommended_actions, '')) <> ''
	ORDER BY e.building, e.room, e.asset_tag;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar relatório de antivírus com problema: %w", err)
	}

	defer rows.Close()

	items := []models.EquipmentAntivirusIssue{}

	for rows.Next() {
		var item models.EquipmentAntivirusIssue

		err := rows.Scan(
			&item.Equipment.ID,
			&item.Equipment.Building,
			&item.Equipment.Room,
			&item.Equipment.BudgetUnit,
			&item.Equipment.AssetTag,
			&item.Equipment.SerialNumber,
			&item.Equipment.Description,
			&item.Equipment.AssetType,
			&item.Equipment.WarrantyEndDate,

			&item.Antivirus.ID,
			&item.Antivirus.EndpointName,
			&item.Antivirus.RecommendedActions,
			&item.Antivirus.EndpointSensor,
			&item.Antivirus.OSName,
			&item.Antivirus.SensorConnectivity,
			&item.Antivirus.LastAgentStatusReported,
			&item.Antivirus.AntiMalware,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler item do relatório de antivírus com problema: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer relatório de antivírus com problema: %w", err)
	}

	return items, nil
}

func (r *ReportRepository) FindInventoryEquipmentAntivirusUnified() ([]models.InventoryEquipmentAntivirusUnified, error) {
	query := `
		SELECT
		e.id,
		e.building,
		e.room,
		e.budget_unit,
		e.asset_tag,
		COALESCE(e.serial_number, ''),
		e.description,
		e.asset_type,
		COALESCE(e.warranty_end_date, ''),

		a.id,
		a.endpoint_name,
		COALESCE(a.recommended_actions, ''),
		COALESCE(a.endpoint_sensor, ''),
		COALESCE(a.os_name, ''),
		COALESCE(a.sensor_connectivity, ''),
		COALESCE(a.last_agent_status_reported, ''),
		COALESCE(a.anti_malware, ''),

		i.id,
		i.host_name,
		i.status,
		i.ip_address,
		i.inventory_status,
		i.os,
		i.device_category,
		COALESCE(i.logged_in_users, 'Sem usuário logado.'),
		i.organization_unit

	FROM equipments e
	INNER JOIN antivirus_endpoints a
		ON TRIM(e.asset_tag) = TRIM(a.endpoint_name)
	INNER JOIN inventory i
		ON TRIM(a.endpoint_name) = TRIM(i.host_name)
	ORDER BY e.asset_tag ASC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erro ao relacionar as três tabelas (Inventário, Antivirus, Equipamentos): %w", err)
	}

	defer rows.Close()

	items := []models.InventoryEquipmentAntivirusUnified{}

	for rows.Next() {
		var item models.InventoryEquipmentAntivirusUnified

		err := rows.Scan(
			&item.Equipment.ID,
			&item.Equipment.Building,
			&item.Equipment.Room,
			&item.Equipment.BudgetUnit,
			&item.Equipment.AssetTag,
			&item.Equipment.SerialNumber,
			&item.Equipment.Description,
			&item.Equipment.AssetType,
			&item.Equipment.WarrantyEndDate,

			&item.Antivirus.ID,
			&item.Antivirus.EndpointName,
			&item.Antivirus.RecommendedActions,
			&item.Antivirus.EndpointSensor,
			&item.Antivirus.OSName,
			&item.Antivirus.SensorConnectivity,
			&item.Antivirus.LastAgentStatusReported,
			&item.Antivirus.AntiMalware,

			&item.Inventory.ID,
			&item.Inventory.HostName,
			&item.Inventory.Status,
			&item.Inventory.IpAddress,
			&item.Inventory.InventoryStatus,
			&item.Inventory.OS,
			&item.Inventory.DeviceCategory,
			&item.Inventory.LoggedInUsers,
			&item.Inventory.OrganizationUnit,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler items do relatório de unificação de tabelas: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer relatório de unificação das tabelas: %w", err)
	}

	return items, nil
}

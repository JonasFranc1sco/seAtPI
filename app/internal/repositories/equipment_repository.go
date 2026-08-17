package repositories

import (
	"app/app/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{
		db: db,
	}
}

func (r *EquipmentRepository) Create(equipment models.Equipment) (models.Equipment, error) {
	query := `
	INSERT INTO equipments (
		building,
		room,
		budget_unit,
		asset_tag,
		serial_number,
		description,
		asset_type,
		warranty_end_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	result, err := r.db.Exec(
		query,
		equipment.Building,
		equipment.Room,
		equipment.BudgetUnit,
		equipment.AssetTag,
		equipment.SerialNumber,
		equipment.Description,
		equipment.AssetType,
		equipment.WarrantyEndDate,
	)
	if err != nil {
		return models.Equipment{}, fmt.Errorf("erro ao cadastrar equipamento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Equipment{}, fmt.Errorf("erro ao obter id do equipamento cadastrado: %w", err)
	}

	equipment.ID = id

	return equipment, nil
}

func (r *EquipmentRepository) Update(assetTag string, equipment models.Equipment) (models.Equipment, bool, error) {
	query := `
	UPDATE equipments
	SET
		building = ?,
		room = ?,
		budget_unit = ?,
		serial_number = ?,
		description = ?,
		asset_type = ?,
		warranty_end_date = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE asset_tag = ?;
	`

	result, err := r.db.Exec(
		query,
		equipment.Building,
		equipment.Room,
		equipment.BudgetUnit,
		equipment.SerialNumber,
		equipment.Description,
		equipment.AssetType,
		equipment.WarrantyEndDate,
		assetTag,
	)
	if err != nil {
		return models.Equipment{}, false, fmt.Errorf("erro ao atualizar equipamento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Equipment{}, false, fmt.Errorf("erro ao verificar atualização do equipamento: %w", err)
	}

	if rowsAffected == 0 {
		return models.Equipment{}, false, nil
	}

	updatedEquipment, found, err := r.FindByAssetTag(assetTag)
	if err != nil {
		return models.Equipment{}, false, err
	}

	if !found {
		return models.Equipment{}, false, nil
	}

	return updatedEquipment, true, nil
}

func (r *EquipmentRepository) Delete(assetTag string) (bool, error) {
	query := `
	DELETE FROM equipments
	WHERE asset_tag = ?;
	`

	result, err := r.db.Exec(query, assetTag)
	if err != nil {
		return false, fmt.Errorf("erro ao remover equipamento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar remoção do equipamento: %w", err)
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *EquipmentRepository) Upsert(equipment models.Equipment) error {
	query := `
	INSERT INTO equipments (
		building,
		room,
		budget_unit,
		asset_tag,
		serial_number,
		description,
		asset_type,
		warranty_end_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (asset_tag) DO UPDATE SET
		building = excluded.building,
		room = excluded.room,
		budget_unit = excluded.budget_unit,
		serial_number = excluded.serial_number,
		description = excluded.description,
		asset_type = excluded.asset_type,
		warranty_end_date = excluded.warranty_end_date,
		updated_at = CURRENT_TIMESTAMP;
	`

	_, err := r.db.Exec(
		query,
		equipment.Building,
		equipment.Room,
		equipment.BudgetUnit,
		equipment.AssetTag,
		equipment.SerialNumber,
		equipment.Description,
		equipment.AssetType,
		equipment.WarrantyEndDate,
	)
	if err != nil {
		return fmt.Errorf("erro ao importar equipamento %s: %w", equipment.AssetTag, err)
	}

	return nil
}

func (r *EquipmentRepository) FindAll() ([]models.Equipment, error) {
	query := `
	SELECT
		id,
		building,
		room,
		budget_unit,
		asset_tag,
		COALESCE(serial_number, ''),
		description,
		asset_type,
		COALESCE(warranty_end_date, '')
	FROM equipments
	ORDER BY id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar equipamentos: %w", err)
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
			return nil, fmt.Errorf("erro ao ler equipamento: %w", err)
		}

		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer equipamentos: %w", err)
	}

	return equipments, nil
}

func (r *EquipmentRepository) FindByAssetTag(assetTag string) (models.Equipment, bool, error) {
	query := `
	SELECT
		id,
		building,
		room,
		budget_unit,
		asset_tag,
		COALESCE(serial_number, ''),
		description,
		asset_type,
		COALESCE(warranty_end_date, '')
	FROM equipments
	WHERE asset_tag = ?;
	`

	var equipment models.Equipment

	err := r.db.QueryRow(query, assetTag).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Equipment{}, false, nil
		}

		return models.Equipment{}, false, fmt.Errorf("erro ao buscar equipamento por tombo: %w", err)
	}

	return equipment, true, nil
}

func (r *EquipmentRepository) FindBySerialNumber(serialNumber string) (models.Equipment, bool, error) {
	query := `
	SELECT
		id,
		building,
		room,
		budget_unit,
		asset_tag,
		COALESCE(serial_number, ''),
		description,
		asset_type,
		COALESCE(warranty_end_date, '')
	FROM equipments
	WHERE serial_number = ?;
	`

	var equipment models.Equipment

	err := r.db.QueryRow(query, serialNumber).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Equipment{}, false, nil
		}

		return models.Equipment{}, false, fmt.Errorf("erro ao buscar equipamento por numero serial: %w", err)
	}

	return equipment, true, nil
}

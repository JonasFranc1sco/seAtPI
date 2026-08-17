package repositories

import (
	"app/app/internal/models"
	"database/sql"
	"fmt"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) FindByHostName(hostName string) (models.Inventory, bool, error) {
	query := `
	SELECT
		id,
		host_name,
		COALESCE(status, ''),
		COALESCE(ip_address, ''),
		COALESCE(inventory_status, ''),
		COALESCE(os, ''),
		COALESCE(device_category, ''),
		COALESCE(logged_in_users, ''),
		COALESCE(organization_unit, '')
	FROM inventory
	WHERE host_name = ?;
	`
	var inventory models.Inventory

	err := r.db.QueryRow(query, hostName).Scan(
		&inventory.ID,
		&inventory.HostName,
		&inventory.Status,
		&inventory.IpAddress,
		&inventory.InventoryStatus,
		&inventory.OS,
		&inventory.DeviceCategory,
		&inventory.LoggedInUsers,
		&inventory.OrganizationUnit,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Inventory{}, false, nil
		}

		return models.Inventory{}, false, fmt.Errorf("erro ao buscar host name do inventário: %w", err)
	}
	return inventory, true, nil
}

func (r *InventoryRepository) Upsert(inventory models.Inventory) error {

	query := `
	INSERT INTO inventory (
	host_name,
	status,
	ip_address,
	inventory_status,
	os,
	device_category,
	logged_in_users,
	organization_unit
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (host_name) DO UPDATE SET
		status = excluded.status,
		ip_address = excluded.ip_address,
		os = excluded.os,
		device_category = excluded.device_category,
		logged_in_users = excluded.logged_in_users,
		organization_unit = excluded.organization_unit,
		updated_at = CURRENT_TIMESTAMP;
	`

	_, err := r.db.Exec(
		query,
		inventory.HostName,
		inventory.Status,
		inventory.IpAddress,
		inventory.InventoryStatus,
		inventory.OS,
		inventory.DeviceCategory,
		inventory.LoggedInUsers,
		inventory.OrganizationUnit,
	)

	if err != nil {
		return fmt.Errorf("erro ao importar inventário %s: %w", inventory.HostName, err)
	}

	return nil
}

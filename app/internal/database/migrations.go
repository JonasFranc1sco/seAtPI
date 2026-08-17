package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`
	CREATE TABLE IF NOT EXISTS equipments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		building TEXT NOT NULL,
		room TEXT NOT NULL,
		budget_unit TEXT NOT NULL,
		asset_tag TEXT NOT NULL UNIQUE,
		serial_number TEXT,
		description TEXT NOT NULL,
		asset_type TEXT NOT NULL,
		warranty_end_date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`,
		`
	CREATE TABLE IF NOT EXISTS antivirus_endpoints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		endpoint_name TEXT NOT NULL UNIQUE,
		recommended_actions TEXT,
		endpoint_sensor TEXT,
		os_name TEXT,
		sensor_connectivity TEXT,
		last_agent_status_reported TEXT,
		anti_malware TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`,
		`
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host_name TEXT NOT NULL UNIQUE,
		status TEXT,
		ip_address TEXT NOT NULL,
		inventory_status TEXT,
		os TEXT,
		device_category TEXT,
		logged_in_users TEXT,
		organization_unit TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("erro ao executar migrations: %w", err)
		}
	}

	return nil
}

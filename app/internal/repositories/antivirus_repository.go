package repositories

import (
	"app/app/internal/models"
	"database/sql"
	"fmt"
)

type AntivirusRepository struct {
	db *sql.DB
}

func NewAntivirusRepository(db *sql.DB) *AntivirusRepository {
	return &AntivirusRepository{
		db: db,
	}
}

func (r *AntivirusRepository) FindByEndpointName(endpointName string) (models.AntivirusEndpoint, bool, error) {
	query := `
	SELECT
		id,
		endpoint_name,
		COALESCE(recommended_actions, ''),
		COALESCE(endpoint_sensor, ''),
		COALESCE(os_name, ''),
		COALESCE(sensor_connectivity, ''),
		COALESCE(last_agent_status_reported, ''),
		COALESCE(anti_malware, '')
	FROM antivirus_endpoints
	WHERE endpoint_name = ?;
	`

	var endpoint models.AntivirusEndpoint

	err := r.db.QueryRow(query, endpointName).Scan(
		&endpoint.ID,
		&endpoint.EndpointName,
		&endpoint.RecommendedActions,
		&endpoint.EndpointSensor,
		&endpoint.OSName,
		&endpoint.SensorConnectivity,
		&endpoint.LastAgentStatusReported,
		&endpoint.AntiMalware,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.AntivirusEndpoint{}, false, nil
		}

		return models.AntivirusEndpoint{}, false, fmt.Errorf("erro ao buscar endpoint do antivírus: %w", err)
	}
	return endpoint, true, nil
}

func (r *AntivirusRepository) Upsert(endpoint models.AntivirusEndpoint) error {
	query := `
	INSERT INTO antivirus_endpoints (
		endpoint_name,
		recommended_actions,
		endpoint_sensor,
		os_name,
		sensor_connectivity,
		last_agent_status_reported,
		anti_malware
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(endpoint_name) DO UPDATE SET
		recommended_actions = excluded.recommended_actions,
		endpoint_sensor = excluded.endpoint_sensor,
		os_name = excluded.os_name,
		sensor_connectivity = excluded.sensor_connectivity,
		last_agent_status_reported = excluded.last_agent_status_reported,
		anti_malware = excluded.anti_malware,
		updated_at = CURRENT_TIMESTAMP;
	`

	_, err := r.db.Exec(
		query,
		endpoint.EndpointName,
		endpoint.RecommendedActions,
		endpoint.EndpointSensor,
		endpoint.OSName,
		endpoint.SensorConnectivity,
		endpoint.LastAgentStatusReported,
		endpoint.AntiMalware,
	)

	if err != nil {
		return fmt.Errorf("erro ao importar endpoint %s: %w", endpoint.EndpointName, err)
	}

	return nil
}

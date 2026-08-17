package models

type AntivirusEndpoint struct {
	ID                      int64  `json:"id"`
	EndpointName            string `json:"endpoint_name"`
	RecommendedActions      string `json:"recommended_actions"`
	EndpointSensor          string `json:"endpoint_sensor"`
	OSName                  string `json:"os_name"`
	SensorConnectivity      string `json:"sensor_connectivity"`
	LastAgentStatusReported string `json:"last_agent_status_reported"`
	AntiMalware             string `json:"anti_malware"`
}

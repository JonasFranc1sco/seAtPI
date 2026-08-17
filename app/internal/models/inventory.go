package models

type Inventory struct {
	ID               int64  `json:"id"`
	HostName         string `json:"host_name"`
	Status           string `json:"status"`
	IpAddress        string `json:"ip_address"`
	InventoryStatus  string `json:"inventory_status"`
	OS               string `json:"os"`
	DeviceCategory   string `json:"device_category"`
	LoggedInUsers    string `json:"logged_in_users"`
	OrganizationUnit string `json:"organization_unit"`
}

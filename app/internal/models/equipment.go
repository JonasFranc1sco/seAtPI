package models

type Equipment struct {
	ID              int64  `json:"id"`
	Building        string `json:"building"`          // Representa o DEPREDIO
	Room            string `json:"room"`              // Representa o DESALA
	BudgetUnit      string `json:"budget_unit"`       // Representa o UO
	AssetTag        string `json:"asset_tag"`         // Representa o TOMBO
	Description     string `json:"description"`       // Representa o DEBEM
	AssetType       string `json:"asset_type"`        // Representa o DESUBGRUPOBEM
	WarrantyEndDate string `json:"warranty_end_date"` // Representa o DTGARANTIAINICIO
	SerialNumber    string `json:"serial_number"`
}

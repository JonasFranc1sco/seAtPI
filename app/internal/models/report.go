package models

type EquipmentsWithoutAntivirusReport struct {
	Total      int         `json:"total"`
	Equipments []Equipment `json:"equipments"`
}

type EquipmentAntivirusIssue struct {
	Equipment Equipment         `json:"equipment"`
	Antivirus AntivirusEndpoint `json:"antivirus"`
	Issues    []string          `json:"issues"`
}

type EquipmentsWithAntivirusProblemsReport struct {
	Total      int                       `json:"total"`
	Equipments []EquipmentAntivirusIssue `json:"equipments"`
}

type InventoryEquipmentAntivirusUnified struct {
	Equipment Equipment         `json:"equipment"`
	Antivirus AntivirusEndpoint `json:"antivirus"`
	Inventory Inventory         `json:"inventory"`
}

type UnifiedFields struct {
	Total      int                                  `json:"total"`
	Equipments []InventoryEquipmentAntivirusUnified `json:"equipments"`
}

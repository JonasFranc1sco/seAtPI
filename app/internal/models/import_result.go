package models

type ImportError struct {
	Line     int    `json:"line"`
	AssetTag string `json:"asset_tag,omitempty"`
	Message  string `json:"message"`
}

type ImportResult struct {
	TotalRows int           `json:"total_rows"`
	Imported  int           `json:"imported"`
	Created   int           `json:"created"`
	Updated   int           `json:"updated"`
	Skipped   int           `json:"skipped"`
	Delimiter string        `json:"delimiter"`
	Errors    []ImportError `json:"errors,omitempty"`
}

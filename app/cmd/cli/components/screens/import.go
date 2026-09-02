package screens

import (
	"fmt"

	upload "app/app/cmd/cli/client"
)

func ImportInputName(option string) error {
	path := PickFile()
	url := fmt.Sprintf("http://localhost:8080/imports/%s", option)
	return upload.UploadCSV(path, url)
}

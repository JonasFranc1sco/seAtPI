package screens

import (
	"fmt"

	"app/app/cmd/cli/components"
	"app/app/cmd/cli/client"
)

func ImportInputName(option string) error {
	path := filepicker.PickFile()
	url := fmt.Sprintf("http://localhost:8080/imports/%s", option)
	return upload.UploadCSV(path, url)
}



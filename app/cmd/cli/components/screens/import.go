package screens

import (
	"app/app/cmd/cli/components"
	"app/app/cmd/cli/client"
)

func ImportInputName() error {
	path := filepicker.PickFile()
	return upload.UploadCSV(path, "http://localhost:8080/imports/antivirus")
}



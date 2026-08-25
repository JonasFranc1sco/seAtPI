package screens

import (
	"app/app/cmd/cli/components/ui"
	"fmt"
	"log"
	"os/exec"
)

func ExportInputName() {
	path := ui.ImportInputName()
	fileName := fmt.Sprintf("%s.csv", path)

	ExportCommand(fileName)
}

func ExportCommand(fileName string) {
	// Monta o comando pra exportar os arquivos
	cmd := exec.Command(
		"curl",
		"-o",
		fileName,
		"http://localhost:8080/reports/juncao",
	)

	// Roda o comando
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}

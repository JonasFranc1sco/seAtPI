package cli

import (
	"app/app/cmd/cli/components/screens"
	"fmt"
	"log"

	"charm.land/huh/v2"
)

func CliInterface() {
	var option string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Escolha qual ação deseja executar.").
				Options(
					huh.NewOption("Importar CSV do inventário.", "inventario"),
					huh.NewOption("Importar CSV do antivírus.", "antivirus"),
					huh.NewOption("Importar CSV do equipamento.", "patrimonio"),
					huh.NewOption("Exportar CSV do relatório completo.", "relatorio"),
				).
				Value(&option),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if option != "relatorio" {
		if err := screens.ImportInputName(option); err != nil {
			log.Fatal(err)
		}
	}
	switch option {
	case "relatorio":
		screens.ExportInputName()
	}

	fmt.Println("OK")
}

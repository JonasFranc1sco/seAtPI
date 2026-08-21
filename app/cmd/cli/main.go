package main

import (
	"charm.land/huh/v2"
	"log"

	"app/app/cmd/cli/components/screens"
)


func main() {
	var option string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
			Title("Escolha qual ação deseja executar.").
			Options(
				huh.NewOption("Importar CSV do inventário.", "inventario"),
				huh.NewOption("Importar CSV do antivírus.", "antivirus"),
				huh.NewOption("Importar CSV do equipamento.", "equipamento"),
				huh.NewOption("Exportar CSV do relatório completo.", "relatorio"),
			).
			Value(&option),
		),
	)
	
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	switch option {
	case "antivirus":
		if err := screens.ImportInputName(); err != nil {
			log.Fatal(err)
		}
	}
}
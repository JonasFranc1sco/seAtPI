package screens

import "charm.land/huh/v2"

func ExportInputName() {
	var name string

	huh.NewGroup(
		huh.NewInput().
			Title("Como quer nomear o relatório?").
			Prompt("?").
			Value(&name),
	)
}



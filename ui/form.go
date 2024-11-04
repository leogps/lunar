package ui

import (
	"github.com/rivo/tview"
)

func StartApp() error {
	app := tview.NewApplication()

	// Function to show the main form
	showMainForm := func(orderType string) {
		var root *tview.Flex
		if orderType == "ESPP" {
			root = loadEspp(app)
		} else {
			root = loadRsu(app)
		}
		app.SetRoot(root, true) // Set the root to the new form layout
	}

	// Create a dropdown for selecting ESPP or RSU
	selectBox := tview.NewDropDown().
		SetLabel("Select Order Type (hit Enter): ").
		SetOptions([]string{"ESPP", "RSU"}, func(option string, index int) {
			// Show the corresponding form when an option is selected
			showMainForm(option)
		})

	// Set up a Flex layout for the select box and result TextView
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(selectBox, 0, 1, true)

	// Set the root to the main Flex layout
	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		return err
	}
	return nil
}

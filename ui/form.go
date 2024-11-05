package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type DataView int

const (
	EsppOrderSummary DataView = iota
	EsppTargetProfits
	EsppError
	RsuOrderSummary
	RsuTargetProfits
	RsuError
)

var currentDataView DataView

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

// clearFlexItems removes all items in the *tview.Flex
func clearFlexItems(summary *tview.Flex) {
	// Loop in reverse to remove all items
	for i := summary.GetItemCount() - 1; i >= 0; i-- {
		summary.RemoveItem(summary.GetItem(i))
	}
}

// acceptIntInputValue validates the input to only allow int values
func acceptIntInputValue(text string, _ rune) bool {
	_, err := strconv.ParseInt(text, 0, 64)
	return err == nil || text == ""
}

// acceptFloat64InputValue validates the input to only allow float64 values
func acceptFloat64InputValue(text string, _ rune) bool {
	_, err := strconv.ParseFloat(text, 64)
	return err == nil || text == ""
}

func enableTableScroll(table *tview.Table) {
	totalRowCount := table.GetRowCount()
	// Manage selection
	selectedRow := 1 // Start with the first data row selected
	table.SetSelectable(true, false)
	table.Select(selectedRow, 0)
	// Set up key event handling for scrolling
	// Set up key event handling for scrolling
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if selectedRow > 1 {
				selectedRow--                // Move up in the table
				table.Select(selectedRow, 0) // Update the selection
			}
			return nil
		case tcell.KeyDown:
			if selectedRow < totalRowCount-1 {
				selectedRow++                // Move down in the table
				table.Select(selectedRow, 0) // Update the selection
			}
			return nil
		}
		return event
	})

	// Handle mouse scroll
	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if event.Buttons() == tcell.WheelUp {
			if selectedRow > 1 {
				selectedRow--                // Move up in the table
				table.Select(selectedRow, 0) // Update the selection
			}
			return action, event
		} else if event.Buttons() == tcell.WheelDown {
			if selectedRow < totalRowCount-1 {
				selectedRow++                // Move down in the table
				table.Select(selectedRow, 0) // Update the selection
			}
			return action, event
		}
		return action, event
	})
}

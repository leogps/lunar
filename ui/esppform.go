/*
 * Copyright (c) 2024, Paul Gundarapu.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/leogps/lunar/pkg/types"
	"github.com/rivo/tview"
	"strconv"
)

func loadEspp(app *tview.Application) *tview.Flex {
	orderType := "ESPP"
	// Create a TextView for displaying results
	status := tview.NewTextView().SetTextAlign(tview.AlignLeft).
		SetText("Please enter data into fields...").SetTextColor(tview.Styles.PrimaryTextColor)
	summary := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Buying Group
	costPerShare := tview.NewInputField().
		SetLabel("Cost price per share ($)").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)

	discountPercent := tview.NewInputField().
		SetLabel("Discounted (buying) price percent per share (%)").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)

	form := tview.NewForm()

	form.AddFormItem(costPerShare).
		AddFormItem(discountPercent)

	// Selling Group
	sellingPricePerShare := tview.NewInputField().
		SetLabel("Selling price per share ($)").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)

	shareQty := tview.NewInputField().
		SetLabel("Number of shares sold").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptIntInputValue)

	form.AddFormItem(sellingPricePerShare).
		AddFormItem(shareQty)

	// Commission Group
	commissionAmountField := tview.NewInputField().
		SetLabel("Commission Fee Amount per Transaction ($): ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	commissionAmountField.SetDisabled(true)

	numTransactionsField := tview.NewInputField().
		SetLabel("Number of Transactions: ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptIntInputValue)
	numTransactionsField.SetDisabled(true)

	commissionCheckbox := tview.NewCheckbox().
		SetLabel("Add commission fee (hit Enter/Space to toggle): ").
		SetChangedFunc(func(checked bool) {
			if !checked {
				commissionAmountField.SetText("")
				numTransactionsField.SetText("")
			}
			commissionAmountField.SetDisabled(!checked)
			numTransactionsField.SetDisabled(!checked)
		})

	form.AddFormItem(commissionCheckbox).
		AddFormItem(commissionAmountField).
		AddFormItem(numTransactionsField)

	// Tax Group
	capitalGainTaxField := tview.NewInputField().
		SetLabel("Capital Gain Tax Percent percent (Short-Term: 10%-35%) (Long-Term: 0%-20%): ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	capitalGainTaxField.SetDisabled(true)

	taxCheckbox := tview.NewCheckbox().SetLabel("Calculate Capital Gain Tax (hit Enter/Space					 to toggle): ").SetChangedFunc(func(checked bool) {
		if !checked {
			capitalGainTaxField.SetText("")
		}
		capitalGainTaxField.SetDisabled(!checked)
	})

	form.AddFormItem(taxCheckbox).
		AddFormItem(capitalGainTaxField)

	// Create a Submit Button
	form.AddButton("Submit", func() {
		costPerShareValue, _ := strconv.ParseFloat(costPerShare.GetText(), 64)
		discountPercentValue, _ := strconv.ParseFloat(discountPercent.GetText(), 64)
		sellingPricePerShareValue, _ := strconv.ParseFloat(sellingPricePerShare.GetText(), 64)
		shareQtyValue, _ := strconv.Atoi(shareQty.GetText())
		considerCommission := commissionCheckbox.IsChecked()
		commissionAmount, _ := strconv.ParseFloat(commissionAmountField.GetText(), 64)
		numOfTransactions, _ := strconv.Atoi(numTransactionsField.GetText())
		considerCapitalGainTax := taxCheckbox.IsChecked()
		capitalGainTax, _ := strconv.ParseFloat(capitalGainTaxField.GetText(), 64)

		esppOrder := buildEsppOrder(costPerShareValue,
			discountPercentValue,
			sellingPricePerShareValue,
			shareQtyValue,
			considerCommission,
			commissionAmount,
			numOfTransactions,
			considerCapitalGainTax,
			capitalGainTax)

		calculateEspp(esppOrder, status, summary)
	})

	form.AddButton("Target Profits", func() {
		costPerShareValue, _ := strconv.ParseFloat(costPerShare.GetText(), 64)
		discountPercentValue, _ := strconv.ParseFloat(discountPercent.GetText(), 64)
		sellingPricePerShareValue, _ := strconv.ParseFloat(sellingPricePerShare.GetText(), 64)
		shareQtyValue, _ := strconv.Atoi(shareQty.GetText())
		considerCommission := commissionCheckbox.IsChecked()
		commissionAmount, _ := strconv.ParseFloat(commissionAmountField.GetText(), 64)
		numOfTransactions, _ := strconv.Atoi(numTransactionsField.GetText())
		considerCapitalGainTax := taxCheckbox.IsChecked()
		capitalGainTax, _ := strconv.ParseFloat(capitalGainTaxField.GetText(), 64)

		esppOrder := buildEsppOrder(costPerShareValue,
			discountPercentValue,
			sellingPricePerShareValue,
			shareQtyValue,
			considerCommission,
			commissionAmount,
			numOfTransactions,
			considerCapitalGainTax,
			capitalGainTax)
		calculateEsppTargetProfits(esppOrder, status, summary, form, app)
	})

	// Create a Exit Button
	form.AddButton("Exit", func() {
		app.Stop() // Close the app without submission
	})

	separator := tview.NewBox().
		SetBorder(false).
		SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw a horizontal line across the middle of the box.
			centerY := y + height/2
			for cx := x + 1; cx < x+width-1; cx++ {
				screen.SetContent(cx, centerY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			// Write some text along the horizontal line.
			tview.Print(screen, "", x+1, centerY, width-2, tview.AlignCenter, tcell.ColorYellow)

			// Space for other content.
			return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		})

	// Set up a Flex layout to arrange the form and the result TextView
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(separator, 1, 1, false).
		AddItem(status, 1, 1, false).
		AddItem(summary, 0, 1, false)
	flex.
		SetBorder(true).
		SetTitle(fmt.Sprintf("** %s Order **", orderType)).
		SetTitleAlign(tview.AlignCenter)

	return flex // Return the flex layout
}

func calculateEsppTargetProfits(esppOrder *types.EsppOrder,
	status *tview.TextView,
	summary *tview.Flex,
	form *tview.Form,
	app *tview.Application) {

	status.SetText("Calculating...")
	clearFlexItems(summary)

	// Create a new table
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1)

	for index, header := range []string{
		"Profit %",
		"Selling price/share ($)",
		"Total Selling Price ($)",
		"Total Cost ($)",
		"Effective Commission ($)",
		"Profit Before Tax ($)",
		"Capital Gain Tax ($)",
		"Profit After Tax ($)",
	} {
		table.SetCell(0, index, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	// Populate the table with selling prices for each target profit percentage
	row := 1
	for percent := 0.0; percent <= 100.0; percent += 5 {
		// Calculate the selling price for the current percentage
		sellingPrice, err := esppOrder.CalculateSellingPriceForTargetProfitPercent(percent)
		if err != nil {
			sellingPrice = -1 // Handle error case by setting to 0 or any fallback
		}

		// Set profit percentage and calculated selling price in the table
		var targetProfitHeader = fmt.Sprintf("%.0f%%", percent)
		if percent == 0 {
			targetProfitHeader = fmt.Sprintf("%.0f%%", percent)
		}
		var col = 0
		table.SetCell(row, col, tview.NewTableCell(targetProfitHeader).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", sellingPrice)).
			SetAlign(tview.AlignCenter))
		col++

		esppOrderClone := esppOrder.Clone()
		esppOrderClone.SellingPricePerShare = sellingPrice
		esppOrderSummary := esppOrderClone.CalculateEsppOrderSummary()
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.TotalSellingPrice)).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.TotalCost)).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.EffectiveCommission)).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.NetResult)).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.CapitalGainTaxAmount)).
			SetAlign(tview.AlignCenter))
		col++
		table.SetCell(row, col, tview.NewTableCell(fmt.Sprintf("$%.2f", esppOrderSummary.ProfitOrLossAfterCapitalGainsTax())).
			SetAlign(tview.AlignCenter))
		col++

		row++
	}

	enableTableScroll(table)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlI:
			app.SetFocus(form)
		}
		return event
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlD:
			if currentDataView == EsppTargetProfits {
				app.SetFocus(table)
			}
		}
		return event
	})

	// Set up the layout with the table
	summary.
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true)

	status.SetText("Target Profits: [ <ctrl+i> to switch focus to input form | <ctrl+d> to switch focus to Data View ]")
	app.SetFocus(table)
	currentDataView = EsppTargetProfits
}

func buildEsppOrder(costPerShare float64,
	discountPercent float64,
	sellingPricePerShare float64,
	shareQty int,
	considerCommission bool,
	commissionAmount float64,
	numOfTransactions int,
	considerCapitalGainTax bool,
	capitalGainTax float64) *types.EsppOrder {
	esppOrder := types.EsppOrder{}
	// Retrieve values
	esppOrder.CostPerShare = costPerShare
	esppOrder.DiscountPercent = discountPercent
	esppOrder.SellingPricePerShare = sellingPricePerShare
	esppOrder.NumberOfSharesSold = shareQty

	if considerCommission {
		esppOrder.ConsiderTransactionCommission = true
		esppOrder.CommissionPaidPerTransaction = commissionAmount
		esppOrder.NumberOfTransactions = numOfTransactions
	}

	if considerCapitalGainTax {
		esppOrder.ConsiderCapitalGainTax = true
		esppOrder.CapitalGainTaxPercent = capitalGainTax
	}
	return &esppOrder
}

func calculateEspp(esppOrder *types.EsppOrder,
	status *tview.TextView,
	summary *tview.Flex) {
	status.SetText("Calculating...")
	clearFlexItems(summary)

	esppOrderSummary := esppOrder.CalculateEsppOrderSummary()
	costField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Cost: $%.2f", esppOrder.CostPerShare)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(costField, 1, 1, false)

	discountField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Discount: %.2f%%", esppOrder.DiscountPercent)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(discountField, 1, 1, false)

	effectiveCostPerShare := esppOrderSummary.EffectiveCostPerShare
	effectiveCostPerShareField := tview.NewTextView().
		SetLabel(fmt.Sprintf("True cost per share: $%.2f", effectiveCostPerShare)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(effectiveCostPerShareField, 1, 1, false)

	sellingPricePerShareField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Selling price per share: $%.2f", esppOrder.SellingPricePerShare)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(sellingPricePerShareField, 1, 1, false)

	numberOfSharesSoldField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Number of shares sold: %d", esppOrder.NumberOfSharesSold)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(numberOfSharesSoldField, 1, 1, false)

	totalSellingPrice := esppOrderSummary.TotalSellingPrice
	totalSellingPriceField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Total selling price (%d * $%.2f): $%.2f",
			esppOrder.NumberOfSharesSold, esppOrder.SellingPricePerShare, totalSellingPrice)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(totalSellingPriceField, 1, 1, false)

	totalCost := esppOrderSummary.TotalCost
	totalCostField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Total cost (%d * $%.2f): $%.2f",
			esppOrder.NumberOfSharesSold, esppOrderSummary.EffectiveCostPerShare, totalCost)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(totalCostField, 1, 1, false)

	if esppOrder.ConsiderTransactionCommission {
		effectiveTransactionCommission := esppOrderSummary.EffectiveCommission
		effectiveTransactionCommissionField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Effective commission fee (%d * $%.2f): $%.2f",
				esppOrder.NumberOfTransactions, esppOrder.CommissionPaidPerTransaction, effectiveTransactionCommission)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(effectiveTransactionCommissionField, 1, 1, false)
	}

	profitOrLoss := esppOrderSummary.NetResult
	if esppOrderSummary.NetResult > 0 {
		profitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Profit (before capital gains tax): $%.2f", profitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(profitOrLossField, 1, 1, false)

		if esppOrder.ConsiderCapitalGainTax {
			capitalGainTaxAmount := esppOrderSummary.CapitalGainTaxAmount
			capitalGainTaxAmountField := tview.NewTextView().
				SetLabel(fmt.Sprintf("Captial gain tax amount: $%.2f", capitalGainTaxAmount)).
				SetTextAlign(tview.AlignLeft)
			summary.AddItem(capitalGainTaxAmountField, 1, 1, false)

			effectiveProfit := esppOrderSummary.ProfitOrLossAfterCapitalGainsTax()
			effectiveProfitField := tview.NewTextView().
				SetLabel(fmt.Sprintf("Profit (after capital gain tax): $%.2f", effectiveProfit)).
				SetTextAlign(tview.AlignLeft)
			summary.AddItem(effectiveProfitField, 1, 1, false)
		}
	} else if profitOrLoss < 0 {
		profitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Loss: $%.2f", profitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(profitOrLossField, 1, 1, false)
	} else {
		profitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Break even: $%.2f", profitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(profitOrLossField, 1, 1, false)
	}

	gainOrLossMargin := esppOrderSummary.ProfitOrLossMargin()
	gainOrLossMarginField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Gain/Loss Margin: %.2f%%", gainOrLossMargin)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(gainOrLossMarginField, 1, 1, false)

	status.SetText("Summary: ")
	currentDataView = EsppOrderSummary
}

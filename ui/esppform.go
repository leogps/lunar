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
		SetAcceptanceFunc(func(text string, lastChar rune) bool {
			_, err := strconv.ParseFloat(text, 64)
			return err == nil || text == ""
		})

	discountPercent := tview.NewInputField().
		SetLabel("Discounted (buying) price percent per share (%)").
		SetFieldWidth(20).
		SetAcceptanceFunc(func(text string, lastChar rune) bool {
			_, err := strconv.ParseFloat(text, 64)
			return err == nil || text == ""
		})

	form := tview.NewForm()

	form.AddFormItem(costPerShare).
		AddFormItem(discountPercent)

	// Selling Group
	sellingPricePerShare := tview.NewInputField().
		SetLabel("Selling price per share ($)").
		SetFieldWidth(20).
		SetAcceptanceFunc(func(text string, lastChar rune) bool {
			_, err := strconv.ParseFloat(text, 64)
			return err == nil || text == ""
		})

	shareQty := tview.NewInputField().
		SetLabel("Number of shares sold").
		SetFieldWidth(20).
		SetAcceptanceFunc(func(text string, lastChar rune) bool {
			_, err := strconv.ParseInt(text, 0, 64)
			return err == nil || text == ""
		})

	form.AddFormItem(sellingPricePerShare).
		AddFormItem(shareQty)

	// Commission Group
	commissionAmountField := tview.NewInputField().
		SetLabel("Commission Amount per Transaction ($): ").
		SetFieldWidth(20)
	commissionAmountField.SetDisabled(true)

	numTransactionsField := tview.NewInputField().
		SetLabel("Number of Transactions: ").
		SetFieldWidth(20)
	numTransactionsField.SetDisabled(true)

	commissionCheckbox := tview.NewCheckbox().
		SetLabel("Add commission fee (hit Enter to toggle): ").
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
		SetFieldWidth(20)
	capitalGainTaxField.SetDisabled(true)

	taxCheckbox := tview.NewCheckbox().SetLabel("Calculate Capital Gain Tax (hit Enter to toggle): ").SetChangedFunc(func(checked bool) {
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
		calculateEspp(costPerShareValue,
			discountPercentValue,
			sellingPricePerShareValue,
			shareQtyValue,
			considerCommission,
			commissionAmount,
			numOfTransactions,
			considerCapitalGainTax,
			capitalGainTax,
			status,
			summary)
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

func calculateEspp(costPerShare float64, discountPercent float64, sellingPricePerShare float64, shareQty int, considerCommission bool, commissionAmount float64, numOfTransactions int, considerCapitalGainTax bool, capitalGainTax float64, status *tview.TextView, summary *tview.Flex) {
	status.SetText("Calculating...")
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

	// Loop in reverse to remove all items
	for i := summary.GetItemCount() - 1; i >= 0; i-- {
		summary.RemoveItem(summary.GetItem(i))
	}

	costField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Cost: $%.2f", esppOrder.CostPerShare)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(costField, 1, 1, false)

	discountField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Discount: %.2f%%", esppOrder.DiscountPercent)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(discountField, 1, 1, false)

	effectiveCostPerShare := esppOrder.CalculateEffectiveCostPerShare()
	effectiveCostPerShareField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Effective cost per share: $%.2f", effectiveCostPerShare)).
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

	totalSellingPrice := esppOrder.SellingPricePerShare * float64(esppOrder.NumberOfSharesSold)
	totalSellingPriceField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Total selling price (%d * $%.2f): $%.2f",
			esppOrder.NumberOfSharesSold, esppOrder.SellingPricePerShare, totalSellingPrice)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(totalSellingPriceField, 1, 1, false)

	if esppOrder.ConsiderTransactionCommission {
		effectiveTransactionCommission := float64(esppOrder.NumberOfTransactions) * esppOrder.CommissionPaidPerTransaction
		effectiveTransactionCommissionField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Effective commission fee (%d * $%.2f): $%.2f",
				esppOrder.NumberOfTransactions, esppOrder.CommissionPaidPerTransaction, effectiveTransactionCommission)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(effectiveTransactionCommissionField, 1, 1, false)
	}

	profitOrLoss := esppOrder.CalculateProfitOrLoss()
	if profitOrLoss > 0 {
		profitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Profit (before capital gains tax): $%.2f", profitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(profitOrLossField, 1, 1, false)

		if esppOrder.ConsiderCapitalGainTax {
			capitalGainTaxAmount, err := esppOrder.CalculateCapitalGainTaxAmount(profitOrLoss)
			if err != nil {
				status.SetText(fmt.Sprintf("Error: %v", err.Error()))
				return
			}
			capitalGainTaxAmountField := tview.NewTextView().
				SetLabel(fmt.Sprintf("Captial gain tax amount: $%.2f", capitalGainTaxAmount)).
				SetTextAlign(tview.AlignLeft)
			summary.AddItem(capitalGainTaxAmountField, 1, 1, false)

			effectiveProfit := profitOrLoss - capitalGainTaxAmount
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

	status.SetText("Summary: ")
}

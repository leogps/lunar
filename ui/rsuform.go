package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/leogps/lunar/pkg/types"
	"github.com/rivo/tview"
	"strconv"
)

func loadRsu(app *tview.Application) *tview.Flex {
	orderType := "RSU"
	// Create a TextView for displaying results
	status := tview.NewTextView().SetTextAlign(tview.AlignLeft).
		SetText("Please enter data into fields...").SetTextColor(tview.Styles.PrimaryTextColor)
	summary := tview.NewFlex().
		SetDirection(tview.FlexRow)

	form := tview.NewForm()

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
		SetLabel("Commission Amount per Transaction ($): ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	commissionAmountField.SetDisabled(true)

	numTransactionsField := tview.NewInputField().
		SetLabel("Number of Transactions: ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptIntInputValue)
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
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	capitalGainTaxField.SetDisabled(true)

	taxCheckbox := tview.NewCheckbox().SetLabel("Calculate Capital Gain Tax (hit Enter to toggle): ").SetChangedFunc(func(checked bool) {
		if !checked {
			capitalGainTaxField.SetText("")
		}
		capitalGainTaxField.SetDisabled(!checked)
	})

	form.AddFormItem(taxCheckbox).
		AddFormItem(capitalGainTaxField)

	// Income tax Group
	incomeTaxField := tview.NewInputField().
		SetLabel("Income Tax percent (State + Federal): ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	incomeTaxField.SetDisabled(true)

	marketPriceOnVestedStockPerShareField := tview.NewInputField().
		SetLabel("Market Price On Vested Stock per share ($): ").
		SetFieldWidth(20).
		SetAcceptanceFunc(acceptFloat64InputValue)
	marketPriceOnVestedStockPerShareField.SetDisabled(true)

	incomeTaxCheckbox := tview.NewCheckbox().SetLabel("Calculate Income Tax (hit Enter to toggle): ").SetChangedFunc(func(checked bool) {
		if !checked {
			incomeTaxField.SetText("")
			marketPriceOnVestedStockPerShareField.SetText("")
		}
		incomeTaxField.SetDisabled(!checked)
		marketPriceOnVestedStockPerShareField.SetDisabled(!checked)
	})

	form.AddFormItem(incomeTaxCheckbox).
		AddFormItem(incomeTaxField).
		AddFormItem(marketPriceOnVestedStockPerShareField)

	// Create a Submit Button
	form.AddButton("Submit", func() {
		sellingPricePerShareValue, _ := strconv.ParseFloat(sellingPricePerShare.GetText(), 64)
		shareQtyValue, _ := strconv.Atoi(shareQty.GetText())
		considerCommission := commissionCheckbox.IsChecked()
		commissionAmount, _ := strconv.ParseFloat(commissionAmountField.GetText(), 64)
		numOfTransactions, _ := strconv.Atoi(numTransactionsField.GetText())
		considerCapitalGainTax := taxCheckbox.IsChecked()
		capitalGainTax, _ := strconv.ParseFloat(capitalGainTaxField.GetText(), 64)

		considerIncomeTax := incomeTaxCheckbox.IsChecked()
		incomeTaxPercent, _ := strconv.ParseFloat(incomeTaxField.GetText(), 64)
		marketPriceOnVestedStockPerShare, _ := strconv.ParseFloat(marketPriceOnVestedStockPerShareField.GetText(), 64)

		rsuOrder := buildRsuOrder(sellingPricePerShareValue,
			shareQtyValue,
			considerCommission,
			commissionAmount,
			numOfTransactions,
			considerCapitalGainTax,
			capitalGainTax,
			considerIncomeTax,
			incomeTaxPercent,
			marketPriceOnVestedStockPerShare)
		calculateRsu(rsuOrder, status, summary)
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

func buildRsuOrder(sellingPricePerShare float64,
	shareQty int,
	considerCommission bool,
	commissionAmount float64,
	numOfTransactions int,
	considerCapitalGainTax bool,
	capitalGainTax float64,
	considerIncomeTax bool,
	incomeTaxPercent float64,
	marketPriceOnVestedStockPerShare float64) *types.RsuOrder {
	rsuOrder := types.RsuOrder{}
	// Retrieve values
	rsuOrder.SellingPricePerShare = sellingPricePerShare
	rsuOrder.NumberOfSharesSold = shareQty

	if considerCommission {
		rsuOrder.ConsiderTransactionCommission = true
		rsuOrder.CommissionPaidPerTransaction = commissionAmount
		rsuOrder.NumberOfTransactions = numOfTransactions
	}

	if considerCapitalGainTax {
		rsuOrder.ConsiderCapitalGainTax = true
		rsuOrder.CapitalGainTaxPercent = capitalGainTax
	}

	if considerIncomeTax {
		rsuOrder.ConsiderIncomeTaxOnVestedStock = true
		rsuOrder.IncomeTaxPercentOnVestedStockPerShare = incomeTaxPercent
		rsuOrder.MarketPriceOnVestedStockPerShare = marketPriceOnVestedStockPerShare
	}
	return &rsuOrder
}

func calculateRsu(rsuOrder *types.RsuOrder,
	status *tview.TextView,
	summary *tview.Flex) {
	status.SetText("Calculating...")
	clearFlexItems(summary)

	rsuOrderSummary := rsuOrder.CalculateRsuOrderSummary()

	sellingPricePerShareField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Selling price per share: $%.2f", rsuOrder.SellingPricePerShare)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(sellingPricePerShareField, 1, 1, false)

	numberOfSharesSoldField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Number of shares sold: %d", rsuOrder.NumberOfSharesSold)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(numberOfSharesSoldField, 1, 1, false)

	totalSellingPrice := rsuOrderSummary.TotalSellingPrice
	totalSellingPriceField := tview.NewTextView().
		SetLabel(fmt.Sprintf("Total selling price (%d * $%.2f): $%.2f",
			rsuOrder.NumberOfSharesSold, rsuOrder.SellingPricePerShare, totalSellingPrice)).
		SetTextAlign(tview.AlignLeft)
	summary.AddItem(totalSellingPriceField, 1, 1, false)

	if rsuOrder.ConsiderTransactionCommission {
		effectiveTransactionCommission := rsuOrderSummary.EffectiveCommission
		effectiveTransactionCommissionField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Effective commission fee (%d * $%.2f): $%.2f",
				rsuOrder.NumberOfTransactions, rsuOrder.CommissionPaidPerTransaction, effectiveTransactionCommission)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(effectiveTransactionCommissionField, 1, 1, false)
	}

	profitOrLoss := rsuOrderSummary.NetResult
	if rsuOrderSummary.IsProfitable {
		profitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Profit (before capital gains tax): $%.2f", profitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(profitOrLossField, 1, 1, false)

		if rsuOrder.ConsiderCapitalGainTax {
			capitalGainTaxAmount := rsuOrderSummary.CapitalGainTaxAmount
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

	if rsuOrder.ConsiderIncomeTaxOnVestedStock {
		totalIncomeTaxIncurred := rsuOrderSummary.TotalIncomeTaxIncurred
		totalIncomeTaxIncurredField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Total income tax amount: $%.2f", totalIncomeTaxIncurred)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(totalIncomeTaxIncurredField, 1, 1, false)

		effectiveProfitOrLoss := rsuOrderSummary.ProfitOrLossAfterIncomeTax
		effectiveProfitOrLossField := tview.NewTextView().
			SetLabel(fmt.Sprintf("Effective Profit/Loss: $%.2f", effectiveProfitOrLoss)).
			SetTextAlign(tview.AlignLeft)
		summary.AddItem(effectiveProfitOrLossField, 1, 1, false)
	}

	status.SetText("Summary: ")
	currentDataView = "RSU_ORDER_SUMMARY"
}

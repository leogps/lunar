package types

import (
	"fmt"
	"math"
	"strings"
)

type EsppOrder struct {
	DiscountPercent      float64
	CostPerShare         float64
	SellingPricePerShare float64
	NumberOfSharesSold   int

	ConsiderTransactionCommission bool
	CommissionPaidPerTransaction  float64
	NumberOfTransactions          int

	ConsiderCapitalGainTax bool
	CapitalGainTaxPercent  float64
}

type EsppOrderSummary struct {
	EsppOrder             *EsppOrder
	EffectiveCostPerShare float64
	TotalSellingPrice     float64
	TotalCost             float64
	EffectiveCommission   float64
	NetResult             float64
	CapitalGainTaxAmount  float64
}

func (e *EsppOrderSummary) IsProfitable() bool {
	return e.TrueProfitOrLoss() > 0
}

func (e *EsppOrderSummary) ProfitOrLossAfterCapitalGainsTax() float64 {
	return e.NetResult - e.CapitalGainTaxAmount
}

func (e *EsppOrderSummary) TrueProfitOrLoss() float64 {
	trueProfitOrLoss := e.NetResult
	if e.EsppOrder.ConsiderCapitalGainTax {
		trueProfitOrLoss -= e.CapitalGainTaxAmount
	}
	return trueProfitOrLoss
}

func (e *EsppOrderSummary) ProfitOrLossMargin() float64 {
	effectiveProfit := e.TrueProfitOrLoss()
	return (effectiveProfit / (e.TotalCost + e.EffectiveCommission + e.CapitalGainTaxAmount)) * 100
}

func (e *EsppOrderSummary) ToString() string {
	var sb strings.Builder

	sb.WriteString("ESPP Order Summary:\n")
	sb.WriteString(fmt.Sprintf("  Effective Cost Per Share:      $%.2f\n", e.EffectiveCostPerShare))
	sb.WriteString(fmt.Sprintf("  Total Selling Price:           $%.2f\n", e.TotalSellingPrice))
	sb.WriteString(fmt.Sprintf("  Total Cost:                    $%.2f\n", e.TotalCost))
	sb.WriteString(fmt.Sprintf("  Effective Commission:          $%.2f\n", e.EffectiveCommission))
	sb.WriteString(fmt.Sprintf("  Net Result: 					$%.2f\n", e.NetResult))
	sb.WriteString(fmt.Sprintf("  Capital Gain Tax Amount:       $%.2f\n", e.CapitalGainTaxAmount))
	sb.WriteString(fmt.Sprintf("  Profit After Capital Gain Tax: $%.2f\n", e.ProfitOrLossAfterCapitalGainsTax()))
	sb.WriteString(fmt.Sprintf("  Profit/Loss Margin: 			%.2f%%\n", e.ProfitOrLossMargin()))
	sb.WriteString(fmt.Sprintf("  Net Result:                    $%.2f\n", e.NetResult))
	sb.WriteString(fmt.Sprintf("  Is Profitable:                 %t\n", e.IsProfitable()))

	return sb.String()
}

// Clone creates a deep copy of the EsppOrder
func (e *EsppOrder) Clone() *EsppOrder {
	return &EsppOrder{
		DiscountPercent:               e.DiscountPercent,
		CostPerShare:                  e.CostPerShare,
		SellingPricePerShare:          e.SellingPricePerShare,
		NumberOfSharesSold:            e.NumberOfSharesSold,
		ConsiderTransactionCommission: e.ConsiderTransactionCommission,
		CommissionPaidPerTransaction:  e.CommissionPaidPerTransaction,
		NumberOfTransactions:          e.NumberOfTransactions,
		ConsiderCapitalGainTax:        e.ConsiderCapitalGainTax,
		CapitalGainTaxPercent:         e.CapitalGainTaxPercent,
	}
}

func (e *EsppOrder) CalculateDiscountAmount() float64 {
	return (e.CostPerShare * e.DiscountPercent) / 100
}

func (e *EsppOrder) CalculateEffectiveCostPerShare() float64 {
	discountAmount := e.CalculateDiscountAmount()
	return e.CostPerShare - discountAmount
}

func (e *EsppOrder) CalculateProfitOrLoss() float64 {
	effectiveCostPerShare := e.CalculateEffectiveCostPerShare()
	var effectiveTransactionCommission float64
	effectiveTransactionCommission = 0
	if e.ConsiderTransactionCommission {
		effectiveTransactionCommission = float64(e.NumberOfTransactions) * e.CommissionPaidPerTransaction
	}
	return (float64(e.NumberOfSharesSold) * (e.SellingPricePerShare - effectiveCostPerShare)) - (effectiveTransactionCommission)
}

func (e *EsppOrder) CalculateCapitalGainTaxAmount(profit float64) (float64, error) {
	if profit < 0 {
		return 0, fmt.Errorf("profit must be greater than or equal to zero")
	}

	return profit * e.CapitalGainTaxPercent / 100, nil
}

func (e *EsppOrder) CalculateEsppOrderSummary() *EsppOrderSummary {
	effectiveCostPerShare := e.CalculateEffectiveCostPerShare()
	totalSellingPrice := e.SellingPricePerShare * float64(e.NumberOfSharesSold)
	totalCost := effectiveCostPerShare * float64(e.NumberOfSharesSold)

	var effectiveTransactionCommission float64
	if e.ConsiderTransactionCommission {
		effectiveTransactionCommission = float64(e.NumberOfTransactions) * e.CommissionPaidPerTransaction
	}
	netResult := e.CalculateProfitOrLoss()
	var capitalGainTaxAmount float64
	if netResult > 0 {
		capitalGainTaxAmount, _ = e.CalculateCapitalGainTaxAmount(netResult)
	}
	return &EsppOrderSummary{
		EsppOrder:             e,
		EffectiveCostPerShare: effectiveCostPerShare,
		TotalSellingPrice:     totalSellingPrice,
		TotalCost:             totalCost,
		EffectiveCommission:   effectiveTransactionCommission,
		NetResult:             netResult,
		CapitalGainTaxAmount:  capitalGainTaxAmount,
	}
}

// CalculateBreakEvenSellingPrice calculates the selling price required to break even.
func (e *EsppOrder) CalculateBreakEvenSellingPrice() float64 {
	breakEvenSellingPrice, _ := e.CalculateSellingPriceForTargetProfitPercent(0)
	return breakEvenSellingPrice
}

// CalculateSellingPriceForTargetProfitPercent calculates the selling price required to achieve a target profit percentage.
func (e *EsppOrder) CalculateSellingPriceForTargetProfitPercent(targetProfitPercent float64) (float64, error) {
	if targetProfitPercent < 0 {
		return 0, fmt.Errorf("target profit percent must be greater than or equal to 0")
	}

	if e.NumberOfSharesSold <= 0 {
		return 0, fmt.Errorf("number of shares sold must be greater than zero")
	}

	effectiveCostPerShare := e.CalculateEffectiveCostPerShare()
	effectiveCost := effectiveCostPerShare * float64(e.NumberOfSharesSold)

	// Initial guess for selling price
	sellingPrice := effectiveCost + (targetProfitPercent/100)*effectiveCost
	for i := 0; i < 100; i++ { // Limit iterations to avoid infinite loops
		profitBeforeTax := (sellingPrice - effectiveCostPerShare) * float64(e.NumberOfSharesSold)
		if e.ConsiderTransactionCommission {
			profitBeforeTax -= float64(e.NumberOfTransactions) * e.CommissionPaidPerTransaction
		}
		capitalGainsTax := 0.0
		if e.ConsiderCapitalGainTax {
			capitalGainsTax = profitBeforeTax * (e.CapitalGainTaxPercent / 100)
		}
		profitAfterTax := profitBeforeTax - capitalGainsTax

		// Calculate the target profit after tax
		targetProfitAfterTax := (targetProfitPercent / 100) * (effectiveCost + capitalGainsTax)

		if math.Abs(profitAfterTax-targetProfitAfterTax) < 0.01 { // Convergence condition
			break
		}
		sellingPrice += (targetProfitAfterTax - profitAfterTax) / float64(e.NumberOfSharesSold) // Adjust selling price
	}

	return sellingPrice, nil
}

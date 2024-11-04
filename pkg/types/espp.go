package types

import "fmt"

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

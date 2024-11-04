package types

import "fmt"

type RsuOrder struct {
	SellingPricePerShare float64
	NumberOfSharesSold   int

	ConsiderTransactionCommission bool
	CommissionPaidPerTransaction  float64
	NumberOfTransactions          int

	ConsiderCapitalGainTax bool
	CapitalGainTaxPercent  float64

	ConsiderIncomeTaxOnVestedStock        bool
	IncomeTaxPercentOnVestedStockPerShare float64
	MarketPriceOnVestedStockPerShare      float64
}

func (r *RsuOrder) CalculateProfitOrLoss() float64 {
	var effectiveTransactionCommission float64
	effectiveTransactionCommission = 0
	if r.ConsiderTransactionCommission {
		effectiveTransactionCommission = float64(r.NumberOfTransactions) * r.CommissionPaidPerTransaction
	}
	return (float64(r.NumberOfSharesSold) * (r.SellingPricePerShare)) - (effectiveTransactionCommission)
}

func (r *RsuOrder) CalculateCapitalGainTaxAmount(profit float64) (float64, error) {
	if profit < 0 {
		return 0, fmt.Errorf("profit must be greater than or equal to zero")
	}

	return profit * r.CapitalGainTaxPercent / 100, nil
}

func (r *RsuOrder) CalculateEffectiveIncomeTaxAmount() float64 {
	return float64(r.NumberOfSharesSold) * r.MarketPriceOnVestedStockPerShare * r.IncomeTaxPercentOnVestedStockPerShare / 100
}

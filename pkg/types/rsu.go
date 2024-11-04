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

type RsuOrderSummary struct {
	RsuOrder                   *RsuOrder
	TotalSellingPrice          float64
	EffectiveCommission        float64
	NetResult                  float64
	IsProfitable               bool
	ProfitBeforeCapitalGainTax float64
	CapitalGainTaxAmount       float64
	ProfitAfterCapitalGainTax  float64
	TotalIncomeTaxIncurred     float64
	ProfitOrLossAfterIncomeTax float64
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

func (r *RsuOrder) CalculateRsuOrderSummary() *RsuOrderSummary {

	totalSellingPrice := r.SellingPricePerShare * float64(r.NumberOfSharesSold)
	effectiveTransactionCommission := float64(r.NumberOfTransactions) * r.CommissionPaidPerTransaction
	netResult := r.CalculateProfitOrLoss()
	isProfitable := netResult > 0

	var capitalGainTaxAmount float64
	var effectiveProfit float64

	if isProfitable {
		capitalGainTaxAmount, _ = r.CalculateCapitalGainTaxAmount(netResult)
		effectiveProfit = netResult - capitalGainTaxAmount
	}

	totalIncomeTaxIncurred := r.CalculateEffectiveIncomeTaxAmount()
	effectiveProfitOrLoss := netResult - totalIncomeTaxIncurred

	return &RsuOrderSummary{
		TotalSellingPrice:          totalSellingPrice,
		EffectiveCommission:        effectiveTransactionCommission,
		NetResult:                  netResult,
		IsProfitable:               isProfitable,
		ProfitBeforeCapitalGainTax: netResult,
		CapitalGainTaxAmount:       capitalGainTaxAmount,
		ProfitAfterCapitalGainTax:  effectiveProfit,
		TotalIncomeTaxIncurred:     totalIncomeTaxIncurred,
		ProfitOrLossAfterIncomeTax: effectiveProfitOrLoss,
	}
}

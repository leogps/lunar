package types

import (
	"fmt"
	"math"
	"strings"
)

type RsuOrder struct {
	SellingPricePerShare float64
	NumberOfSharesSold   int

	ConsiderTransactionCommission bool
	CommissionPaidPerTransaction  float64
	NumberOfTransactions          int

	ConsiderCapitalGainTax bool
	CapitalGainTaxPercent  float64

	ConsiderIncomeTaxOnVestedStock   bool
	IncomeTaxIncurredWhenStockVested float64
	NumberOfStocksVested             int
	MarketValuePerShare              float64
}

type RsuOrderSummary struct {
	RsuOrder               *RsuOrder
	TotalSellingPrice      float64
	EffectiveCommission    float64
	NetResult              float64
	CapitalGainTaxAmount   float64
	TotalIncomeTaxIncurred float64
}

func (r *RsuOrderSummary) ProfitOrLossAfterIncomeTax() float64 {
	return r.NetResult - r.TotalIncomeTaxIncurred
}

func (r *RsuOrderSummary) ProfitOrLossAfterCapitalGainsTax() float64 {
	return r.NetResult - r.CapitalGainTaxAmount
}

func (r *RsuOrderSummary) IsProfitable() bool {
	return r.TrueProfitOrLoss() > 0
}

func (r *RsuOrderSummary) TrueProfitOrLoss() float64 {
	trueProfitOrLoss := r.NetResult
	if r.RsuOrder.ConsiderCapitalGainTax {
		trueProfitOrLoss -= r.CapitalGainTaxAmount
	}
	if r.RsuOrder.ConsiderIncomeTaxOnVestedStock {
		trueProfitOrLoss -= r.TotalIncomeTaxIncurred
	}
	return trueProfitOrLoss
}

func (r *RsuOrderSummary) ProfitOrLossMargin() float64 {
	return (r.TrueProfitOrLoss() / r.TotalIncomeTaxIncurred) * 100
}

func (r *RsuOrderSummary) ToString() string {
	var sb strings.Builder

	sb.WriteString("RSU Order Summary:\n")
	sb.WriteString(fmt.Sprintf("  Total Selling Price:          		$%.2f\n", r.TotalSellingPrice))
	sb.WriteString(fmt.Sprintf("  Effective Commission:         		$%.2f\n", r.EffectiveCommission))
	sb.WriteString(fmt.Sprintf("  Capital Gain Tax Amount:      		$%.2f\n", r.CapitalGainTaxAmount))
	sb.WriteString(fmt.Sprintf("  Total Income Tax Incurred:    		$%.2f\n", r.TotalIncomeTaxIncurred))
	sb.WriteString(fmt.Sprintf("  Net Result:                   		$%.2f\n", r.NetResult))
	sb.WriteString(fmt.Sprintf("  Is Profitable:                		%t\n", r.IsProfitable()))
	sb.WriteString(fmt.Sprintf("  Profit After Capital Gains Tax: 	$%.2f\n", r.ProfitOrLossAfterCapitalGainsTax()))
	sb.WriteString(fmt.Sprintf("  Profit/Loss After Income Tax: 		$%.2f\n", r.ProfitOrLossAfterIncomeTax()))
	sb.WriteString(fmt.Sprintf("  True Profit/Loss: 					$%.2f\n", r.TrueProfitOrLoss()))
	sb.WriteString(fmt.Sprintf("  Profit/Loss Margin: 		   	   	%.2f%%\n", r.ProfitOrLossMargin()))
	return sb.String()
}

// Clone creates a deep copy of the EsppOrder
func (r *RsuOrder) Clone() *RsuOrder {
	return &RsuOrder{
		SellingPricePerShare:             r.SellingPricePerShare,
		NumberOfSharesSold:               r.NumberOfSharesSold,
		ConsiderTransactionCommission:    r.ConsiderTransactionCommission,
		CommissionPaidPerTransaction:     r.CommissionPaidPerTransaction,
		NumberOfTransactions:             r.NumberOfTransactions,
		ConsiderCapitalGainTax:           r.ConsiderCapitalGainTax,
		CapitalGainTaxPercent:            r.CapitalGainTaxPercent,
		ConsiderIncomeTaxOnVestedStock:   r.ConsiderIncomeTaxOnVestedStock,
		IncomeTaxIncurredWhenStockVested: r.IncomeTaxIncurredWhenStockVested,
		NumberOfStocksVested:             r.NumberOfStocksVested,
		MarketValuePerShare:              r.MarketValuePerShare,
	}
}

func (r *RsuOrder) CalculateEffectiveProfitOrLoss() float64 {
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

func (r *RsuOrder) CalculateIncomeTaxPerShare() (float64, error) {
	if r.NumberOfStocksVested <= 0 {
		return 0, fmt.Errorf("number of stocks vested must be greater than zero")
	}
	return r.IncomeTaxIncurredWhenStockVested / float64(r.NumberOfStocksVested), nil
}

func (r *RsuOrder) CalculateTotalIncomeTaxAmount() (float64, error) {
	incomeTaxPerShare, err := r.CalculateIncomeTaxPerShare()
	if err != nil {
		return 0, err
	}
	return incomeTaxPerShare * float64(r.NumberOfSharesSold), nil
}

func (r *RsuOrder) CalculateRsuOrderSummary() (*RsuOrderSummary, error) {

	totalSellingPrice := r.SellingPricePerShare * float64(r.NumberOfSharesSold)
	effectiveTransactionCommission := float64(r.NumberOfTransactions) * r.CommissionPaidPerTransaction
	netResult := r.CalculateEffectiveProfitOrLoss()

	var capitalGainTaxAmount float64

	profitOrLossForCapitalGain := r.CalculateProfitOrLossForCapitalGain()
	if profitOrLossForCapitalGain > 0 {
		capitalGainTaxAmount, _ = r.CalculateCapitalGainTaxAmount(profitOrLossForCapitalGain)
	}

	totalIncomeTaxIncurred, err := r.CalculateTotalIncomeTaxAmount()
	if err != nil {
		return nil, err
	}
	return &RsuOrderSummary{
		RsuOrder:               r,
		TotalSellingPrice:      totalSellingPrice,
		EffectiveCommission:    effectiveTransactionCommission,
		NetResult:              netResult,
		CapitalGainTaxAmount:   capitalGainTaxAmount,
		TotalIncomeTaxIncurred: totalIncomeTaxIncurred,
	}, nil
}

// CalculateSellingPriceForTargetProfitPercent calculates the selling price required to achieve a target profit percentage.
func (r *RsuOrder) CalculateSellingPriceForTargetProfitPercent(targetProfitPercent float64) (float64, error) {
	if targetProfitPercent < 0 {
		return 0, fmt.Errorf("target profit percent must be greater than or equal to 0")
	}
	if r.NumberOfSharesSold <= 0 {
		return 0, fmt.Errorf("number of shares sold must be greater than zero")
	}

	// Calculate income tax per share incurred at vesting, if applicable
	incomeTaxPerShare := 0.0
	var err error
	if r.ConsiderIncomeTaxOnVestedStock {
		incomeTaxPerShare, err = r.CalculateIncomeTaxPerShare()
		if err != nil {
			return 0, err
		}
	}

	// Initial effective cost per share is only the income tax per share
	effectiveCostPerShare := incomeTaxPerShare
	totalEffectiveCost := effectiveCostPerShare * float64(r.NumberOfSharesSold)

	// Initial guess for selling price per share to achieve the target profit
	estimatedSellingPrice := totalEffectiveCost / float64(r.NumberOfSharesSold) * (1 + targetProfitPercent/100)

	// Iteratively adjust the selling price to achieve the target profit percent
	for i := 0; i < 1000; i++ { // Limit iterations for safety
		// Calculate total selling price and profit before tax
		totalSellingPrice := estimatedSellingPrice * float64(r.NumberOfSharesSold)
		profitBeforeTax := totalSellingPrice - totalEffectiveCost

		// Subtract transaction commission, if applicable
		totalTransactionCommission := 0.0
		if r.ConsiderTransactionCommission {
			totalTransactionCommission = float64(r.NumberOfTransactions) * r.CommissionPaidPerTransaction
			profitBeforeTax -= totalTransactionCommission
		}

		// Calculate capital gains tax, only if selling price exceeds market value at vesting
		capitalGainsTax := 0.0
		if r.ConsiderCapitalGainTax && estimatedSellingPrice > r.MarketValuePerShare {
			capitalGainPerShare := estimatedSellingPrice - r.MarketValuePerShare
			totalCapitalGains := capitalGainPerShare * float64(r.NumberOfSharesSold)
			capitalGainsTax = totalCapitalGains * (r.CapitalGainTaxPercent / 100)
			profitBeforeTax -= capitalGainsTax
		}

		// Calculate net profit percentage after all deductions
		actualProfitPercent := (profitBeforeTax / totalEffectiveCost) * 100

		// Check if the achieved profit percentage matches the target profit percentage
		if math.Abs(actualProfitPercent-targetProfitPercent) < 0.01 { // Convergence condition
			break
		}

		// Adjust selling price per share based on the difference between actual and target profit percentages
		estimatedSellingPrice += (targetProfitPercent - actualProfitPercent) / 100 * effectiveCostPerShare
	}

	return estimatedSellingPrice, nil
}

func (r *RsuOrder) CalculateProfitOrLossForCapitalGain() float64 {
	return (float64(r.NumberOfSharesSold) * (r.SellingPricePerShare)) - (r.MarketValuePerShare * float64(r.NumberOfSharesSold))
}

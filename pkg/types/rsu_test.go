package types

import (
	"fmt"
	"testing"
)

func TestRsuOrder_CalculateRsuOrderSummary(t *testing.T) {
	rsuOrder := &RsuOrder{
		SellingPricePerShare:             200.00,
		NumberOfSharesSold:               1,
		ConsiderTransactionCommission:    true,
		CommissionPaidPerTransaction:     5.00,
		NumberOfTransactions:             1,
		ConsiderCapitalGainTax:           true,
		CapitalGainTaxPercent:            24,
		ConsiderIncomeTaxOnVestedStock:   true,
		IncomeTaxIncurredWhenStockVested: 2166.12,
		NumberOfStocksVested:             33,
		MarketValuePerShare:              120.34,
	}

	summary, err := rsuOrder.CalculateRsuOrderSummary()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(summary.ToString())
}

func TestRsuOrder_CalculateSellingPriceForTargetProfitPercent(t *testing.T) {
	rsuOrder := &RsuOrder{
		NumberOfSharesSold:               1,
		ConsiderTransactionCommission:    true,
		CommissionPaidPerTransaction:     5.00,
		NumberOfTransactions:             1,
		ConsiderCapitalGainTax:           true,
		CapitalGainTaxPercent:            24,
		ConsiderIncomeTaxOnVestedStock:   true,
		IncomeTaxIncurredWhenStockVested: 2166.12,
		NumberOfStocksVested:             33,
		MarketValuePerShare:              120.34,
	}

	sellingPrice, err := rsuOrder.CalculateSellingPriceForTargetProfitPercent(0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Selling Price: $%.2f", sellingPrice))

	rsuClone := rsuOrder.Clone()
	rsuClone.SellingPricePerShare = sellingPrice
	summary, err := rsuClone.CalculateRsuOrderSummary()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(summary.ToString())
}

func TestRsuOrder_CalculateSellingPriceFor5PercentProfit(t *testing.T) {
	rsuOrder := &RsuOrder{
		NumberOfSharesSold:               1,
		ConsiderTransactionCommission:    true,
		CommissionPaidPerTransaction:     5.00,
		NumberOfTransactions:             1,
		ConsiderCapitalGainTax:           true,
		CapitalGainTaxPercent:            24,
		ConsiderIncomeTaxOnVestedStock:   true,
		IncomeTaxIncurredWhenStockVested: 2166.12,
		NumberOfStocksVested:             33,
		MarketValuePerShare:              120.34,
	}

	sellingPrice, err := rsuOrder.CalculateSellingPriceForTargetProfitPercent(5)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sellingPrice)

	rsuClone := rsuOrder.Clone()
	rsuClone.SellingPricePerShare = sellingPrice
	summary, err := rsuClone.CalculateRsuOrderSummary()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(summary.ToString())
}

func TestRsuOrder_CalculateSellingPriceFor290PercentProfit(t *testing.T) {
	rsuOrder := &RsuOrder{
		NumberOfSharesSold:               1,
		ConsiderTransactionCommission:    true,
		CommissionPaidPerTransaction:     5.00,
		NumberOfTransactions:             1,
		ConsiderCapitalGainTax:           true,
		CapitalGainTaxPercent:            24,
		ConsiderIncomeTaxOnVestedStock:   true,
		IncomeTaxIncurredWhenStockVested: 2166.12,
		MarketValuePerShare:              190.39,
	}

	sellingPrice, err := rsuOrder.CalculateSellingPriceForTargetProfitPercent(280)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sellingPrice)

	rsuClone := rsuOrder.Clone()
	rsuClone.SellingPricePerShare = sellingPrice
	summary, err := rsuClone.CalculateRsuOrderSummary()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(summary.ToString())
}

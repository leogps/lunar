package types

import (
	"fmt"
	"testing"
)

func TestEsppOrder_CalculateSellingPriceForTargetProfitPercent(t *testing.T) {
	esppOrder := &EsppOrder{
		DiscountPercent:               15,
		CostPerShare:                  100,
		NumberOfSharesSold:            1,
		ConsiderTransactionCommission: true,
		CommissionPaidPerTransaction:  5,
		NumberOfTransactions:          1,
		ConsiderCapitalGainTax:        true,
		CapitalGainTaxPercent:         24,
	}

	breakEvenSellingPrice, err := esppOrder.CalculateSellingPriceForTargetProfitPercent(0)
	if err != nil {
		t.Fatal(err)
	}

	esppClone := esppOrder.Clone()
	esppClone.SellingPricePerShare = breakEvenSellingPrice
	summary := esppClone.CalculateEsppOrderSummary()
	fmt.Println(fmt.Sprintf("Break-Even selling price: $%.2f", breakEvenSellingPrice))
	fmt.Println(summary.ToString())
	if summary.CapitalGainTaxAmount != 0 {
		t.Fail()
	}
}

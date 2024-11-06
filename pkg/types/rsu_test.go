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
		NumberOfStocksVested:             33,
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

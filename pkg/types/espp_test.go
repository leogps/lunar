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

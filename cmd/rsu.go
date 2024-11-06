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

package cmd

import (
	"github.com/leogps/lunar/pkg/types"
	"github.com/leogps/lunar/pkg/utils"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	rootCmd.AddCommand(rsuCmd)
}

var rsuCmd = &cobra.Command{
	Use:   "rsu",
	Short: "calculate profit/loss on RSU orders interactively",
	Long:  `calculate profit/loss on RSU orders interactively`,
	Run: func(cmd *cobra.Command, _ []string) {
		silent, _ := cmd.Flags().GetBool("silent")
		var level slog.Level
		if silent {
			level = slog.LevelInfo
		} else {
			level = slog.LevelDebug
		}
		utils.InitLogger(level)

		handleRsu()
	},
}

func handleRsu() {
	rsuOrder := types.RsuOrder{}

	sellingPrice, err := PromptAndValidate[float64]("What is the selling price per share ($)? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.SellingPricePerShare = sellingPrice

	numberOfShares, err := PromptAndValidate[int]("How many shares sold? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.NumberOfSharesSold = numberOfShares

	considerTransactionCommission, err := PromptAndValidate[bool]("Consider transaction commission[Y/N]? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.ConsiderTransactionCommission = considerTransactionCommission

	if considerTransactionCommission {
		commissionPaidPerTransaction, err := PromptAndValidate[float64]("What is the commission paid per transaction ($)? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		rsuOrder.CommissionPaidPerTransaction = commissionPaidPerTransaction

		numberOfTransactions, err := PromptAndValidate[int]("Number of transactions? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		rsuOrder.NumberOfTransactions = numberOfTransactions
	}

	profitOrLoss := rsuOrder.CalculateEffectiveProfitOrLoss()
	var capitalGainTaxAmount float64
	if profitOrLoss > 0 {
		utils.LogInfo("Profit: $%.2f", profitOrLoss)

		deductCapitalGains, err := PromptAndValidate[bool]("Do you want to calculate capital gain tax and deduct from the profit[Y/N]? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		if deductCapitalGains {
			rsuOrder.ConsiderTransactionCommission = true
			capitalGainTaxPercent, err := PromptAndValidate[float64]("What is the capital gain tax percent (Short-Term: 10%-35%) (Long-Term: 0%-20%)? ")
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			rsuOrder.CapitalGainTaxPercent = capitalGainTaxPercent

			marketPriceOnVestedStockPerShare, err := PromptAndValidate[float64]("What is the (FMV) market price on vested stock per share ($)? ")
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			rsuOrder.MarketValuePerShare = marketPriceOnVestedStockPerShare

			capitalGainTaxableAmount := rsuOrder.CalculateProfitOrLossForCapitalGain()
			if capitalGainTaxableAmount <= 0 {
				utils.LogInfo("Sold at a loss ($%.2f). No Capital Gain.", capitalGainTaxableAmount)
			} else {
				capitalGainTaxAmount, _ = rsuOrder.CalculateCapitalGainTaxAmount(capitalGainTaxableAmount)
				utils.LogInfo("Capital Gain tax amount: $%.2f", capitalGainTaxAmount)
				effectiveProfit := profitOrLoss - capitalGainTaxAmount
				utils.LogInfo("Effective profit: $%.2f", effectiveProfit)
			}
		} else if profitOrLoss < 0 {
			utils.LogInfo("Loss: $%.2f", profitOrLoss)
		} else {
			utils.LogInfo("Broke even: $%.2f", profitOrLoss)
		}

		considerIncomeTaxOnVestedStock, err := PromptAndValidate[bool]("Calculate and deduct income tax on vested stock[Y/N]? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		if !considerIncomeTaxOnVestedStock {
			return
		}

		rsuOrder.ConsiderIncomeTaxOnVestedStock = true
		incomeTaxIncurredWhenStockVested, err := PromptAndValidate[float64]("What is the income tax paid on vested stock\n(no. of shares traded * income tax %) ($)? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		rsuOrder.IncomeTaxIncurredWhenStockVested = incomeTaxIncurredWhenStockVested

		var noOfStocksVested int
		for {
			noOfStocksVested, err = PromptAndValidate[int]("Number of stocks vested? ")
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			if noOfStocksVested <= 0 {
				utils.LogWarn("Number of stocks vested must be greater than 0")
			} else {
				break
			}
		}
		rsuOrder.NumberOfStocksVested = noOfStocksVested

		incomeTaxPerShare, _ := rsuOrder.CalculateIncomeTaxPerShare()
		utils.LogInfo("Income tax per share: $%.2f", incomeTaxPerShare)

		totalIncomeTaxIncurred, err := rsuOrder.CalculateTotalIncomeTaxAmount()
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		utils.LogInfo("Total Income Tax: $%.2f", totalIncomeTaxIncurred)

		effectiveProfitOrLoss := profitOrLoss - totalIncomeTaxIncurred
		if deductCapitalGains && capitalGainTaxAmount > 0 {
			effectiveProfitOrLoss -= capitalGainTaxAmount
		}
		utils.LogInfo("True profit/loss: $%.2f", effectiveProfitOrLoss)
	}
}

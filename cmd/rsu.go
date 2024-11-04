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

	sellingPrice, err := PromptAndValidate[float64]("What is the selling price per share? ")
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

	considerTransactionCommission, err := PromptAndValidate[bool]("Consider transaction commission? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.ConsiderTransactionCommission = considerTransactionCommission

	if considerTransactionCommission {
		commissionPaidPerTransaction, err := PromptAndValidate[float64]("What is the commission paid per transaction? ")
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

	profitOrLoss := rsuOrder.CalculateProfitOrLoss()
	if profitOrLoss > 0 {
		utils.LogInfo("Profit: $%.2f", profitOrLoss)

		deductCapitalGains, err := PromptAndValidate[bool]("Do you want to calculate capital gain tax and deduct from the profit? ")
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

			capitalGainTaxAmount, err := rsuOrder.CalculateCapitalGainTaxAmount(profitOrLoss)
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			utils.LogInfo("Capital Gain Amount: $%.2f", capitalGainTaxAmount)

			effectiveProfit := profitOrLoss - capitalGainTaxAmount
			utils.LogInfo("Effective profit: $%.2f", effectiveProfit)
		}

	} else if profitOrLoss < 0 {
		utils.LogInfo("Loss: $%.2f", profitOrLoss)
	} else {
		utils.LogInfo("Broke even: $%.2f", profitOrLoss)
	}

	considerIncomeTaxOnVestedStock, err := PromptAndValidate[bool]("Calculate and deduct income tax on vested stock? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	if !considerIncomeTaxOnVestedStock {
		return
	}

	rsuOrder.ConsiderIncomeTaxOnVestedStock = true

	incomeTaxPercentOnVestedStockPerShare, err := PromptAndValidate[float64]("What is the income tax percent on vested stock? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.IncomeTaxPercentOnVestedStockPerShare = incomeTaxPercentOnVestedStockPerShare

	marketPriceOnVestedStockPerShare, err := PromptAndValidate[float64]("What is the market price on vested stock per share? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	rsuOrder.MarketPriceOnVestedStockPerShare = marketPriceOnVestedStockPerShare

	totalIncomeTaxIncurred := rsuOrder.CalculateEffectiveIncomeTaxAmount()
	utils.LogInfo("Total Income Tax: $%.2f", totalIncomeTaxIncurred)

	effectiveProfitOrLoss := profitOrLoss - totalIncomeTaxIncurred
	utils.LogInfo("Effective capital gain/loss: $%.2f", effectiveProfitOrLoss)
}

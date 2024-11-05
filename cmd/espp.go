package cmd

import (
	"github.com/leogps/lunar/pkg/types"
	"github.com/leogps/lunar/pkg/utils"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	rootCmd.AddCommand(esppCmd)
}

var esppCmd = &cobra.Command{
	Use:   "espp",
	Short: "calculate profit/loss on ESPP orders interactively",
	Long:  `calculate profit/loss on ESPP orders interactively`,
	Run: func(cmd *cobra.Command, _ []string) {
		silent, _ := cmd.Flags().GetBool("silent")
		var level slog.Level
		if silent {
			level = slog.LevelInfo
		} else {
			level = slog.LevelDebug
		}
		utils.InitLogger(level)

		handleEspp()
	},
}

func handleEspp() {
	discountPercent, err := PromptAndValidate[float64]("What is the discounted (buying) price percent per share (%)? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	esppOrder := types.EsppOrder{
		DiscountPercent: discountPercent,
	}

	costPricePerShare, err := PromptAndValidate[float64]("What is the cost price per share (with/without look-back) ($)? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	esppOrder.CostPerShare = costPricePerShare

	discountAmount := esppOrder.CalculateDiscountAmount()
	utils.LogInfo("Discount Amount: $%.2f", discountAmount)
	effectiveCostPerShare := esppOrder.CalculateEffectiveCostPerShare()
	utils.LogInfo("Effective Cost per share: $%.2f", effectiveCostPerShare)

	sellingPrice, err := PromptAndValidate[float64]("What is the selling price per share ($)? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	esppOrder.SellingPricePerShare = sellingPrice

	numberOfShares, err := PromptAndValidate[int]("How many shares sold? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	esppOrder.NumberOfSharesSold = numberOfShares

	considerTransactionCommission, err := PromptAndValidate[bool]("Deduct transaction commission[Y/N]? ")
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
	esppOrder.ConsiderTransactionCommission = considerTransactionCommission

	if considerTransactionCommission {
		commissionPaidPerTransaction, err := PromptAndValidate[float64]("What is the commission paid per transaction ($)? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		esppOrder.CommissionPaidPerTransaction = commissionPaidPerTransaction

		numberOfTransactions, err := PromptAndValidate[int]("Number of transactions? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		esppOrder.NumberOfTransactions = numberOfTransactions
	}

	profitOrLoss := esppOrder.CalculateProfitOrLoss()
	if profitOrLoss < 0 {
		utils.LogInfo("Loss: $%.2f", profitOrLoss)
		return
	} else if profitOrLoss == 0 {
		utils.LogInfo("Broke even: $%.2f", profitOrLoss)
		return
	}
	if profitOrLoss > 0 {
		utils.LogInfo("Profit: $%.2f", profitOrLoss)

		deductCapitalGains, err := PromptAndValidate[bool]("Do you want to calculate capital gain tax and deduct from the profit[Y/N]? ")
		if err != nil {
			utils.LogError("error occurred", err)
			os.Exit(1)
		}
		if deductCapitalGains {
			esppOrder.ConsiderTransactionCommission = true
			capitalGainTaxPercent, err := PromptAndValidate[float64]("What is the capital gain tax percent (Short-Term: 10%-35%) (Long-Term: 0%-20%)? ")
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			esppOrder.CapitalGainTaxPercent = capitalGainTaxPercent

			capitalGainTaxAmount, err := esppOrder.CalculateCapitalGainTaxAmount(profitOrLoss)
			if err != nil {
				utils.LogError("error occurred", err)
				os.Exit(1)
			}
			utils.LogInfo("Capital Gain Amount: $%.2f", capitalGainTaxAmount)

			effectiveProfit := profitOrLoss - capitalGainTaxAmount
			utils.LogInfo("True profit: $%.2f", effectiveProfit)
		}

	}
}

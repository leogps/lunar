## ESPP

---

 	- Prompt for Discounted price percentage?
 		store: discountPercentage

 	- Share Cost Price (With/Without Lookback): Each Share price?
 		store: shareCostPrice

 	- Effective Cost Per Share: 
 		Discount Amount per Share: 
 			discountAmount = (shareCostPrice *  discount)/100
 			effectiveCostPerShare = (shareCostPrice) - discountAmount

	- What is the selling price?
		store: sellingPrice

	- How many shares did you sell?
		store: numOfSharesSold

	- Do you want to include transaction commission? (dollar amount)
		YES?
			store: commission
		- How many transactions?
			store: numOfTransactions
	    profitOrLossAmount = (numOfSharesSold * (sellingPrice - effectiveCostPerShare)) - (numOfTransactions * transactionComission)

	- If profitOrLossAmount > 0
		Do you want to calculate capital gain tax and deduct from the profit?
		YES?
			Short-Term Capital Gain Tax?
			Long-Term Capital Gain Tax?
			effectiveProfit = profitOrLossAmount - (profitOrLossAmount * capitalGainTax/100)

---

## RSU

---

	- What is the selling price per share?
		store: sellingPrice
	
	- How many shares did you sell?
		store: numOfSharesSold
		
	- Do you want to include transaction commission? (dollar amount)
		YES?
			store: commission
		- How many transactions?
			store: numOfTransactions
	    profitOrLossAmount = (numOfSharesSold * (sellingPrice)) - (numOfTransactions * transactionComission)

	- If profitOrLossAmount > 0
		Do you want to calculate capital gain tax and deduct from the profit?
		YES?
			Short-Term Capital Gain Tax? 20% (default)
			Long-Term Capital Gain Tax? 24% (default)
			effectiveProfit = profitOrLossAmount - (profitOrLossAmount * capitalGainTax/100)

	- Do you want to calculate income tax on vested RSUs and deduct the tax amount?
		YES?
			store: incomeTaxPercent

		- What is the Market Price per share at the time of vesting?
			store: vestingPricePerShare
		profitOrLoss - (numOfSharesSold*vestingPricePerShare * incomeTaxPercent/100)
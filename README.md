# Lunar

---

lunar is a CLI tool to perform calculations for stocks.

```
Usage:
[flags]
[command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
espp        calculate profit/loss on ESPP orders interactively
help        Help about any command
rsu         calculate profit/loss on RSU orders interactively
ui          Starts Terminal UI

Flags:
-h, --help   help for this command
```

---

### ESPP

---

#### Usage

    lunar espp # For interactive 
        OR
    lunar ui # Choose ESPP

#### ESPP Order Summary Calculation

This program calculates and displays a summary of an Employee Stock Purchase Plan (ESPP) order. It is intended to breakdown the costs and profit involved in ESPP to calculate the 'true' profit. 

##### Features

Legend:

1. **Cost per Share**: The base cost of each share.
2. **Discount**: The discount percentage applied to the cost per share as part of ESPP with/without look-back.
3. **True Cost per Share**: The effective cost per share after discount.
4. **Selling Price per Share**: The price at which each share is sold.
5. **Total Selling Price**: The total revenue from selling the shares.
6. **Total Cost**: The total cost incurred for the shares sold.
7. **Effective Commission Fee (Optional)**: Total commission fee, if transaction commissions are included.
8. **Profit or Loss Before Capital Gains Tax**: The net result of the sale before applying capital gains tax.
9. **Capital Gains Tax (Optional)**: Tax on profits if the selling price exceeds the cost price.
10. **Profit or Loss After Capital Gains Tax**: The net result after deducting capital gains tax.
11. **Gain/Loss Margin**: The percentage of gain or loss on the transaction.

#### Target Profit Calculation

In the context of an ESPP, the **Target Profit** represents a specified percentage of profit aimed to achieve from selling shares, calculated relative to the total effective cost. This includes capital gains tax (if applicable) and optional commission fees.

---

### RSU

---

#### Usage

    lunar rsu # For interactive
        OR
    lunar ui # Choose RSU

---

* Optionally, considers Fair Market Value (FMV) at the time of vesting and 'true' profit considered only based on the number of shares traded to cover for income tax. 

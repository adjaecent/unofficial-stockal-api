package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adjaecent/unofficial-stockal-api"
)

func main() {
	// Create client with custom timeout
	client := stockal.NewClient(
		stockal.WithTimeout(60*time.Second),
		stockal.WithUserAgent("example-app/1.0"),
	)

	// Create context with timeout for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example credentials (replace with actual credentials)
	username := "nid90"
	password := "P93WmwB7P7658uy"

	// Attempt to login
	resp, err := client.Login(ctx, username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Printf("Login successful!\n")
	fmt.Printf("Response Code: %d\n", resp.Code)
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Access Token: %s\n", resp.Data.AccessToken)
	fmt.Printf("Token Expires: %s\n", resp.Data.ExpiryAccessToken)

	// Get account summary
	fmt.Printf("\nFetching account summary...\n")
	summary, err := client.GetAccountSummary(ctx)
	if err != nil {
		log.Fatalf("Failed to get account summary: %v", err)
	}

	fmt.Printf("Account Summary Response Code: %d\n", summary.Code)
	fmt.Printf("Message: %s\n", summary.Message)
	fmt.Printf("UTC Time: %s\n", summary.Data.UTCTime)

	// Account details
	fmt.Printf("\n--- Account Details ---\n")
	fmt.Printf("Cash Available for Trade: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
	fmt.Printf("Cash Available for Withdrawal: $%.2f\n", summary.Data.AccountSummary.CashAvailableForWithdrawal)
	fmt.Printf("Cash Balance: $%.2f\n", summary.Data.AccountSummary.CashBalance)
	fmt.Printf("Good Faith Violations: %s\n", summary.Data.AccountSummary.GoodFaithViolations)
	fmt.Printf("Restricted: %t\n", summary.Data.AccountSummary.Restricted)

	// Portfolio summary
	fmt.Printf("\n--- Portfolio Summary ---\n")
	fmt.Printf("Total Current Value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)
	fmt.Printf("Total Investment Amount: $%.2f\n", summary.Data.PortfolioSummary.TotalInvestmentAmount)
	fmt.Printf("Total Gain/Loss: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue-summary.Data.PortfolioSummary.TotalInvestmentAmount)

	// Stock portfolio
	fmt.Printf("Stock Portfolio - Current: $%.2f, Investment: $%.2f\n",
		summary.Data.PortfolioSummary.StockPortfolio.CurrentValue,
		summary.Data.PortfolioSummary.StockPortfolio.InvestmentAmount)

	// ETF portfolio
	fmt.Printf("ETF Portfolio - Current: $%.2f, Investment: $%.2f\n",
		summary.Data.PortfolioSummary.ETFPortfolio.CurrentValue,
		summary.Data.PortfolioSummary.ETFPortfolio.InvestmentAmount)

	// Get portfolio details
	fmt.Printf("\nFetching portfolio details...\n")
	portfolio, err := client.GetPortfolioDetail(ctx)
	if err != nil {
		log.Fatalf("Failed to get portfolio details: %v", err)
	}

	fmt.Printf("Portfolio Detail Response Code: %d\n", portfolio.Code)
	fmt.Printf("Message: %s\n", portfolio.Message)
	fmt.Printf("Total Holdings: %d\n", portfolio.Data.TotalRecords)
	fmt.Printf("Pending Transactions: %d\n", len(portfolio.Data.PendingData))

	// Display holdings
	fmt.Printf("\n--- Holdings ---\n")
	for i, holding := range portfolio.Data.Holdings {
		currentValue := holding.TotalUnit * holding.Price
		gainLoss := currentValue - holding.TotalInvestment
		gainLossPercent := (gainLoss / holding.TotalInvestment) * 100

		fmt.Printf("%d. %s (%s)\n", i+1, holding.Company, holding.Symbol)
		fmt.Printf("   Units: %.4f @ $%.2f = $%.2f\n", holding.TotalUnit, holding.Price, currentValue)
		fmt.Printf("   Investment: $%.2f\n", holding.TotalInvestment)
		fmt.Printf("   Gain/Loss: $%.2f (%.2f%%)\n", gainLoss, gainLossPercent)
		fmt.Printf("   Category: %s | Status: %s\n", holding.Category, holding.Status)
		if holding.SellOnly {
			fmt.Printf("   ** SELL ONLY **\n")
		}
		fmt.Printf("\n")

		// Show only first 5 holdings to avoid too much output
		if i >= 4 {
			fmt.Printf("... and %d more holdings\n", len(portfolio.Data.Holdings)-5)
			break
		}
	}
}

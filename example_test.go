package stockal_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adjaecent/unofficial-stockal-api"
)

// ExampleNewClient demonstrates how to create a new Stockal API client.
func ExampleNewClient() {
	client := stockal.NewClient()
	_ = client // Use the client variable
	fmt.Printf("Client created successfully\n")

	// With custom options
	customClient := stockal.NewClient(
		stockal.WithTimeout(60*time.Second),
		stockal.WithUserAgent("my-app/1.0"),
	)
	_ = customClient // Use the customClient variable
	fmt.Printf("Custom client created successfully\n")
	// Output: Client created successfully
	// Custom client created successfully
}

// ExampleClient_Login demonstrates user authentication with the Stockal API.
func ExampleClient_Login() {
	client := stockal.NewClient()
	ctx := context.Background()

	// Note: Use actual credentials in real usage
	resp, err := client.Login(ctx, "your_username", "your_password")
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	fmt.Printf("Login successful! Response code: %d\n", resp.Code)
	fmt.Printf("Token expires: %s\n", resp.Data.ExpiryAccessToken)
}

// ExampleClient_GetAccountSummary demonstrates fetching account summary information.
func ExampleClient_GetAccountSummary() {
	client := stockal.NewClient()
	ctx := context.Background()

	// Login first (credentials would be real in actual usage)
	_, err := client.Login(ctx, "your_username", "your_password")
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	// Get account summary
	summary, err := client.GetAccountSummary(ctx)
	if err != nil {
		log.Printf("Failed to get account summary: %v", err)
		return
	}

	fmt.Printf("Cash available for trade: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
	fmt.Printf("Total portfolio value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)
	fmt.Printf("Total investment: $%.2f\n", summary.Data.PortfolioSummary.TotalInvestmentAmount)

	// Calculate total gain/loss
	gainLoss := summary.Data.PortfolioSummary.TotalCurrentValue - summary.Data.PortfolioSummary.TotalInvestmentAmount
	fmt.Printf("Total gain/loss: $%.2f\n", gainLoss)
}

// ExampleClient_GetPortfolioDetail demonstrates fetching detailed portfolio information.
func ExampleClient_GetPortfolioDetail() {
	client := stockal.NewClient()
	ctx := context.Background()

	// Login first (credentials would be real in actual usage)
	_, err := client.Login(ctx, "your_username", "your_password")
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	// Get portfolio details
	portfolio, err := client.GetPortfolioDetail(ctx)
	if err != nil {
		log.Printf("Failed to get portfolio details: %v", err)
		return
	}

	fmt.Printf("Total holdings: %d\n", portfolio.Data.TotalRecords)
	fmt.Printf("Pending transactions: %d\n", len(portfolio.Data.PendingData))

	// Show details for first few holdings
	for i, holding := range portfolio.Data.Holdings {
		if i >= 3 { // Limit output for example
			break
		}

		currentValue := holding.TotalUnit * holding.Price
		gainLoss := currentValue - holding.TotalInvestment
		gainLossPercent := (gainLoss / holding.TotalInvestment) * 100

		fmt.Printf("\n%s (%s):\n", holding.Company, holding.Symbol)
		fmt.Printf("  Units: %.4f @ $%.2f = $%.2f\n", holding.TotalUnit, holding.Price, currentValue)
		fmt.Printf("  Investment: $%.2f\n", holding.TotalInvestment)
		fmt.Printf("  Gain/Loss: $%.2f (%.2f%%)\n", gainLoss, gainLossPercent)

		if holding.SellOnly {
			fmt.Printf("  Status: SELL ONLY\n")
		}
	}
}

// Example demonstrates a complete workflow: login, get account summary, and portfolio details.
func Example() {
	// Create a new client with custom options
	client := stockal.NewClient(
		stockal.WithTimeout(30*time.Second),
		stockal.WithUserAgent("example-workflow/1.0"),
	)
	ctx := context.Background()

	// Step 1: Login
	resp, err := client.Login(ctx, "your_username", "your_password")
	if err != nil {
		log.Fatal("Login failed:", err)
	}
	fmt.Printf("✓ Login successful (token expires: %s)\n", resp.Data.ExpiryAccessToken)

	// Step 2: Get account summary
	summary, err := client.GetAccountSummary(ctx)
	if err != nil {
		log.Fatal("Failed to get account summary:", err)
	}

	fmt.Printf("✓ Account Summary:\n")
	fmt.Printf("  Cash available: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
	fmt.Printf("  Portfolio value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)

	// Step 3: Get portfolio details
	portfolio, err := client.GetPortfolioDetail(ctx)
	if err != nil {
		log.Fatal("Failed to get portfolio details:", err)
	}

	fmt.Printf("✓ Portfolio Details:\n")
	fmt.Printf("  Total holdings: %d\n", portfolio.Data.TotalRecords)

	// Calculate and show top performing stocks
	var bestPerformer stockal.Holding
	var bestPerformance float64

	for _, holding := range portfolio.Data.Holdings {
		if holding.TotalInvestment > 0 {
			currentValue := holding.TotalUnit * holding.Price
			performance := ((currentValue - holding.TotalInvestment) / holding.TotalInvestment) * 100

			if performance > bestPerformance {
				bestPerformance = performance
				bestPerformer = holding
			}
		}
	}

	if bestPerformer.Symbol != "" {
		fmt.Printf("  Best performer: %s (%.2f%% gain)\n", bestPerformer.Symbol, bestPerformance)
	}
}

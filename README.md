# Unofficial Stockal API Go Library (Read-Only)

[![Go Reference](https://pkg.go.dev/badge/github.com/adjaecent/unofficial-stockal-api.svg)](https://pkg.go.dev/github.com/adjaecent/unofficial-stockal-api)

An unofficial **read-only** Go client library for the [Stockal](https://globalinvesting.in/) trading platform API. This library provides a simple and clean interface to **view and analyze** your Stockal account data, portfolio holdings, and performance metrics.

> **🔍 Read-Only Library**: This library is designed for **data retrieval and analysis only**. It cannot place trades, modify orders, or perform any account changes. All trading operations must be done through the official Stockal platform.

> **⚠️ Disclaimer**: This is an unofficial library and is not affiliated with or endorsed by Stockal. Use at your own risk and always verify account information through the official Stockal platform.

## 🚀 Features (Read-Only Data Access)

- ✅ **User Authentication** - Login with username/password to get access tokens
- ✅ **Account Summary** - View cash balances, restrictions, and portfolio totals
- ✅ **Portfolio Analysis** - Analyze detailed holdings with real-time prices and P&L

## 📦 Installation

```bash
go get github.com/adjaecent/unofficial-stockal-api
```

## 🏃‍♂️ Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/adjaecent/unofficial-stockal-api"
)

func main() {
    // Create a new client
    client := stockal.NewClient()

    // Login to get access token
    resp, err := client.Login("your_username", "your_password")
    if err != nil {
        log.Fatal("Login failed:", err)
    }
    fmt.Printf("✓ Login successful! Token expires: %s\n", resp.Data.ExpiryAccessToken)

    // Retrieve account summary (read-only)
    summary, err := client.GetAccountSummary()
    if err != nil {
        log.Fatal("Failed to get account summary:", err)
    }

    fmt.Printf("💰 Cash available: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
    fmt.Printf("📈 Portfolio value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)

    // Analyze detailed portfolio (read-only)
    portfolio, err := client.GetPortfolioDetail()
    if err != nil {
        log.Fatal("Failed to get portfolio details:", err)
    }

    fmt.Printf("📊 Total holdings: %d\n", portfolio.Data.TotalRecords)

    // Analyze top 3 holdings performance
    for i, holding := range portfolio.Data.Holdings {
        if i >= 3 { break }

        currentValue := holding.TotalUnit * holding.Price
        gainLoss := currentValue - holding.TotalInvestment
        gainLossPercent := (gainLoss / holding.TotalInvestment) * 100

        fmt.Printf("🏢 %s (%s): $%.2f (%.2f%%)\n",
            holding.Company, holding.Symbol, currentValue, gainLossPercent)
    }
}
```

## 🔧 Examples

### Running the Examples

#### 1. **Working Example** (`example/main.go`)
- **Purpose**: Demonstrates real API usage with actual credentials
- **Usage**: Update credentials in the file, then run:
  ```bash
  go run example/main.go
  ```

#### 2. **Documentation Examples** (`example_test.go`)
- **Purpose**: Provides godoc examples and API documentation
- **Usage**: View in documentation or run tests:
  ```bash
  go test -run Example
  godoc -http=:6060  # View at http://localhost:6060
  ```

## 📖 Documentation

- **Go Documentation**: Run `godoc -http=:6060` and visit http://localhost:6060
- **API Reference**: See the godoc examples in `example_test.go`
- **OpenAPI Spec**: Complete specification available in `openapi.yaml`

## 🔒 Security Considerations

- **Never commit credentials** to version control
- **Use environment variables** for sensitive information:
  ```go
  username := os.Getenv("STOCKAL_USERNAME")
  password := os.Getenv("STOCKAL_PASSWORD")
  ```
- **Validate all responses** before using data for trading decisions
- **Test thoroughly** with small amounts before scaling

## 🛠️ Development

### Prerequisites

- Go 1.25.0+ (see `.tool-versions`)
- Valid Stockal account credentials

### Building

```bash
# Build the library
go build

# Run tests
go test

# Generate documentation
godoc -http=:6060

# Run example
go run example/main.go
```

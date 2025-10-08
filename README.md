# Unofficial Stockal API Go Library (Read-Only)

[![Go Reference](https://pkg.go.dev/badge/github.com/kitallis/unofficial-stockal-api.svg)](https://pkg.go.dev/github.com/kitallis/unofficial-stockal-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/kitallis/unofficial-stockal-api)](https://goreportcard.com/report/github.com/kitallis/unofficial-stockal-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An unofficial **read-only** Go client library for the [Stockal](https://globalinvesting.in/) trading platform API. This library provides a simple and clean interface to **view and analyze** your Stockal account data, portfolio holdings, and performance metrics.

> **🔍 Read-Only Library**: This library is designed for **data retrieval and analysis only**. It cannot place trades, modify orders, or perform any account changes. All trading operations must be done through the official Stockal platform.

> **⚠️ Disclaimer**: This is an unofficial library and is not affiliated with or endorsed by Stockal. Use at your own risk and always verify account information through the official Stockal platform.

## 🚀 Features (Read-Only Data Access)

- ✅ **User Authentication** - Login with username/password to get access tokens
- ✅ **Account Summary** - View cash balances, restrictions, and portfolio totals  
- ✅ **Portfolio Analysis** - Analyze detailed holdings with real-time prices and P&L
- ✅ **Performance Metrics** - Calculate gains/losses and portfolio performance
- ✅ **Comprehensive Documentation** - Full godoc documentation with examples
- ✅ **Type Safety** - Complete Go structs for all API responses
- ✅ **Error Handling** - Detailed error messages for network and API issues
- ✅ **Zero Dependencies** - Uses only Go standard library
- ✅ **OpenAPI Specification** - Complete OpenAPI 3.0 spec included

### 🚫 What This Library Cannot Do

- ❌ **Place trades** or execute buy/sell orders
- ❌ **Modify account settings** or personal information  
- ❌ **Cancel or update orders**
- ❌ **Transfer funds** or withdraw money
- ❌ **Access real-time market data** (only portfolio data)
- ❌ **Perform any account modifications**

> **💡 Use Case**: This library is perfect for **portfolio tracking**, **performance analysis**, **tax reporting**, **personal finance apps**, and **investment research tools**.

## 📦 Installation

```bash
go get github.com/kitallis/unofficial-stockal-api
```

## 🏃‍♂️ Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kitallis/unofficial-stockal-api"
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

## 📚 API Reference

### Authentication

#### `Login(username, password string) (*LoginResponse, error)`

Authenticate with Stockal and retrieve access token.

```go
client := stockal.NewClient()
resp, err := client.Login("username", "password")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Access token: %s\n", resp.Data.AccessToken)
fmt.Printf("Expires: %s\n", resp.Data.ExpiryAccessToken)
```

### Account Information

#### `GetAccountSummary() (*AccountSummaryResponse, error)`

Retrieve comprehensive account summary including cash balances and portfolio totals (read-only).

```go
summary, err := client.GetAccountSummary()
if err != nil {
    log.Fatal(err)
}

// Cash information
fmt.Printf("Cash for trading: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
fmt.Printf("Cash for withdrawal: $%.2f\n", summary.Data.AccountSummary.CashAvailableForWithdrawal)

// Portfolio totals
fmt.Printf("Total value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)
fmt.Printf("Total invested: $%.2f\n", summary.Data.PortfolioSummary.TotalInvestmentAmount)
```

### Portfolio Analysis

#### `GetPortfolioDetail() (*PortfolioDetailResponse, error)`

Get detailed information about all holdings including current prices and performance (read-only analysis).

```go
portfolio, err := client.GetPortfolioDetail()
if err != nil {
    log.Fatal(err)
}

for _, holding := range portfolio.Data.Holdings {
    currentValue := holding.TotalUnit * holding.Price
    gainLoss := currentValue - holding.TotalInvestment
    gainLossPercent := (gainLoss / holding.TotalInvestment) * 100
    
    fmt.Printf("%s: %.4f shares @ $%.2f = $%.2f (%.2f%%)\n",
        holding.Symbol, holding.TotalUnit, holding.Price, currentValue, gainLossPercent)
}
```

## 🏗️ Project Structure

```
unofficial-stockal-api/
├── README.md              # This file
├── go.mod                 # Go module definition
├── .tool-versions         # Go version specification
├── stockal.go             # Main library code
├── example_test.go        # Godoc examples
├── openapi.yaml           # OpenAPI 3.0 specification
└── example/
    └── main.go            # Working example program
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

### Sample Output

```
Login successful!
Response Code: 200
Message: Success
Access Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Token Expires: 2025-10-08T08:04:10.681Z

Fetching account summary...
Account Summary Response Code: 200
Message: Success

--- Account Details ---
Cash Available for Trade: $1459.50
Cash Available for Withdrawal: $1459.50
Cash Balance: $1459.50
Good Faith Violations: 0 of 3
Restricted: false

--- Portfolio Summary ---
Total Current Value: $39698.25
Total Investment Amount: $26214.85
Total Gain/Loss: $13483.40

--- Holdings ---
1. Advanced Micro Devices, Inc. (AMD)
   Units: 17.0000 @ $211.51 = $3595.67
   Investment: $2401.60
   Gain/Loss: $1194.07 (49.72%)
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

### Code Quality

- **Comprehensive tests**: Includes example tests and godoc examples
- **Full documentation**: Every public function and type is documented
- **Clean architecture**: Minimal dependencies, clear separation of concerns
- **Error handling**: Detailed error messages with context

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository** and create a feature branch
2. **Add tests** for any new functionality
3. **Update documentation** including godoc comments
4. **Follow Go conventions** and run `gofmt`
5. **Submit a pull request** with a clear description

### Areas for Contribution

- [ ] Additional API endpoints (trading, orders, etc.)
- [ ] Rate limiting and retry logic
- [ ] WebSocket support for real-time data
- [ ] CLI tool built on this library
- [ ] Integration tests
- [ ] Performance optimizations

## 📋 API Coverage

### ✅ Implemented Endpoints

- `POST /v3/auth/login` - User authentication
- `GET /v2/users/accountSummary/summary` - Account summary
- `GET /v2/users/portfolio/detail` - Portfolio details

### 🚧 Potential Future Read-Only Endpoints

- Transaction history and trade records
- Watchlist viewing
- Performance analytics and reporting
- Account statements and documents
- Historical portfolio snapshots

> **Note**: This library will remain read-only. Trading functionality is intentionally excluded for safety and security reasons.

## ⚠️ Known Limitations

- **Read-only operations only**: This library is designed exclusively for data retrieval and analysis
- **No trading capabilities**: Cannot place orders, execute trades, or modify account settings
- **Limited to portfolio data**: Does not provide real-time market quotes or market data
- **Rate limiting**: No built-in rate limiting (be respectful with API calls)
- **Error recovery**: Limited retry logic for failed requests
- **WebSocket**: No real-time data streaming support

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Stockal team for providing the trading platform
- Go community for excellent HTTP and JSON libraries
- Contributors and users of this library

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/kitallis/unofficial-stockal-api/issues)
- **Discussions**: [GitHub Discussions](https://github.com/kitallis/unofficial-stockal-api/discussions)
- **Documentation**: Available via `godoc` and inline comments

---

**Disclaimer**: This is an unofficial library and is not affiliated with or endorsed by Stockal. Always verify your trades and account information through the official Stockal platform. Trading involves risk, and past performance does not guarantee future results.
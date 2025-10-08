// Package stockal provides an unofficial Go client library for the Stockal API.
//
// Stockal is a platform that allows trading in the US stock market. This library
// provides a simple interface to interact with Stockal's REST API for authentication,
// account management, and portfolio operations.
//
// # Basic Usage
//
//	client := stockal.NewClient()
//
//	// Login to get access token
//	resp, err := client.Login("username", "password")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get account summary
//	summary, err := client.GetAccountSummary()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get portfolio details
//	portfolio, err := client.GetPortfolioDetail()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Authentication
//
// All API calls except Login require authentication. The Client automatically
// stores and includes the access token from a successful login in subsequent
// requests.
//
// # Error Handling
//
// All methods return detailed error information. Network errors, JSON parsing
// errors, and API errors are wrapped with descriptive messages.
package stockal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BaseURL is the base URL for the Stockal API v2.
const (
	BaseURL = "https://api-v2.stockal.com"
	DefaultTimeout = 30 * time.Second
	DefaultUserAgent = "unofficial-stockal-api/1.0"
)

// Common errors
var (
	ErrNotAuthenticated = errors.New("not authenticated: please login first")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyPassword = errors.New("password cannot be empty")
)

// APIError represents an error response from the Stockal API.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	if e.Err != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// StockalClient defines the interface for Stockal API operations.
type StockalClient interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	GetAccountSummary(ctx context.Context) (*AccountSummaryResponse, error)
	GetPortfolioDetail(ctx context.Context) (*PortfolioDetailResponse, error)
}

// ClientOption is a function that configures a Client.
type ClientOption func(*clientConfig)

// clientConfig holds configuration for the client.
type clientConfig struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *clientConfig) {
		c.userAgent = userAgent
	}
}

// WithTimeout sets a custom timeout for requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// Client represents a Stockal API client with authentication and HTTP configuration.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	userAgent   string
	accessToken string
}

// LoginRequest represents the request payload for user authentication.
type LoginRequest struct {
	// Username is the user's login username
	Username string `json:"username"`
	// Password is the user's login password
	Password string `json:"password"`
}

// LoginData represents the data payload of a login response.
type LoginData struct {
	// AccessToken is the JWT token used for authenticated API calls
	AccessToken           string `json:"accessToken"`
	// RefreshToken is used to refresh the access token when it expires
	RefreshToken          string `json:"refreshToken"`
	// ExpiryAccessToken is the access token expiration time
	ExpiryAccessToken     string `json:"expiryAccessToken"`
	// ExpiryRefreshToken is the refresh token expiration time
	ExpiryRefreshToken    string `json:"expiryRefreshToken"`
}

// LoginResponse represents the response from the login API endpoint.
type LoginResponse struct {
	// Code is the HTTP response code
	Code    int       `json:"code"`
	// Message contains the response message (usually "Success")
	Message string    `json:"message"`
	// Data contains the actual login data with tokens
	Data    LoginData `json:"data,omitempty"`
	// Error contains error code if login failed
	Error   string    `json:"error,omitempty"`
}

// CashSettlement represents a scheduled cash settlement in the account.
type CashSettlement struct {
	// UTCTime is the settlement date and time in UTC
	UTCTime string  `json:"utcTime"`
	// Cash is the amount to be settled
	Cash    float64 `json:"cash"`
}

// AccountSummary represents the user's account summary including cash balances and restrictions.
type AccountSummary struct {
	// CashAvailableForTrade is the cash available for placing new trades
	CashAvailableForTrade       float64          `json:"cashAvailableForTrade"`
	// CashAvailableForWithdrawal is the cash available for withdrawal
	CashAvailableForWithdrawal  float64          `json:"cashAvailableForWithdrawal"`
	// CashBalance is the total cash balance in the account
	CashBalance                 float64          `json:"cashBalance"`
	// GoodFaithViolations shows current violations count (e.g., "0 of 3")
	GoodFaithViolations         string           `json:"goodFaithViolations"`
	// Restricted indicates if the account has trading restrictions
	Restricted                  bool             `json:"restricted"`
	// CashSettlement contains scheduled cash settlements
	CashSettlement              []CashSettlement `json:"cashSettlement"`
}

// Portfolio represents a portfolio category with current value and investment amount.
type Portfolio struct {
	// CurrentValue is the current market value of holdings in this portfolio
	CurrentValue     float64 `json:"currentValue"`
	// InvestmentAmount is the total amount invested in this portfolio
	InvestmentAmount float64 `json:"investmentAmount"`
}

// PortfolioSummary represents a summary of all portfolio categories and totals.
type PortfolioSummary struct {
	// StockPortfolio contains stock holdings summary
	StockPortfolio         Portfolio `json:"stockPortfolio"`
	// StackPortfolio contains stack holdings summary
	StackPortfolio         Portfolio `json:"stackPortfolio"`
	// ETFPortfolio contains ETF holdings summary
	ETFPortfolio          Portfolio `json:"etfPortfolio"`
	// TotalCurrentValue is the total current value across all portfolios
	TotalCurrentValue     float64   `json:"totalCurrentValue"`
	// TotalInvestmentAmount is the total invested amount across all portfolios
	TotalInvestmentAmount float64   `json:"totalInvestmentAmount"`
}

// AccountSummaryData represents the data payload of an account summary response.
type AccountSummaryData struct {
	// UTCTime is the timestamp when the summary was generated
	UTCTime          string           `json:"utcTime"`
	// AccountSummary contains account-level information
	AccountSummary   AccountSummary   `json:"accountSummary"`
	// UnsettledAmount is the amount of unsettled funds
	UnsettledAmount  float64          `json:"unsettledAmount"`
	// PortfolioSummary contains portfolio-level summaries
	PortfolioSummary PortfolioSummary `json:"portfolioSummary"`
}

// AccountSummaryResponse represents the complete response from the account summary API.
type AccountSummaryResponse struct {
	// Code is the HTTP response code
	Code    int                `json:"code"`
	// Message is the response message (usually "Success")
	Message string             `json:"message"`
	// Data contains the actual account summary data
	Data    AccountSummaryData `json:"data"`
}

// Holding represents a single stock or asset holding in the portfolio.
type Holding struct {
	// Symbol is the stock symbol (e.g., "AAPL")
	Symbol           string  `json:"symbol"`
	// Ticker is the trading ticker symbol
	Ticker           string  `json:"ticker"`
	// UserID is the user's unique identifier
	UserID           string  `json:"userID"`
	// Date is the last update date for this holding
	Date             string  `json:"Date"`
	// V is the version field from MongoDB
	V                int     `json:"__v"`
	// Category is the asset category (e.g., "stock")
	Category         string  `json:"category"`
	// Status is the holding status (e.g., "successful")
	Status           string  `json:"status"`
	// Timestamp is the Unix timestamp of the last update
	Timestamp        int64   `json:"timestamp"`
	// TotalInvestment is the total amount invested in this holding
	TotalInvestment  float64 `json:"totalInvestment"`
	// TotalUnit is the number of shares/units owned
	TotalUnit        float64 `json:"totalUnit"`
	// Type is the asset type (e.g., "stock")
	Type             string  `json:"type"`
	// Code is the asset code
	Code             string  `json:"code"`
	// Company is the full company name
	Company          string  `json:"company"`
	// Price is the current price per share
	Price            float64 `json:"price"`
	// Listed indicates if the asset is currently listed/tradeable
	Listed           bool    `json:"listed"`
	// Close is the current closing price
	Close            float64 `json:"close"`
	// PriorClose is the previous day's closing price
	PriorClose       float64 `json:"priorClose"`
	// Logo is the URL to the company's logo image (optional)
	Logo             string  `json:"logo,omitempty"`
	// SellOnly indicates if only sell operations are allowed (optional)
	SellOnly         bool    `json:"sellOnly,omitempty"`
}

// PortfolioDetailData represents the data payload of a portfolio detail response.
type PortfolioDetailData struct {
	// PendingData contains any pending transactions (usually empty)
	PendingData  []interface{} `json:"pendingData"`
	// Holdings contains all current holdings in the portfolio
	Holdings     []Holding     `json:"holdings"`
	// Timestamp is the Unix timestamp when the data was generated
	Timestamp    int64         `json:"timestamp"`
	// TotalRecords is the total number of holdings
	TotalRecords int           `json:"totalRecords"`
}

// PortfolioDetailResponse represents the complete response from the portfolio detail API.
type PortfolioDetailResponse struct {
	// Code is the HTTP response code
	Code    int                 `json:"code"`
	// Message is the response message (usually "Success")
	Message string              `json:"message"`
	// Data contains the actual portfolio detail data
	Data    PortfolioDetailData `json:"data"`
}

// NewClient creates a new Stockal API client with the given options.
//
// Default configuration:
//   - BaseURL: Official Stockal API endpoint
//   - Timeout: 30 seconds
//   - UserAgent: unofficial-stockal-api/1.0
//
// Example:
//
//	client := stockal.NewClient()
//	resp, err := client.Login(ctx, "username", "password")
//
// With custom options:
//
//	client := stockal.NewClient(
//		stockal.WithTimeout(60*time.Second),
//		stockal.WithUserAgent("my-app/1.0"),
//	)
func NewClient(options ...ClientOption) StockalClient {
	config := &clientConfig{
		baseURL:   BaseURL,
		userAgent: DefaultUserAgent,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, option := range options {
		option(config)
	}

	return &Client{
		baseURL:    config.baseURL,
		httpClient: config.httpClient,
		userAgent:  config.userAgent,
	}
}

// makeRequest is an internal helper method that handles HTTP request creation and execution.
// It automatically adds all necessary headers including authentication and browser simulation.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) (*http.Response, error) {
	// Validate URL
	apiURL, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers similar to browser request
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// Removed Accept-Encoding to avoid compression issues
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", "https://globalinvesting.in")
	req.Header.Set("Referer", "https://globalinvesting.in/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	// Add Authorization header if access token is available
	if c.accessToken != "" {
		req.Header.Set("Authorization", c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// handleResponse is an internal helper method that processes HTTP responses.
// It handles response body reading, JSON unmarshaling, and status code validation.
func (c *Client) handleResponse(resp *http.Response, result interface{}, operation string) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to parse as JSON first
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		// Try to extract API error details
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Code != 0 {
			return &apiErr
		}
		return fmt.Errorf("%s failed with status code: %d", operation, resp.StatusCode)
	}

	return nil
}

// Login authenticates the user with Stockal and stores the access token for subsequent requests.
//
// The method sends a POST request to the authentication endpoint with the provided
// credentials. On successful authentication, the access token is automatically stored
// in the client and will be included in all subsequent API calls.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - username: The user's Stockal username
//   - password: The user's Stockal password
//
// Returns:
//   - LoginResponse: Contains access token, expiration info, and any error details
//   - error: Any network, parsing, authentication, or validation errors
//
// Example:
//
//	ctx := context.Background()
//	client := stockal.NewClient()
//	resp, err := client.Login(ctx, "myusername", "mypassword")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Logged in successfully. Token expires: %s\n", resp.Data.ExpiryAccessToken)
func (c *Client) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	// Input validation
	if strings.TrimSpace(username) == "" {
		return nil, ErrEmptyUsername
	}
	if strings.TrimSpace(password) == "" {
		return nil, ErrEmptyPassword
	}

	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.makeRequest(ctx, "POST", "/v3/auth/login", loginReq)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}

	var loginResp LoginResponse

	// Handle response parsing
	if err := c.handleResponse(resp, &loginResp, "login"); err != nil {
		// Check for specific login error messages
		if loginResp.Error != "" {
			return &loginResp, ErrInvalidCredentials
		}
		return &loginResp, err
	}

	// Store access token in client for subsequent requests
	c.accessToken = loginResp.Data.AccessToken

	return &loginResp, nil
}

// GetAccountSummary retrieves a comprehensive summary of the user's account.
//
// This method fetches account-level information including cash balances, trading restrictions,
// portfolio summaries, and unsettled amounts. The user must be authenticated (logged in)
// before calling this method.
//
// Returns:
//   - AccountSummaryResponse: Contains account details, cash balances, and portfolio summaries
//   - error: Authentication errors, network errors, or API errors
//
// The response includes:
//   - Cash available for trading and withdrawal
//   - Account restrictions and violations
//   - Portfolio summaries by asset type (stocks, ETFs)
//   - Total portfolio values and investment amounts
//
// Example:
//
//	summary, err := client.GetAccountSummary()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Cash available: $%.2f\n", summary.Data.AccountSummary.CashAvailableForTrade)
//	fmt.Printf("Total portfolio value: $%.2f\n", summary.Data.PortfolioSummary.TotalCurrentValue)
func (c *Client) GetAccountSummary(ctx context.Context) (*AccountSummaryResponse, error) {
	if c.accessToken == "" {
		return nil, ErrNotAuthenticated
	}

	resp, err := c.makeRequest(ctx, "GET", "/v2/users/accountSummary/summary", nil)
	if err != nil {
		return nil, fmt.Errorf("account summary request failed: %w", err)
	}

	var summaryResp AccountSummaryResponse
	if err := c.handleResponse(resp, &summaryResp, "account summary"); err != nil {
		return &summaryResp, err
	}

	return &summaryResp, nil
}

// GetPortfolioDetail retrieves detailed information about all holdings in the user's portfolio.
//
// This method fetches comprehensive details for each individual holding including current prices,
// investment amounts, units owned, gain/loss information, and trading restrictions. The user
// must be authenticated (logged in) before calling this method.
//
// Returns:
//   - PortfolioDetailResponse: Contains detailed holdings information and portfolio metadata
//   - error: Authentication errors, network errors, or API errors
//
// Each holding includes:
//   - Stock symbol, company name, and current price
//   - Number of units owned and total investment amount
//   - Current market value and performance vs investment
//   - Trading status and any restrictions (e.g., sell-only)
//   - Asset category (stock, ETF, etc.)
//
// Example:
//
//	portfolio, err := client.GetPortfolioDetail()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Total holdings: %d\n", portfolio.Data.TotalRecords)
//	for _, holding := range portfolio.Data.Holdings {
//		currentValue := holding.TotalUnit * holding.Price
//		gainLoss := currentValue - holding.TotalInvestment
//		fmt.Printf("%s: $%.2f (%.2f%% gain/loss)\n",
//			holding.Symbol, currentValue, (gainLoss/holding.TotalInvestment)*100)
//	}
func (c *Client) GetPortfolioDetail(ctx context.Context) (*PortfolioDetailResponse, error) {
	if c.accessToken == "" {
		return nil, ErrNotAuthenticated
	}

	resp, err := c.makeRequest(ctx, "GET", "/v2/users/portfolio/detail", nil)
	if err != nil {
		return nil, fmt.Errorf("portfolio detail request failed: %w", err)
	}

	var portfolioResp PortfolioDetailResponse
	if err := c.handleResponse(resp, &portfolioResp, "portfolio detail"); err != nil {
		return &portfolioResp, err
	}

	return &portfolioResp, nil
}

package tiqs

import (
	"encoding/json"
	"fmt"
)

// HolidaysResponse represents the API response structure for market holidays.
type HolidaysResponse struct {
	Data struct {
		Holidays           map[string]string          `json:"holidays"`           // A map of holiday dates and their descriptions.
		SpecialTradingDays map[string][][]interface{} `json:"specialTradingDays"` // A map of special trading days with additional details.
	} `json:"data"`
	Status string `json:"status"` // API response status (e.g., "success" or "error").
}

// IndexListResponse represents the API response structure for retrieving a list of market indices.
type IndexListResponse struct {
	Data []struct {
		Name  string `json:"name"`  // Name of the index (e.g., NIFTY 50, SENSEX).
		Token string `json:"token"` // Unique identifier (token) for the index.
	} `json:"data"`
	Status string `json:"status"` // API response status (e.g., "success" or "error").
}

// OptionChainSymbolResponse represents the API response structure for retrieving option chain symbols.
type OptionChainSymbolResponse struct {
	Data   map[string][]string `json:"data"`   // A map containing option chain symbols grouped by category.
	Status string              `json:"status"` // API response status (e.g., "success" or "error").
}

// GetHolidays fetches the list of market holidays and special trading days.
//
// It sends a GET request to the "/info/holidays" endpoint to retrieve market holiday
// schedules and special trading days.
//
// Returns:
//   - A pointer to a HolidaysResponse struct containing holiday details if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetHolidays() (*HolidaysResponse, error) {
	endpoint := "/info/holidays"

	// Send a GET request to fetch market holidays.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into the HolidaysResponse struct.
	var holidaysResponse HolidaysResponse
	if err := json.Unmarshal(resp, &holidaysResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal holidays response: %w", err)
	}

	return &holidaysResponse, nil
}

// GetIndexList fetches the list of available stock market indices.
//
// It sends a GET request to the "/info/index-list" endpoint to retrieve details of
// available indices, including their names and unique tokens.
//
// Returns:
//   - A pointer to an IndexListResponse struct containing index details if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetIndexList() (*IndexListResponse, error) {
	endpoint := "/info/index-list"

	// Send a GET request to fetch the list of indices.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into the IndexListResponse struct.
	var indexListResponse IndexListResponse
	if err := json.Unmarshal(resp, &indexListResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index list response: %w", err)
	}

	return &indexListResponse, nil
}

// GetOptionChainSymbol fetches the available option chain symbols.
//
// It sends a GET request to the "/info/option-chain-symbols" endpoint to retrieve
// the available option chain symbols categorized by different asset types.
//
// Returns:
//   - A pointer to an OptionChainSymbolResponse struct containing option chain symbols if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetOptionChainSymbol() (*OptionChainSymbolResponse, error) {
	endpoint := "/info/option-chain-symbols"

	// Send a GET request to fetch option chain symbols.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into the OptionChainSymbolResponse struct.
	var optionChainSymbolResponse OptionChainSymbolResponse
	if err := json.Unmarshal(resp, &optionChainSymbolResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal option chain symbols response: %w", err)
	}

	return &optionChainSymbolResponse, nil
}

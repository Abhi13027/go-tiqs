// market.go
package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// MarketQuote represents the response structure for market quotes.
type MarketQuote struct {
	Token        int64  `json:"token"`        // Unique identifier for the instrument.
	LTP          int64  `json:"ltp"`          // Last traded price of the instrument.
	Open         int64  `json:"open"`         // Opening price of the instrument for the trading session.
	High         int64  `json:"high"`         // Highest price of the instrument in the current session.
	Low          int64  `json:"low"`          // Lowest price of the instrument in the current session.
	Close        int64  `json:"close"`        // Closing price of the instrument from the previous session.
	Volume       int64  `json:"volume"`       // Total traded volume of the instrument.
	TotalBuyQty  int64  `json:"totalBuyQty"`  // Total quantity of buy orders in the market.
	TotalSellQty int64  `json:"totalSellQty"` // Total quantity of sell orders in the market.
	LTT          int64  `json:"ltt"`          // Last trade time of the instrument (epoch timestamp).
	Status       string `json:"status"`       // API response status (e.g., "success" or "error").
}

// GetMarketQuote fetches market data for a single instrument.
//
// It sends a POST request to the "/info/quote/{mode}" endpoint to retrieve
// the latest market details for a given token.
//
// Parameters:
//   - token: The unique identifier of the instrument.
//   - mode: Market mode (e.g., "full", "ltp", "depth").
//
// Returns:
//   - A pointer to MarketQuote struct containing market data if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetMarketQuote(token int64, mode string) (*MarketQuote, error) {
	endpoint := fmt.Sprintf("/info/quote/%s", mode)
	payload := fmt.Sprintf(`{"token": %d}`, token)

	// Send a POST request to fetch market data.
	resp, err := c.request(endpoint, "POST", []byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch market quote")
		return nil, err
	}

	var result struct {
		Status string      `json:"status"`
		Data   MarketQuote `json:"data"`
	}

	// Parse the JSON response into the MarketQuote struct.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse market quote response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("market data retrieval failed")
	}

	log.Info().Int64("token", token).Msg("Market quote retrieved successfully")
	return &result.Data, nil
}

// GetMarketQuotes fetches market data for multiple instruments.
//
// It sends a POST request to the "/info/quotes/{mode}" endpoint to retrieve
// market data for a list of tokens.
//
// Parameters:
//   - tokens: A slice of unique identifiers representing instruments.
//   - mode: Market mode (e.g., "full", "ltp", "depth").
//
// Returns:
//   - A slice of MarketQuote structs containing market data if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetMarketQuotes(tokens []int64, mode string) ([]MarketQuote, error) {
	endpoint := fmt.Sprintf("/info/quotes/%s", mode)

	// Construct JSON payload for multiple tokens.
	payload := "["
	for i, token := range tokens {
		payload += fmt.Sprintf("%d", token)
		if i < len(tokens)-1 {
			payload += ","
		}
	}
	payload += "]"

	// Send a POST request to fetch market data for multiple tokens.
	resp, err := c.request(endpoint, "POST", []byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch market quotes")
		return nil, err
	}

	var result struct {
		Status string        `json:"status"`
		Data   []MarketQuote `json:"data"`
	}

	// Parse the JSON response into a slice of MarketQuote structs.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse market quotes response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("market data retrieval failed")
	}

	log.Info().Msg("Market quotes retrieved successfully")
	return result.Data, nil
}

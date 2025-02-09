// historical.go
package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// HistoricalCandle represents a single OHLCV (Open, High, Low, Close, Volume) data point.
type HistoricalCandle struct {
	Time   string `json:"time"`         // Timestamp of the candle in ISO 8601 format.
	Open   int64  `json:"open"`         // Open price of the candle.
	High   int64  `json:"high"`         // Highest price during the candle period.
	Low    int64  `json:"low"`          // Lowest price during the candle period.
	Close  int64  `json:"close"`        // Closing price of the candle.
	Volume int64  `json:"volume"`       // Trading volume during the candle period.
	OI     *int64 `json:"oi,omitempty"` // Open Interest (optional, included if requested).
}

// HistoricalDataResponse represents the structure of the historical data API response.
type HistoricalDataResponse struct {
	Status string             `json:"status"` // API response status (e.g., "success" or "error").
	Data   []HistoricalCandle `json:"data"`   // List of historical OHLCV candles.
}

// GetHistoricalData fetches historical OHLCV data for a given instrument.
//
// It sends a GET request to the "/candle/{exchange}/{token}/{interval}?from={from}&to={to}" endpoint
// to retrieve OHLCV data for the specified time range. If Open Interest (OI) is requested, it is appended
// as a query parameter.
//
// Parameters:
//   - exchange: The exchange where the instrument is listed (e.g., NSE, BSE).
//   - token: The unique identifier of the instrument.
//   - interval: The timeframe of the candles (e.g., "1m", "5m", "1d").
//   - from: The start date/time for historical data (ISO 8601 format).
//   - to: The end date/time for historical data (ISO 8601 format).
//   - includeOI: Boolean flag to include Open Interest (OI) data if available.
//
// Returns:
//   - A slice of HistoricalCandle structs containing OHLCV data if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetHistoricalData(exchange, token, interval, from, to string, includeOI bool) ([]HistoricalCandle, error) {
	endpoint := fmt.Sprintf("/candle/%s/%s/%s?from=%s&to=%s", exchange, token, interval, from, to)

	// If Open Interest (OI) is requested, append it as a query parameter.
	if includeOI {
		endpoint += "&oi=1"
	}

	// Send a GET request to the API to fetch historical data.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch historical data")
		return nil, err
	}

	var result HistoricalDataResponse
	// Parse the JSON response into the HistoricalDataResponse struct.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse historical data response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("historical data retrieval failed")
	}

	log.Info().
		Str("exchange", exchange).
		Str("token", token).
		Str("interval", interval).
		Bool("includeOI", includeOI).
		Msg("Historical data retrieved successfully")

	return result.Data, nil
}

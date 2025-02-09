package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Position represents a trading position held by the user.
type Position struct {
	AvgPrice                 string `json:"avgPrice"`                 // Average price of the position.
	BreakEvenPrice           string `json:"breakEvenPrice"`           // Break-even price for the position.
	CarrtForwardAvgPrice     string `json:"carrtForwardAvgPrice"`     // Average price carried forward from the previous session.
	CarryForwardBuyAmount    string `json:"carryForwardBuyAmount"`    // Amount spent on buy orders carried forward.
	CarryForwardBuyAvgPrice  string `json:"carryForwardBuyAvgPrice"`  // Average buy price of carry-forward positions.
	CarryForwardBuyQty       string `json:"carryForwardBuyQty"`       // Quantity of buy orders carried forward.
	CarryForwardSellAmount   string `json:"carryForwardSellAmount"`   // Amount received from sell orders carried forward.
	CarryForwardSellAvgPrice string `json:"carryForwardSellAvgPrice"` // Average sell price of carry-forward positions.
	CarryForwardSellQty      string `json:"carryForwardSellQty"`      // Quantity of sell orders carried forward.
	DayBuyAmount             string `json:"dayBuyAmount"`             // Total amount spent on buy orders during the day.
	DayBuyAvgPrice           string `json:"dayBuyAvgPrice"`           // Average price of buy orders executed during the day.
	DayBuyQty                string `json:"dayBuyQty"`                // Total quantity of buy orders executed during the day.
	DaySellAmount            string `json:"daySellAmount"`            // Total amount received from sell orders during the day.
	DaySellAvgPrice          string `json:"daySellAvgPrice"`          // Average price of sell orders executed during the day.
	DaySellQty               string `json:"daySellQty"`               // Total quantity of sell orders executed during the day.
	Exchange                 string `json:"exchange"`                 // The exchange where the position is held (e.g., NSE, BSE).
	LotSize                  string `json:"lotSize"`                  // Lot size of the instrument.
	Ltp                      string `json:"ltp"`                      // Last traded price of the instrument.
	Multiplier               string `json:"multiplier"`               // Multiplier factor for derivative positions.
	NetUploadPrice           string `json:"netUploadPrice"`           // Net upload price for the position.
	OpenBuyAmount            string `json:"openBuyAmount"`            // Total amount spent on open buy orders.
	OpenBuyAvgPrice          string `json:"openBuyAvgPrice"`          // Average price of open buy orders.
	OpenBuyQty               string `json:"openBuyQty"`               // Quantity of open buy orders.
	OpenSellAmount           string `json:"openSellAmount"`           // Total amount received from open sell orders.
	OpenSellAvgPrice         string `json:"openSellAvgPrice"`         // Average price of open sell orders.
	OpenSellQty              string `json:"openSellQty"`              // Quantity of open sell orders.
	PriceFactor              string `json:"priceFactor"`              // Price factor applied to the instrument.
	PricePrecision           string `json:"pricePrecision"`           // Precision of price representation.
	Product                  string `json:"product"`                  // Trading product type (e.g., MIS, CNC, NRML).
	Qty                      string `json:"qty"`                      // Total quantity held in the position.
	RealisedPnL              string `json:"realisedPnL"`              // Realized profit and loss from the position.
	Symbol                   string `json:"symbol"`                   // Trading symbol of the instrument.
	TickSize                 string `json:"tickSize"`                 // Tick size of the instrument.
	Token                    string `json:"token"`                    // Unique token identifier for the instrument.
	UnrealisedMarkToMarket   string `json:"unrealisedMarkToMarket"`   // Unrealized Mark-to-Market (MTM) value of the position.
	UploadPrice              string `json:"uploadPrice"`              // Uploaded price for position tracking.
}

// PositionsResponse represents the API response for user positions.
type PositionsResponse struct {
	Data   []Position `json:"data"`   // List of positions held by the user.
	Status string     `json:"status"` // API response status (e.g., "success" or "error").
}

// GetPositions fetches the positions for the authenticated user.
//
// It sends a GET request to the "/user/positions" endpoint to retrieve all open
// and carry-forward positions.
//
// Returns:
//   - A slice of Position structs containing all active positions if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetPositions() ([]Position, error) {
	endpoint := "/user/positions"

	// Send a GET request to the API to fetch position details.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch positions")
		return nil, err
	}

	var result PositionsResponse
	// Parse the JSON response into the PositionsResponse struct.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse positions response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("positions retrieval failed")
	}

	log.Info().Msg("Positions retrieved successfully")
	return result.Data, nil
}

package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Position represents a trading position held by the user.

type Position struct {
	AvgPrice                 string `json:"avgPrice"`                 // Average price of the position.
	BreakEvenPrice           string `json:"breakEvenPrice"`           // Break-even price of the position.
	CarryForwarAvgPrice      string `json:"carryForwarAvgPrice"`      // Carry-forward average price of the position.
	CarryForwardBuyAmount    string `json:"carryForwardBuyAmount"`    // Carry-forward buy amount.
	CarryForwardBuyAvgPrice  string `json:"carryForwardBuyAvgPrice"`  // Carry-forward buy average price.
	CarryForwardBuyQty       string `json:"carryForwardBuyQty"`       //	Carry-forward buy quantity.
	CarryForwardSellAmount   string `json:"carryForwardSellAmount"`   // Carry-forward sell amount.
	CarryForwardSellAvgPrice string `json:"carryForwardSellAvgPrice"` // Carry-forward sell average price.
	CarryForwardSellQty      string `json:"carryForwardSellQty"`      // Carry-forward sell quantity.
	DayBuyAmount             string `json:"dayBuyAmount"`             // Day buy amount.
	DayBuyAvgPrice           string `json:"dayBuyAvgPrice"`           // Day buy average price.
	DayBuyQty                string `json:"dayBuyQty"`                //	Day buy quantity.
	DaySellAmount            string `json:"daySellAmount"`            // Day sell amount.
	DaySellAvgPrice          string `json:"daySellAvgPrice"`          // Day sell average price.
	DaySellQty               string `json:"daySellQty"`               // Day sell quantity.
	Exchange                 string `json:"exchange"`                 // Exchange where the position is held.
	LotSize                  string `json:"lotSize"`                  // Lot size of the position.
	Ltp                      string `json:"ltp"`                      // Last traded price of the position.
	Multiplier               string `json:"multiplier"`               // Multiplier of the position.
	NetBuyQty                string `json:"netBuyQty"`                // Net buy quantity.
	NetSellQty               string `json:"netSellQty"`               // Net sell quantity.
	NetUploadPrice           string `json:"netUploadPrice"`           // Net upload price.
	OpenBuyAmount            string `json:"openBuyAmount"`            // Open buy amount.
	OpenBuyAvgPrice          string `json:"openBuyAvgPrice"`          // Open buy average price.
	OpenBuyQty               string `json:"openBuyQty"`               // Open buy quantity.
	OpenSellAmount           string `json:"openSellAmount"`           // Open sell amount.
	OpenSellAvgPrice         string `json:"openSellAvgPrice"`         // Open sell average price.
	OpenSellQty              string `json:"openSellQty"`              // Open sell quantity.
	Pnl                      string `json:"pnl"`                      // Profit and loss of the position.
	PriceFactor              string `json:"priceFactor"`              // Price factor of the position.
	PricePrecision           string `json:"pricePrecision"`           // Price precision of the position.
	Product                  string `json:"product"`                  // Product type of the position.
	Qty                      string `json:"qty"`                      // Quantity of the position.
	RealisedPnL              string `json:"realisedPnL"`              // Realised profit and loss of the position.
	Symbol                   string `json:"symbol"`                   // Trading symbol of the position.
	TickSize                 string `json:"tickSize"`                 // Tick size of the position.
	Token                    string `json:"token"`                    // Token of the position.
	UnRealisedPnl            string `json:"unRealisedPnl"`            // Unrealised profit and loss of the position.
	UnrealisedMarkToMarket   string `json:"unrealisedMarkToMarket"`   // Unrealised mark-to-market of the position.
	UploadPrice              string `json:"uploadPrice"`              // Upload price of the position.
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

package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Holding represents a user's stock or asset holding details.
type Holding struct {
	AuthorizedQty       string  `json:"authorizedQty"`       // Authorized quantity of the holding.
	AvgPrice            string  `json:"avgPrice"`            // Average price at which the holding was acquired.
	BrokerCollateralQty string  `json:"brokerCollateralQty"` // Quantity pledged as collateral with the broker.
	Close               float64 `json:"close"`               // Closing price of the holding from the previous session.
	CollateralQty       string  `json:"collateralQty"`       // Total collateral quantity.
	DepositoryQty       string  `json:"depositoryQty"`       // Quantity held in the depository.
	EffectiveQty        string  `json:"effectiveQty"`        // Effective quantity available for trading.
	Exchange            string  `json:"exchange"`            // Exchange where the holding is listed (e.g., NSE, BSE).
	Haircut             string  `json:"haircut"`             // Haircut percentage applied to the collateral.
	Ltp                 float64 `json:"ltp"`                 // Last traded price of the holding.
	NonPOAQty           string  `json:"nonPOAQty"`           // Quantity not under Power of Attorney (POA).
	Pnl                 string  `json:"pnl"`                 // Profit and Loss (PnL) on the holding.
	Qty                 string  `json:"qty"`                 // Total quantity held.
	SellableQty         string  `json:"sellableQty"`         // Quantity available for selling.
	Symbol              string  `json:"symbol"`              // Trading symbol of the instrument.
	T1Qty               string  `json:"t1Qty"`               // T+1 quantity, which is yet to be settled.
	Token               string  `json:"token"`               // Unique token identifier for the holding.
	TradingSymbol       string  `json:"tradingSymbol"`       // Full trading symbol of the instrument.
	UnPledgedQty        string  `json:"unPledgedQty"`        // Quantity that is not pledged as collateral.
	UsedQty             string  `json:"usedQty"`             // Quantity already used (e.g., for margin or pledging).
}

// HoldingsResponse represents the API response containing user holdings.
type HoldingsResponse struct {
	Data   []Holding `json:"data"`   // List of user's holdings.
	Status string    `json:"status"` // API response status (e.g., "success" or "error").
}

// GetHoldings fetches the holdings for the authenticated user.
//
// It sends a GET request to the "/user/holdings" endpoint to retrieve all holdings
// associated with the user's account.
//
// Returns:
//   - A slice of Holding structs containing all available holdings if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetHoldings() ([]Holding, error) {
	endpoint := "/user/holdings"

	// Send a GET request to the API to fetch holdings.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch holdings")
		return nil, err
	}

	var result HoldingsResponse
	// Parse the JSON response into the HoldingsResponse struct.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse holdings response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("holdings retrieval failed")
	}

	log.Info().Msg("Holdings retrieved successfully")
	return result.Data, nil
}

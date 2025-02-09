package tiqs

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// MarginRequest represents the structure for a single order margin request.
type MarginRequest struct {
	Exchange        string `json:"exchange"`        // Exchange where the order is placed (e.g., NSE, BSE).
	Token           string `json:"token"`           // Unique identifier for the instrument.
	Quantity        string `json:"quantity"`        // Order quantity.
	Product         string `json:"product"`         // Product type (e.g., MIS, CNC, NRML).
	Price           string `json:"price"`           // Order price (applicable for LIMIT orders).
	TransactionType string `json:"transactionType"` // Order transaction type (BUY/SELL).
	OrderType       string `json:"order"`           // Type of order (e.g., MARKET, LIMIT).
	Symbol          string `json:"symbol"`          // Trading symbol of the instrument.
}

// BasketMarginRequest represents a collection of margin requests for multiple orders.
type BasketMarginRequest []MarginRequest

// MarginResponse represents the response structure for margin calculations.
type MarginResponse struct {
	Cash   string `json:"cash"` // Available cash balance.
	Charge struct {
		Brokerage      int     `json:"brokerage"`      // Brokerage fees.
		SebiCharges    float64 `json:"sebiCharges"`    // SEBI charges.
		ExchangeTxnFee float64 `json:"exchangeTxnFee"` // Exchange transaction fees.
		StampDuty      float64 `json:"stampDuty"`      // Stamp duty applicable.
		Ipft           float64 `json:"ipft"`           // Investor Protection Fund Trust (IPFT) fees.
		TransactionTax int     `json:"transactionTax"` // Transaction tax applied.

		Gst struct {
			Cgst  int     `json:"cgst"`  // Central GST amount.
			Sgst  int     `json:"sgst"`  // State GST amount.
			Igst  float64 `json:"igst"`  // Integrated GST amount.
			Total float64 `json:"total"` // Total GST amount.
		} `json:"gst"`

		Total float64 `json:"total"` // Total charge applied.
	} `json:"charge"`

	Margin     string `json:"margin"`     // Required margin for the order.
	MarginUsed string `json:"marginUsed"` // Margin already used.
}

// OrderMargin represents the API response for a single order margin request.
type OrderMargin struct {
	Data   MarginResponse `json:"data"`   // Margin calculation details.
	Status string         `json:"status"` // API response status (e.g., "success" or "error").
}

// BasketOrderMargin represents the API response for multiple order margin requests.
type BasketOrderMargin struct {
	Data struct {
		MarginUsed           string `json:"marginUsed"`           // Total margin used before placing orders.
		MarginUsedAfterTrade string `json:"marginUsedAfterTrade"` // Total margin used after trade execution.
	} `json:"data"`
	Status string `json:"status"` // API response status (e.g., "success" or "error").
}

// GetMargin fetches the margin details for a single order.
//
// It sends a POST request to the "/margin/order" endpoint with the order details
// to calculate the required margin for the specified transaction.
//
// Parameters:
//   - order: A MarginRequest struct containing the order details.
//
// Returns:
//   - A pointer to an OrderMargin struct with margin details if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetMargin(order MarginRequest) (*OrderMargin, error) {
	endpoint := "/margin/order"

	// Convert order details into JSON payload.
	payload, err := json.Marshal(order)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize margin request")
		return nil, err
	}

	// Send the request to the API.
	resp, err := c.request(endpoint, "POST", []byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch margin")
		return nil, err
	}

	// Parse the JSON response into the OrderMargin struct.
	var result OrderMargin
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse margin response")
		return nil, err
	}

	return &result, nil
}

// GetBasketMargin fetches the margin details for multiple orders.
//
// This function sends a POST request to the "/margin/basket" endpoint with a collection
// of orders to calculate the combined margin requirements.
//
// Parameters:
//   - order: A BasketMarginRequest struct containing multiple orders.
//
// Returns:
//   - A pointer to a BasketOrderMargin struct with total margin details if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetBasketMargin(order BasketMarginRequest) (*BasketOrderMargin, error) {
	endpoint := "/margin/basket"

	// Convert order details into JSON payload.
	payload, err := json.Marshal(order)
	log.Info().Msgf("Payload: %s", payload) // Log the payload for debugging.
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize margin request")
		return nil, err
	}

	// Send the request to the API.
	resp, err := c.request(endpoint, "POST", []byte(payload))
	log.Info().Msgf("Response: %s", resp) // Log the response for debugging.
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch margin")
		return nil, err
	}

	// Parse the JSON response into the BasketOrderMargin struct.
	var result BasketOrderMargin
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse margin response")
		return nil, err
	}

	return &result, nil
}

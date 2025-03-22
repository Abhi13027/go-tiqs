// orders.go
package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// OrderRequest represents the structure for placing an order.
type OrderRequest struct {
	Exchange        string `json:"exchange"`                // Exchange where the order is placed (e.g., NSE, BSE).
	Token           string `json:"token"`                   // Unique identifier for the instrument.
	Quantity        string `json:"quantity"`                // Order quantity.
	DisclosedQty    string `json:"disclosedQty,omitempty"`  // Disclosed quantity (optional).
	Product         string `json:"product"`                 // Product type (e.g., MIS, CNC, NRML).
	Symbol          string `json:"symbol"`                  // Trading symbol of the instrument.
	TransactionType string `json:"transactionType"`         // Order transaction type (BUY/SELL).
	OrderType       string `json:"order"`                   // Type of order (e.g., MARKET, LIMIT).
	Price           string `json:"price"`                   // Order price (applicable for LIMIT orders).
	Validity        string `json:"validity"`                // Order validity (e.g., DAY, IOC).
	Tags            string `json:"tags,omitempty"`          // Custom tags for order tracking (optional).
	AMO             bool   `json:"amo,omitempty"`           // Indicates if the order is an After Market Order (AMO).
	TriggerPrice    string `json:"triggerPrice,omitempty"`  // Trigger price for stop-loss or conditional orders.
	BookLossPrice   string `json:"bookLossPrice,omitempty"` // Book loss price for risk management.
}

// OrderResponse represents the API response after placing an order.
type OrderResponse struct {
	Status    string `json:"status"`              // API response status (e.g., "success", "error").
	Message   string `json:"message,omitempty"`   // Message from the API (if any).
	ErrorCode string `json:"errorCode,omitempty"` // Error code in case of failure.
	Data      struct {
		OrderNo     string `json:"orderNo,omitempty"`     // Order number assigned by the exchange.
		RequestTime string `json:"requestTime,omitempty"` // Timestamp of the order request.
	} `json:"data,omitempty"`
}

type OrderDetailsResponse struct {
	Data []struct {
		Status             string `json:"status"`
		Exchange           string `json:"exchange"`
		Symbol             string `json:"symbol"`
		ID                 string `json:"id"`
		Price              string `json:"price"`
		Quantity           string `json:"quantity"`
		Product            string `json:"product"`
		OrderStatus        string `json:"orderStatus"`
		ReportType         string `json:"reportType"`
		TransactionType    string `json:"transactionType"`
		Order              string `json:"order"`
		FillShares         string `json:"fillShares"`
		AveragePrice       string `json:"averagePrice"`
		RejectReason       string `json:"rejectReason"`
		ExchangeOrderID    string `json:"exchangeOrderID"`
		CancelQuantity     string `json:"cancelQuantity"`
		Remarks            string `json:"remarks"`
		DisclosedQuantity  string `json:"disclosedQuantity"`
		OrderTriggerPrice  string `json:"orderTriggerPrice"`
		Retention          string `json:"retention"`
		BookProfitPrice    string `json:"bookProfitPrice"`
		BookLossPrice      string `json:"bookLossPrice"`
		TrailingPrice      string `json:"trailingPrice"`
		Amo                string `json:"amo"`
		PricePrecision     string `json:"pricePrecision"`
		TickSize           string `json:"tickSize"`
		LotSize            string `json:"lotSize"`
		Token              string `json:"token"`
		TimeStamp          string `json:"timeStamp"`
		OrderTime          string `json:"orderTime"`
		ExchangeUpdateTime string `json:"exchangeUpdateTime"`
		RequestTime        string `json:"requestTime"`
		ErrorMessage       string `json:"errorMessage"`
	} `json:"data"`
	Status string `json:"status"`
}

// PlaceOrder places a new order in the market.
//
// It sends a POST request to the API endpoint "/order/{orderType}" with the order details.
//
// Parameters:
//   - orderType: Type of order (e.g., MARKET, LIMIT).
//   - order: OrderRequest struct containing the order details.
//
// Returns:
//   - A pointer to OrderResponse with the order confirmation details if successful.
//   - An error if the order placement fails.
func (c *Client) PlaceOrder(orderType string, order OrderRequest) (*OrderResponse, error) {
	endpoint := fmt.Sprintf("/order/%s", orderType)

	payload, err := json.Marshal(order)
	log.Info().Str("payload", string(payload)).Msg("Placing order")
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize order request")
		return nil, err
	}

	resp, err := c.request(endpoint, "POST", payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to place order")
		return nil, err
	}

	var result OrderResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse order response")
		return nil, err
	}

	if result.Status != "success" {
		log.Error().Str("errorCode", result.ErrorCode).Str("message", result.Message).Msg("Order placement failed")
		return nil, fmt.Errorf("order placement failed")
	}

	log.Info().Str("orderNo", result.Data.OrderNo).Msg("Order placed successfully")
	return &result, nil
}

// ModifyOrder modifies an existing order.
//
// It sends a PATCH request to the API endpoint "/order/{orderType}/{orderID}" with the modified order details.
//
// Parameters:
//   - orderType: Type of the order being modified (e.g., MARKET, LIMIT).
//   - orderID: Unique identifier of the order to be modified.
//   - order: OrderRequest struct containing updated order details.
//
// Returns:
//   - A pointer to OrderResponse with the updated order details if successful.
//   - An error if the modification fails.
func (c *Client) ModifyOrder(orderType, orderID string, order OrderRequest) (*OrderResponse, error) {
	endpoint := fmt.Sprintf("/order/%s/%s", orderType, orderID)

	payload, err := json.Marshal(order)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize modify order request")
		return nil, err
	}

	resp, err := c.request(endpoint, "PATCH", payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to modify order")
		return nil, err
	}

	var result OrderResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse modify order response")
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("order modification failed")
	}

	log.Info().Str("orderNo", result.Data.OrderNo).Msg("Order modified successfully")
	return &result, nil
}

// CancelOrder cancels an existing order.
//
// It sends a DELETE request to the API endpoint "/order/{orderType}/{orderID}".
//
// Parameters:
//   - orderType: Type of the order to be canceled (e.g., MARKET, LIMIT).
//   - orderID: Unique identifier of the order.
//
// Returns:
//   - An error if the cancellation fails; otherwise, nil.
func (c *Client) CancelOrder(orderType, orderID string) error {
	endpoint := fmt.Sprintf("/order/%s/%s", orderType, orderID)

	resp, err := c.request(endpoint, "DELETE", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cancel order")
		return err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Message string `json:"message"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse cancel order response")
		return err
	}

	if result.Status != "success" {
		return fmt.Errorf("order cancellation failed")
	}

	log.Info().Str("message", result.Data.Message).Msg("Order cancelled successfully")
	return nil
}

// GetOrder retrieves details of a specific order.
//
// It sends a GET request to the API endpoint "/order/{orderID}".
//
// Parameters:
//   - orderID: Unique identifier of the order.
//
// Returns:
//   - A pointer to OrderResponse containing order details if successful.
//   - An error if the retrieval fails.
func (c *Client) GetOrder(orderID string) (*OrderDetailsResponse, error) {
	endpoint := fmt.Sprintf("/order/%s", orderID)

	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get order details")
		return nil, err
	}

	var result OrderDetailsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse order details response")
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("failed to retrieve order details")
	}

	log.Info().Str("orderNo", orderID).Msg("Order details retrieved successfully")
	return &result, nil
}

// GetOrderBook retrieves all orders for the current trading day.
//
// It sends a GET request to the API endpoint "/user/orders" and returns a list of orders.
//
// Returns:
//   - A slice of OrderResponse structs containing all orders if successful.
//   - An error if the retrieval fails.
func (c *Client) GetOrderBook() ([]OrderResponse, error) {
	endpoint := "/user/orders"

	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch order book")
		return nil, err
	}

	var result struct {
		Status string          `json:"status"`
		Data   []OrderResponse `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse order book response")
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("failed to retrieve order book")
	}

	log.Info().Msg("Order book retrieved successfully")
	return result.Data, nil
}

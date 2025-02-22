package ticks

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	WSS_URL = "wss://wss.tiqs.trading"
)

// DepthLevel represents a single level in the market depth
type DepthLevel struct {
	Quantity int64 `json:"quantity"`
	Price    int32 `json:"price"`
	Orders   int16 `json:"orders"`
}

// MarketDepth represents the full market depth with bids and asks
type MarketDepth struct {
	Bids [5]DepthLevel `json:"bids"`
	Asks [5]DepthLevel `json:"asks"`
}

// TickData represents the complete market data for a token
type TickData struct {
	Token              int32       `json:"token"`
	LTP                int32       `json:"ltp"`
	NetChangeIndicator int32       `json:"net_change_indicator"`
	NetChange          int32       `json:"net_change"`
	LTQ                int32       `json:"ltq"`
	AvgPrice           int32       `json:"avg_price"`
	TotalBuyQty        int64       `json:"total_buy_qty"`
	TotalSellQty       int64       `json:"total_sell_qty"`
	Open               int32       `json:"open"`
	High               int32       `json:"high"`
	Close              int32       `json:"close"`
	Low                int32       `json:"low"`
	Volume             int64       `json:"volume"`
	LTT                int32       `json:"ltt"`
	Time               int32       `json:"time"`
	OI                 int32       `json:"oi"`
	OIDayHigh          int32       `json:"oi_day_high"`
	OIDayLow           int32       `json:"oi_day_low"`
	LowerLimit         int32       `json:"lower_limit"`
	UpperLimit         int32       `json:"upper_limit"`
	MarketDepth        MarketDepth `json:"market_depth"`
}

// WS represents the WebSocket client
type WS struct {
	AppID         string
	Token         string
	TokenList     []int
	Conn          *websocket.Conn
	URL           string
	RetryDelay    time.Duration
	MaxRetries    int
	ctx           context.Context
	cancel        context.CancelFunc
	logger        *zerolog.Logger
	DataChan      chan TickData
	errChan       chan error
	subscriptions sync.Map
	mu            sync.RWMutex
}

// NewWS creates a new WebSocket client instance
func NewWS(appId, token string) *WS {
	ctx, cancel := context.WithCancel(context.Background())
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &WS{
		AppID:      appId,
		Token:      token,
		TokenList:  make([]int, 0),
		URL:        WSS_URL,
		RetryDelay: 5 * time.Second,
		MaxRetries: 25,
		ctx:        ctx,
		cancel:     cancel,
		logger:     &logger,
		DataChan:   make(chan TickData, 1000),
		errChan:    make(chan error, 100),
	}
}

// Connect establishes a WebSocket connection
func (ws *WS) Connect() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	var err error
	for attempt := 1; attempt <= ws.MaxRetries; attempt++ {
		ws.logger.Info().Msgf("Attempting to connect to WebSocket (attempt %d/%d)", attempt, ws.MaxRetries)

		url := fmt.Sprintf("%s?appId=%s&token=%s", ws.URL, ws.AppID, ws.Token)
		ws.Conn, _, err = websocket.DefaultDialer.Dial(url, nil)

		if err == nil {
			ws.logger.Info().Msg("Connected to WebSocket")

			// Resubscribe to existing subscriptions
			ws.resubscribeAll()

			// Start message handler
			go ws.handleMessages()
			return nil
		}

		ws.logger.Error().Err(err).Msgf("Failed to connect. Retrying in %s...", ws.RetryDelay)
		time.Sleep(ws.RetryDelay)
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", ws.MaxRetries, err)
}

// Subscribe subscribes to market data for given tokens
func (ws *WS) Subscribe(tokens []int, mode string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	message := map[string]interface{}{
		"code": "sub",
		"mode": mode,
		mode:   tokens,
	}

	// Store subscription
	for _, token := range tokens {
		ws.subscriptions.Store(token, mode)
	}

	ws.TokenList = append(ws.TokenList, tokens...)
	return ws.sendJSONMessage(message)
}

// Unsubscribe removes subscription for given tokens
func (ws *WS) Unsubscribe(tokens []int, mode string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	message := map[string]interface{}{
		"code": "unsub",
		"mode": mode,
		mode:   tokens,
	}

	// Remove subscription
	for _, token := range tokens {
		ws.subscriptions.Delete(token)
	}

	return ws.sendJSONMessage(message)
}

// GetDataChannel returns the channel for receiving market data
func (ws *WS) GetDataChannel() <-chan TickData {
	return ws.DataChan
}

// GetErrorChannel returns the channel for receiving errors
func (ws *WS) GetErrorChannel() <-chan error {
	return ws.errChan
}

// Close closes the WebSocket connection and cleanup
func (ws *WS) Close() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.cancel() // Stop all goroutines

	// Close channels
	close(ws.DataChan)
	close(ws.errChan)

	if ws.Conn != nil {
		ws.logger.Info().Msg("Closing WebSocket connection")
		return ws.Conn.Close()
	}
	return nil
}

// handleMessages processes incoming WebSocket messages
func (ws *WS) handleMessages() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			if ws.Conn == nil {
				return
			}

			messageType, message, err := ws.Conn.ReadMessage()
			if err != nil {
				ws.logger.Error().Err(err).Msg("Error reading message")
				ws.errChan <- err
				ws.reconnect()
				return
			}

			if messageType == websocket.BinaryMessage {
				tickData, err := ws.parseBinaryToTickData(message)
				if err != nil {
					ws.logger.Error().Err(err).Msg("Error parsing binary data")
					continue
				}

				// Send data to channel (non-blocking)
				select {
				case ws.DataChan <- tickData:
					// Data sent successfully
				default:
					ws.logger.Warn().Msg("Data channel is full, skipping message")
				}
			}
		}
	}
}

// parseBinaryToTickData converts binary message to TickData struct
func (ws *WS) parseBinaryToTickData(data []byte) (TickData, error) {
	var tick TickData

	if len(data) == 1 {
		tick.Token = int32(data[0])
		return tick, nil
	}

	if len(data) < 17 {
		return tick, fmt.Errorf("invalid data length: %d", len(data))
	}

	// Parse basic fields
	tick.Token = bigEndianToInt(data[:4])
	tick.LTP = bigEndianToInt(data[4:8])

	if len(data) == 17 {
		tick.Close = bigEndianToInt(data[13:17])
		tick.NetChange = int32((float64(tick.LTP-tick.Close) / float64(tick.Close)) * 100)

		if tick.LTP > tick.Close {
			tick.NetChangeIndicator = 43 // '+'
		} else if tick.LTP < tick.Close {
			tick.NetChangeIndicator = 45 // '-'
		} else {
			tick.NetChangeIndicator = 32 // ' '
		}
	}

	if len(data) >= 81 {
		tick.AvgPrice = bigEndianToInt(data[17:21])
		tick.TotalBuyQty = int64(bigEndianToInt(data[21:29]))
		tick.TotalSellQty = int64(bigEndianToInt(data[29:37]))
		tick.Open = bigEndianToInt(data[37:41])
		tick.High = bigEndianToInt(data[41:45])
		tick.Close = bigEndianToInt(data[45:49])
		tick.Low = bigEndianToInt(data[49:53])
		tick.Volume = int64(bigEndianToInt(data[53:61]))
		tick.LTT = bigEndianToInt(data[61:65])
		tick.Time = bigEndianToInt(data[65:69])
		tick.OI = bigEndianToInt(data[69:73])
		tick.OIDayHigh = bigEndianToInt(data[73:77])
		tick.OIDayLow = bigEndianToInt(data[77:81])
	}

	if len(data) == 229 {
		tick.LowerLimit = bigEndianToInt(data[81:85])
		tick.UpperLimit = bigEndianToInt(data[85:89])

		// Parse market depth
		offset := 89
		for i := 0; i < 5; i++ {
			// Parse bids
			tick.MarketDepth.Bids[i] = DepthLevel{
				Quantity: int64(bigEndianToInt(data[offset : offset+8])),
				Price:    bigEndianToInt(data[offset+8 : offset+12]),
				Orders:   int16(bigEndianToInt(data[offset+12 : offset+14])),
			}
			offset += 14
		}

		for i := 0; i < 5; i++ {
			// Parse asks
			tick.MarketDepth.Asks[i] = DepthLevel{
				Quantity: int64(bigEndianToInt(data[offset : offset+8])),
				Price:    bigEndianToInt(data[offset+8 : offset+12]),
				Orders:   int16(bigEndianToInt(data[offset+12 : offset+14])),
			}
			offset += 14
		}
	}

	return tick, nil
}

// Helper function to convert big endian bytes to int32
func bigEndianToInt(data []byte) int32 {
	var value int32
	buffer := bytes.NewReader(data)
	binary.Read(buffer, binary.BigEndian, &value)
	return value
}

// sendJSONMessage sends a JSON message through the WebSocket connection
func (ws *WS) sendJSONMessage(data interface{}) error {
	if ws.Conn == nil {
		return websocket.ErrCloseSent
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	return ws.Conn.WriteMessage(websocket.TextMessage, jsonData)
}

// reconnect attempts to reconnect to the WebSocket server
func (ws *WS) reconnect() {
	ws.logger.Info().Msg("Attempting to reconnect...")

	if err := ws.Connect(); err != nil {
		ws.logger.Error().Err(err).Msg("Failed to reconnect")
		ws.errChan <- fmt.Errorf("reconnection failed: %w", err)
	}
}

// resubscribeAll resubscribes to all stored subscriptions
func (ws *WS) resubscribeAll() {
	tokensByMode := make(map[string][]int)

	ws.subscriptions.Range(func(key, value interface{}) bool {
		token := key.(int)
		mode := value.(string)
		tokensByMode[mode] = append(tokensByMode[mode], token)
		return true
	})

	for mode, tokens := range tokensByMode {
		if err := ws.Subscribe(tokens, mode); err != nil {
			ws.logger.Error().Err(err).
				Str("mode", mode).
				Interface("tokens", tokens).
				Msg("Failed to resubscribe")
		}
	}
}

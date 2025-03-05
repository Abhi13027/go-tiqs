package tiqs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

type Instrument struct {
	ExchSeg            string  `csv:"ExchSeg,omitempty"`
	Token              int64   `csv:"Token,omitempty"`
	LotSize            int64   `csv:"LotSize,omitempty"`
	Symbol             string  `csv:"Symbol,omitempty"`
	CompanyName        string  `csv:"CompanyName,omitempty"`
	Exchange           string  `csv:"Exchange,omitempty"`
	Segment            string  `csv:"Segment,omitempty"`
	TradingSymbol      string  `csv:"TradingSymbol,omitempty"`
	Instrument         string  `csv:"Instrument,omitempty"`
	ExpiryDate         *string `csv:"ExpiryDate,omitempty"` // Nullable field
	Isin               string  `csv:"Isin,omitempty"`
	TickSize           float64 `csv:"TickSize,omitempty"`
	PricePrecision     int     `csv:"PricePrecision,omitempty"`
	Multiplier         int     `csv:"Multiplier,omitempty"`
	PriceMultiplier    float64 `csv:"PriceMultiplier,omitempty"`
	OptionType         *string `csv:"OptionType,omitempty"` // Nullable field
	UnderlyingExchange *string `csv:"UnderlyingExchange,omitempty"`
	UnderlyingToken    *string `csv:"UnderlyingToken,omitempty"`
	StrikePrice        int64   `csv:"StrikePrice,omitempty"`
	ExchExpiryDate     int64   `csv:"ExchExpiryDate,omitempty"`
	UpdateTime         int64   `csv:"UpdateTime,omitempty"`
	MessageFlag        int     `csv:"MessageFlag,omitempty"`
	FoFlag             int     `csv:"FoFlag,omitempty"`
}

// GetInstrumentList fetches the list of all available instruments.
//
// It sends a GET request to the "/all" endpoint to retrieve a list of all available
// instruments on the platform.
//
// Returns:
//   - A slice of Instrument structs containing all available instruments if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetInstrumentList() ([]Instrument, error) {
	endpoint := "/all"

	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch instrument list")
		return nil, err
	}

	// Preprocess CSV to clean up any malformed lines
	cleanCSV, err := preprocessCSV(resp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to preprocess CSV response")
		return nil, err
	}

	var instruments []Instrument
	if err := gocsv.UnmarshalBytes(cleanCSV, &instruments); err != nil {
		log.Error().Err(err).Msg("Failed to parse CSV response")
		return nil, err
	}

	log.Info().Msg("Successfully parsed instrument list")
	return instruments, nil
}

func preprocessCSV(data []byte) ([]byte, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Allows variable fields per record

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %w", err)
	}

	// Ensure all rows have the correct number of fields
	expectedCols := len(records[0])
	var cleanedRecords [][]string

	for i, row := range records {
		if len(row) == 0 || strings.TrimSpace(strings.Join(row, "")) == "" {
			continue // Skip empty lines
		}
		if len(row) != expectedCols {
			log.Warn().
				Int("line", i+1).
				Int("expected_fields", expectedCols).
				Int("actual_fields", len(row)).
				Str("raw_data", fmt.Sprintf("%q", row)). // Print the actual content
				Msg("Skipping malformed CSV row due to incorrect field count")
			continue // Skip malformed rows
		}
		cleanedRecords = append(cleanedRecords, row)
	}

	// Convert back to CSV format
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.WriteAll(cleanedRecords)
	writer.Flush()

	return buffer.Bytes(), nil
}

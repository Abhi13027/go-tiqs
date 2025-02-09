package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Limits represents the trading limits and margin details for a user.
type Limits struct {
	Data []struct {
		Cash                          string `json:"cash"`
		DayCash                       string `json:"dayCash"`
		BlockedAmount                 string `json:"blockedAmount"`
		UnClearedCash                 string `json:"unClearedCash"`
		BrokerCollateralAmount        string `json:"brokerCollateralAmount"`
		LiquidCollateralAmount        string `json:"liquidCollateralAmount"`
		EquityCollateralAmount        string `json:"equityCollateralAmount"`
		PayIn                         string `json:"payIn"`
		PayOut                        string `json:"payOut"`
		MarginUsed                    string `json:"marginUsed"`
		CashNCarryBuyUsed             string `json:"cashNCarryBuyUsed"`
		CashNCarrySellCredits         string `json:"cashNCarrySellCredits"`
		Turnover                      string `json:"turnover"`
		PendingOrderValue             string `json:"pendingOrderValue"`
		Span                          string `json:"span"`
		Exposure                      string `json:"exposure"`
		DeliveryMargin                string `json:"deliveryMargin"`
		MtomCurrentPct                string `json:"mtomCurrentPct"`
		RealisedPnL                   string `json:"realisedPnL"`
		UnRealisedMtoM                string `json:"unRealisedMtoM"`
		ProductMargin                 string `json:"productMargin"`
		Premium                       string `json:"premium"`
		VarELMMargin                  string `json:"varELMMargin"`
		GrossExposure                 string `json:"grossExposure"`
		GrossExposureDerivate         string `json:"grossExposureDerivate"`
		ScripBasketMargin             string `json:"scripBasketMargin"`
		AdditionalScriptBasketMargin  string `json:"additionalScriptBasketMargin"`
		Brokerage                     string `json:"brokerage"`
		Collateral                    string `json:"collateral"`
		GrossCollateral               string `json:"grossCollateral"`
		TurnOverLimit                 string `json:"turnOverLimit"`
		PendingOrderValueAmount       string `json:"pendingOrderValueAmount"`
		CurrentRealizedPnLei          string `json:"currentRealizedPnLei"`
		CurrentRealizedPnLem          string `json:"currentRealizedPnLem"`
		CurrentRealizedPnLc           string `json:"currentRealizedPnLc"`
		CurrentRealizedPnLdi          string `json:"currentRealizedPnLdi"`
		CurrentRealizedPnLdm          string `json:"currentRealizedPnLdm"`
		CurrentRealizedPnLfi          string `json:"currentRealizedPnLfi"`
		CurrentRealizedPnLfm          string `json:"currentRealizedPnLfm"`
		CurrentRealizedPnLci          string `json:"currentRealizedPnLci"`
		CurrentRealizedPnLcm          string `json:"currentRealizedPnLcm"`
		CurrentUnRealizedPnLei        string `json:"currentUnRealizedPnLei"`
		CurrentUnRealizedPnLem        string `json:"currentUnRealizedPnLem"`
		CurrentUnRealizedPnLc         string `json:"currentUnRealizedPnLc"`
		CurrentUnRealizedPnLdi        string `json:"currentUnRealizedPnLdi"`
		CurrentUnRealizedPnLdm        string `json:"currentUnRealizedPnLdm"`
		CurrentUnRealizedPnLfi        string `json:"currentUnRealizedPnLfi"`
		CurrentUnRealizedPnLfm        string `json:"currentUnRealizedPnLfm"`
		CurrentUnRealizedPnLci        string `json:"currentUnRealizedPnLci"`
		CurrentUnRealizedPnLcm        string `json:"currentUnRealizedPnLcm"`
		SpanDi                        string `json:"spanDi"`
		SpanDm                        string `json:"spanDm"`
		SpanFi                        string `json:"spanFi"`
		SpanFm                        string `json:"spanFm"`
		SpanCi                        string `json:"spanCi"`
		SpanCm                        string `json:"spanCm"`
		ExposureMarginDi              string `json:"exposureMarginDi"`
		ExposureMarginDm              string `json:"exposureMarginDm"`
		ExposureMarginFi              string `json:"exposureMarginFi"`
		ExposureMarginFm              string `json:"exposureMarginFm"`
		ExposureMarginCi              string `json:"exposureMarginCi"`
		ExposureMarginCm              string `json:"exposureMarginCm"`
		PremiumDi                     string `json:"premiumDi"`
		PremiumDm                     string `json:"premiumDm"`
		PremiumFi                     string `json:"premiumFi"`
		PremiumFm                     string `json:"premiumFm"`
		PremiumCi                     string `json:"premiumCi"`
		PremiumCm                     string `json:"premiumCm"`
		VarELMei                      string `json:"varELMei"`
		VarELMem                      string `json:"varELMem"`
		VarELMc                       string `json:"varELMc"`
		CoveredProductMarginEh        string `json:"coveredProductMarginEh"`
		CoveredProductMarginEb        string `json:"coveredProductMarginEb"`
		CoveredProductMarginDh        string `json:"coveredProductMarginDh"`
		CoveredProductMarginDb        string `json:"coveredProductMarginDb"`
		CoveredProductMarginFh        string `json:"coveredProductMarginFh"`
		CoveredProductMarginFb        string `json:"coveredProductMarginFb"`
		CoveredProductMarginCh        string `json:"coveredProductMarginCh"`
		CoveredProductMarginCb        string `json:"coveredProductMarginCb"`
		ScripBasketMarginEi           string `json:"scripBasketMarginEi"`
		ScripBasketMarginEm           string `json:"scripBasketMarginEm"`
		ScripBasketMarginEc           string `json:"scripBasketMarginEc"`
		AdditionalScripBasketMarginDi string `json:"additionalScripBasketMarginDi"`
		AdditionalScripBasketMarginDm string `json:"additionalScripBasketMarginDm"`
		AdditionalScripBasketMarginFi string `json:"additionalScripBasketMarginFi"`
		AdditionalScripBasketMarginFm string `json:"additionalScripBasketMarginFm"`
		AdditionalScripBasketMarginCi string `json:"additionalScripBasketMarginCi"`
		AdditionalScripBasketMarginCm string `json:"additionalScripBasketMarginCm"`
		BrokerageEi                   string `json:"brokerageEi"`
		BrokerageEm                   string `json:"brokerageEm"`
		BrokerageEc                   string `json:"brokerageEc"`
		BrokerageEh                   string `json:"brokerageEh"`
		BrokerageEb                   string `json:"brokerageEb"`
		BrokerageDi                   string `json:"brokerageDi"`
		BrokerageDm                   string `json:"brokerageDm"`
		BrokerageDh                   string `json:"brokerageDh"`
		BrokerageDb                   string `json:"brokerageDb"`
		BrokerageFi                   string `json:"brokerageFi"`
		BrokerageFm                   string `json:"brokerageFm"`
		BrokerageFh                   string `json:"brokerageFh"`
		BrokerageFb                   string `json:"brokerageFb"`
		BrokerageCi                   string `json:"brokerageCi"`
		BrokerageCm                   string `json:"brokerageCm"`
		BrokerageCh                   string `json:"brokerageCh"`
		BrokerageCb                   string `json:"brokerageCb"`
		PeakMargin                    string `json:"peakMargin"`
		RequestTime                   string `json:"requestTime"`
	} `json:"data"`
	Status string `json:"status"`
}

// GetLimits fetches the trading limits and margin details for the authenticated user.
//
// This function sends a GET request to the "/user/limits" endpoint to retrieve available margins,
// blocked funds, collateral, pending orders, and other financial details.
//
// Returns:
//   - A pointer to a Limits struct containing the trading limits if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetLimits() (*Limits, error) {
	endpoint := "/user/limits"

	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch trading limits")
		return nil, err
	}

	var result Limits
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse trading limits response")
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("failed to retrieve trading limits")
	}

	log.Info().Msg("Trading limits retrieved successfully")
	return &result, nil
}

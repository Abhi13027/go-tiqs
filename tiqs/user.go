package tiqs

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// User represents the structure of user details received from the Tiqs API.
type User struct {
	Data struct {
		AccountID   string `json:"accountID"` // Unique identifier for the user's account.
		Address     string `json:"address"`   // Residential address of the user.
		BankDetails []struct {
			Vpa           string `json:"vpa"`           // Virtual Payment Address (UPI).
			BankName      string `json:"bankName"`      // Name of the bank associated with the account.
			AccountType   string `json:"accountType"`   // Type of bank account (e.g., Savings, Current).
			AccountNumber string `json:"accountNumber"` // Account number linked to the user.
		} `json:"bankDetails"` // List of bank accounts associated with the user.

		Blocked       bool   `json:"blocked"` // Indicates if the user's account is blocked.
		City          string `json:"city"`    // City of residence.
		DepositoryIDs struct {
			String string `json:"String"` // Depository participant ID as a string.
			Valid  bool   `json:"Valid"`  // Boolean indicating if the depository ID is valid.
		} `json:"depositoryIDs"` // User's depository information.

		Email       string   `json:"email"`       // User's registered email address.
		Exchanges   []string `json:"exchanges"`   // List of exchanges the user has access to.
		ID          string   `json:"id"`          // Unique identifier for the user.
		Image       string   `json:"image"`       // Profile image URL of the user.
		Name        string   `json:"name"`        // Full name of the user.
		OrdersTypes []string `json:"ordersTypes"` // Types of orders the user can place.
		Pan         string   `json:"pan"`         // Permanent Account Number (PAN) of the user.
		Phone       string   `json:"phone"`       // User's registered phone number.
		Products    []string `json:"products"`    // List of financial products the user has access to.
		State       string   `json:"state"`       // State of residence.
		TotpEnabled bool     `json:"totpEnabled"` // Indicates whether TOTP-based 2FA is enabled.
		UserType    string   `json:"userType"`    // Type of user (e.g., Retail, Institutional).
	} `json:"data"` // Data field containing user details.

	Status string `json:"status"` // Status of the API response (e.g., "success" or "error").
}

// GetUserDetails fetches user profile details from the Tiqs API.
//
// It makes a GET request to the "/user/details" endpoint and returns a User struct
// containing all user-related information.
//
// Returns:
//   - A pointer to a User struct with the retrieved details if successful.
//   - An error if the request fails or the response cannot be parsed.
func (c *Client) GetUserDetails() (*User, error) {
	endpoint := "/user/details"

	// Send a GET request to the API to retrieve user details.
	resp, err := c.request(endpoint, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user profile")
		return nil, err
	}

	var result User
	// Parse the JSON response into the User struct.
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Error().Err(err).Msg("Failed to parse user profile response")
		return nil, err
	}

	// Check if the API response status indicates success.
	if result.Status != "success" {
		return nil, fmt.Errorf("user profile retrieval failed")
	}

	log.Info().Msg("User profile retrieved successfully")
	return &result, nil
}

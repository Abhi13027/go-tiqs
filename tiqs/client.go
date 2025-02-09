// client.go
package tiqs

import (
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// Config holds the SDK configuration settings.
type Config struct {
	AppID        string // Application ID for API authentication.
	AppSecret    string // Application secret key for API authentication.
	Token        string // Authentication token for API requests.
	BaseURL      string // Base URL of the Tiqs API.
	RefreshToken string // Token used to refresh authentication when expired.
}

// Client is the main struct for interacting with the Tiqs API.
//
// It contains the configuration settings and an HTTP client for making API requests.
type Client struct {
	Config     Config           // Configuration settings for the API client.
	HTTPClient *fasthttp.Client // HTTP client for executing requests.
}

// NewClient initializes a new SDK client with the provided application credentials.
//
// Parameters:
//   - appID: The application ID used for authentication.
//   - appSecret: The application secret key used for authentication.
//
// Returns:
//   - A pointer to a newly created Client instance.
func NewClient(appID, appSecret string) *Client {
	return &Client{
		Config: Config{
			AppID:     appID,
			AppSecret: appSecret,
			BaseURL:   "https://api.tiqs.trading",
		},
		HTTPClient: &fasthttp.Client{},
	}
}

// request sends an HTTP API request to the Tiqs server and retrieves the response.
//
// This function constructs an HTTP request with the required authentication headers
// and executes it using the `fasthttp` client.
//
// Parameters:
//   - endpoint: The API endpoint (relative to BaseURL) to send the request to.
//   - method: The HTTP method ("GET" or "POST").
//   - payload: The request body (for POST requests).
//
// Returns:
//   - A byte slice containing the response body if successful.
//   - An error if the request fails.
func (c *Client) request(endpoint string, method string, payload []byte) ([]byte, error) {
	url := c.Config.BaseURL + endpoint
	log.Info().Str("url", url).Msg("Making request")

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.Set("appId", c.Config.AppID)
	req.Header.Set("token", c.Config.Token)

	if method == "POST" {
		req.Header.SetMethod("POST")
		req.SetBody(payload)
	} else {
		req.Header.SetMethod("GET")
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Execute the request using the fasthttp client.
	err := c.HTTPClient.Do(req, resp)
	if err != nil {
		log.Error().Err(err).Msg("API request failed")
		return nil, err
	}

	return resp.Body(), nil
}

// rawRequest sends an HTTP request to a fully specified URL and retrieves the response.
//
// Unlike `request()`, this function allows specifying an absolute URL rather than an endpoint.
//
// Parameters:
//   - url: The full API URL to send the request to.
//   - method: The HTTP method ("GET" or "POST").
//   - payload: The request body (for POST requests).
//
// Returns:
//   - A byte slice containing the response body if successful.
//   - An error if the request fails.
func (c *Client) rawRequest(url string, method string, payload []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	if method == "POST" {
		req.Header.SetMethod("POST")
		req.SetBody(payload)
	} else {
		req.Header.SetMethod("GET")
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Execute the request using the fasthttp client.
	err := c.HTTPClient.Do(req, resp)
	if err != nil {
		log.Error().Err(err).Msg("API request failed")
		return nil, err
	}

	return resp.Body(), nil
}

// SetToken updates the authentication token dynamically.
//
// This function allows updating the API token at runtime without needing to recreate the client.
//
// Parameters:
//   - token: The new authentication token.
func (c *Client) SetToken(token string) {
	c.Config.Token = token
}

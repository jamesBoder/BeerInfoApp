package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// BeerAPI Client struct
type BeerAPIClient struct {
	APIKey  string
	BaseUrl string
	Client  *http.Client
}

// client constructor function
func NewBeerAPIClient(apiKey string) *BeerAPIClient {
	return &BeerAPIClient{
		APIKey:  apiKey,
		BaseUrl: "https://beer9.p.rapidapi.com/",
		Client:  &http.Client{},
	}
}

// SearchBeers method to search for beers by name
func (c *BeerAPIClient) SearchBeers(query string) (models.APIResponse, error) {
	// build the request URL with query parameters
	reqUrl, err := url.Parse(c.BaseUrl)
	if err != nil {
		return models.APIResponse{}, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := reqUrl.Query()
	q.Set("name", query)
	reqUrl.RawQuery = q.Encode()

	// make the HTTP GET request
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return models.APIResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// set required headers
	req.Header.Set("X-RapidAPI-Key", c.APIKey)
	req.Header.Set("X-RapidAPI-Host", "beer9.p.rapidapi.com")

	// send the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return models.APIResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.APIResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	// parse the JSON response
	var apiResp models.APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return models.APIResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return apiResp, nil
}

// helper function for error handling
func (c *BeerAPIClient) handleAPIError(code int) error {
	return fmt.Errorf("API returned error code %d", code)
}

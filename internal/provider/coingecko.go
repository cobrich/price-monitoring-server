package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"price-monitoring-server/internal/models"
	"strings"
	"time"
)

// CoinGeckoProvider fetches real-time cryptocurrency prices from the CoinGecko API.
type CoinGeckoProvider struct {
	client *http.Client
}

// NewCoinGeckoProvider creates a new CoinGeckoProvider.
func NewCoinGeckoProvider() *CoinGeckoProvider {
	return &CoinGeckoProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchPrices gets prices for a list of products from the CoinGecko API in a single request.
func (p *CoinGeckoProvider) FetchPrices(ctx context.Context, products []models.Product) (map[string]int, error) {
	if len(products) == 0 {
		return make(map[string]int), nil
	}

	// Create a comma-separated list of crypto IDs for the API call
	ids := make([]string, len(products))
	for i, p := range products {
		ids[i] = strings.ToLower(p.ProductName)
	}
	idString := strings.Join(ids, ",")

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", idString)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map the results back to the original product names
	prices := make(map[string]int)
	for _, product := range products {
		id := strings.ToLower(product.ProductName)
		if priceData, ok := result[id]; ok {
			if price, ok := priceData["usd"]; ok {
				// Use the original casing from product.ProductName for the map key
				prices[product.ProductName] = int(price)
			}
		}
	}

	return prices, nil
}

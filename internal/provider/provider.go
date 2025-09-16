package provider

import (
	"context"
	"price-monitoring-server/internal/models"
)

// PriceProvider defines the interface for any price data source.
// It fetches prices for a list of products in a single batch.
type PriceProvider interface {
	FetchPrices(ctx context.Context, products []models.Product) (map[string]int, error)
}

package provider

import (
	"context"
	"math/rand"
	"price-monitoring-server/internal/models"
	"time"
)

// RandomProvider generates random price data.
type RandomProvider struct {
	r *rand.Rand
}

// NewRandomProvider creates a new RandomProvider.
func NewRandomProvider() *RandomProvider {
	return &RandomProvider{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// FetchPrices generates a random price for each product in the list.
// It complies with the PriceProvider interface.
func (p *RandomProvider) FetchPrices(_ context.Context, products []models.Product) (map[string]int, error) {
	results := make(map[string]int)
	for _, product := range products {
		if product.MinPrice >= product.MaxPrice {
			results[product.ProductName] = product.MinPrice
			continue
		}
		price := p.r.Intn(product.MaxPrice-product.MinPrice+1) + product.MinPrice
		results[product.ProductName] = price
	}
	return results, nil
}

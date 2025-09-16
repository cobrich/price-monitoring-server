package service

import (
	"context"
	"log"
	"price-monitoring-server/internal/models"
	"sync"
)

type PriceRepository interface {
	SavePrice(ctx context.Context, pp models.ProductPrice) error
	GetPriceHistory(ctx context.Context, productName string) ([]models.PriceHistoryEntry, error)
}

// PriceService is responsible for storing and providing access to the latest price data and statistics.
// It is safe for concurrent use.
type PriceService struct {
	mu           sync.RWMutex
	latestPrices map[string]models.ProductPrice // Key: ProductName
	stats        map[string]models.Stat         // Key: ProductName
	repo         PriceRepository
}

// NewPriceService creates a new PriceService.
func NewPriceService(repo PriceRepository) *PriceService {
	return &PriceService{
		latestPrices: make(map[string]models.ProductPrice),
		stats:        make(map[string]models.Stat),
		repo:         repo,
	}
}

// UpdateWithPrice processes a new price update.
// It updates the latest price for the product and recalculates the statistics.
func (s *PriceService) UpdateWithPrice(pp models.ProductPrice) {
	// First, save to the database
	err := s.repo.SavePrice(context.Background(), pp) // Using background context for simplicity
	if err != nil {
		// In a real app, you'd use a structured logger
		log.Printf("Error saving price to DB: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.latestPrices[pp.ProductName] = pp

	// Update stats
	stat, ok := s.stats[pp.ProductName]
	if !ok {
		stat = models.Stat{
			ProductName: pp.ProductName,
			Min:         pp.ProductPrice,
			Max:         pp.ProductPrice,
		}
	}

	stat.Count++
	stat.Sum += int64(pp.ProductPrice)
	if pp.ProductPrice < stat.Min {
		stat.Min = pp.ProductPrice
	}
	if pp.ProductPrice > stat.Max {
		stat.Max = pp.ProductPrice
	}
	s.stats[pp.ProductName] = stat
}

// GetLatestPrices returns a copy of the latest prices for all products.
func (s *PriceService) GetLatestPrices() map[string]models.ProductPrice {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent race conditions on the caller's side
	pricesCopy := make(map[string]models.ProductPrice, len(s.latestPrices))
	for k, v := range s.latestPrices {
		pricesCopy[k] = v
	}
	return pricesCopy
}

// GetStats returns a copy of the statistics for all products.
func (s *PriceService) GetStats() map[string]models.Stat {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy
	statsCopy := make(map[string]models.Stat, len(s.stats))
	for k, v := range s.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

// GetHistory retrieves the price history for a product from the repository.
func (s *PriceService) GetHistory(ctx context.Context, productName string) ([]models.PriceHistoryEntry, error) {
	return s.repo.GetPriceHistory(ctx, productName)
}

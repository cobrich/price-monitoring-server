package monitoring

import (
	"fmt"
	"price-monitoring-server/internal/models"
)

// PriceUpdater is an interface that can process price updates.
type PriceUpdater interface {
	UpdateWithPrice(pp models.ProductPrice)
}

func ConsumeAndAggregate(updates <-chan models.ProductPrice, updater PriceUpdater) {
	for u := range updates {
		// Log the update
		fmt.Printf("[INFO] Store:%s Product:%s Price:%d\n", u.StoreName, u.ProductName, u.ProductPrice)
		// Update the service with the new price
		updater.UpdateWithPrice(u)
	}
}

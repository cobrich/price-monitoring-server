package monitoring

import (
	"context"
	"log"
	"price-monitoring-server/internal/models"
	"price-monitoring-server/internal/provider"
	"sync"
	"time"
)

func Run(ctx context.Context, cfg models.Config) <-chan models.ProductPrice {
	out := make(chan models.ProductPrice, 100)
	var wg sync.WaitGroup

	providers := map[string]provider.PriceProvider{
		"random":    provider.NewRandomProvider(),
		"coingecko": provider.NewCoinGeckoProvider(),
	}

	// Group products by their provider type
	productsByProvider := make(map[string][]models.Product)
	storeNameMap := make(map[string]string) // Map product name to its store for context

	for storeName, products := range cfg.Stores {
		for productName, prod := range products {
			prod.ProductName = productName
			providerName := prod.Provider
			if providerName == "" {
				providerName = "random"
			}
			productsByProvider[providerName] = append(productsByProvider[providerName], prod)
			storeNameMap[productName] = storeName
		}
	}

	// Launch one worker per provider
	for providerName, productList := range productsByProvider {
		prov, ok := providers[providerName]
		if !ok {
			log.Printf("WARN: Unknown provider '%s'. Skipping.", providerName)
			continue
		}

		wg.Add(1)
		go func(name string, p provider.PriceProvider, prods []models.Product) {
			defer wg.Done()
			runProviderWorker(ctx, name, p, prods, storeNameMap, out)
		}(providerName, prov, productList)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// runProviderWorker manages fetching all prices from a single provider at a set interval.
func runProviderWorker(ctx context.Context, providerName string, p provider.PriceProvider, products []models.Product, storeNames map[string]string, out chan<- models.ProductPrice) {
	var tickerDuration time.Duration
	switch providerName {
	case "coingecko":
		// Fetch once a minute to stay well within rate limits
		tickerDuration = 60 * time.Second
	default: // for "random" and others
		tickerDuration = 2 * time.Second
	}

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	// Perform an initial fetch immediately on startup
	fetchAndSend(ctx, providerName, p, products, storeNames, out)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fetchAndSend(ctx, providerName, p, products, storeNames, out)
		}
	}
}

// fetchAndSend performs a single batch fetch and sends the results to the output channel.
func fetchAndSend(ctx context.Context, providerName string, p provider.PriceProvider, products []models.Product, storeNames map[string]string, out chan<- models.ProductPrice) {
	prices, err := p.FetchPrices(ctx, products)
	if err != nil {
		log.Printf("ERROR fetching prices for provider '%s': %v", providerName, err)
		return
	}

	for productName, price := range prices {
		select {
		case <-ctx.Done():
			return
		case out <- models.ProductPrice{
			StoreName:    storeNames[productName],
			ProductName:  productName,
			ProductPrice: price,
		}:
		}
	}
}

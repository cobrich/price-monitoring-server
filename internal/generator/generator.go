package generator

import (
	"context"
	"math/rand"
	"price-monitoring-server/internal/models"
	"time"
)

// StartProvider запускает генерацию цен для всех продуктов
// cfg - конфиг с магазинами и продуктами
// ctx - контекст, чтобы можно было остановить все горутины
// returns канал, в котором будут появляться цены
func StartProvider(ctx context.Context, cfg models.Config) <-chan models.ProductPrice {
	out := make(chan models.ProductPrice)

	for storeName, products := range cfg.Stores {
		for productName, product := range products {
			// на каждый продукт своя горутина
			go func(store, pName string, p models.Product) {
				for {
					select {
					case <-ctx.Done():
						return // останавливаем генерацию
					default:
						// случайная пауза 200–1000 мс
						sleepTime := time.Duration(rand.Intn(801)+200) * time.Millisecond
						time.Sleep(sleepTime)

						// случайная цена в диапазоне
						price := rand.Intn(p.MaxPrice-p.MinPrice+1) + p.MinPrice

						// пишем результат
						out <- models.ProductPrice{
							StoreName:   store,
							ProductName: pName,
							ProductPrice:    price,
						}
					}
				}
			}(storeName, productName, product)
		}
	}

	return out
}

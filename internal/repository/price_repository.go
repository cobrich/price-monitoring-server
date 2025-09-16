package repository

import (
	"context"
	"price-monitoring-server/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PriceRepository handles database operations for price data.
type PriceRepository struct {
	db *pgxpool.Pool
}

// NewPriceRepository creates a new PriceRepository.
func NewPriceRepository(db *pgxpool.Pool) *PriceRepository {
	return &PriceRepository{db: db}
}

// InitSchema creates the necessary database table if it doesn't exist.
func (r *PriceRepository) InitSchema(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS price_history (
			id SERIAL PRIMARY KEY,
			store_name VARCHAR(255) NOT NULL,
			product_name VARCHAR(255) NOT NULL,
			price INT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

// SavePrice saves a new price record to the database.
func (r *PriceRepository) SavePrice(ctx context.Context, pp models.ProductPrice) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO price_history (store_name, product_name, price) VALUES ($1, $2, $3)",
		pp.StoreName, pp.ProductName, pp.ProductPrice,
	)
	return err
}

// GetPriceHistory retrieves the price history for a specific product.
func (r *PriceRepository) GetPriceHistory(ctx context.Context, productName string) ([]models.PriceHistoryEntry, error) {
	rows, err := r.db.Query(ctx,
		"SELECT store_name, product_name, price, timestamp FROM price_history WHERE product_name = $1 ORDER BY timestamp DESC LIMIT 100",
		productName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PriceHistoryEntry
	for rows.Next() {
		var entry models.PriceHistoryEntry
		if err := rows.Scan(&entry.StoreName, &entry.ProductName, &entry.ProductPrice, &entry.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}

	return history, nil
}

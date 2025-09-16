package models

import "time"

type Config struct {
	Stores map[string]map[string]Product `json:"stores"`
}

type Product struct {
	ProductName    string
	Provider       string `json:"provider"` // e.g., "random" or "coingecko"
	MinPrice       int    `json:"minPrice"`
	MaxPrice       int    `json:"maxPrice"`
	AlertThreshold int    `json:"alertThreshold"`
}

type ProductPrice struct {
	StoreName    string
	ProductName  string
	ProductPrice int
}

type PriceHistoryEntry struct {
	StoreName    string    `json:"storeName"`
	ProductName  string    `json:"productName"`
	ProductPrice int       `json:"productPrice"`
	Timestamp    time.Time `json:"timestamp"`
}

type ProductAvarage struct {
	ProductName string
	Avg         float32
}

type Stat struct {
	ProductName string
	Count       int
	Sum         int64
	Min         int
	Max         int
}

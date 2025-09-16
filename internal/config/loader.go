package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"price-monitoring-server/internal/models"
)

func readConf(data []byte) models.Config {
	var cfg models.Config
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		slog.Error(err.Error())
		return cfg
	}

	return cfg
}

func LoadData(path string) (models.Config, error) {
	// Reading all datas from jsonfile
	data, err := os.ReadFile(path)
	if err != nil {
		return models.Config{}, err
	}

	// Reading stores
	cfg := readConf(data)

	return cfg, nil
}

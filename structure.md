
---

# 📂 Структура проекта

```
price-monitoring/
├── cmd/
│   └── app/
│       └── main.go        # Точка входа
├── config/
│   └── config.json        # JSON с магазинами и продуктами
├── internal/
│   ├── config/
│   │   └── loader.go      # Загрузка и парсинг конфигурации
│   ├── generator/
│   │   └── generator.go   # Генерация случайных цен
│   ├── monitoring/
│   │   ├── aggregator.go  # Читает из каналов, пишет логи
│   │   ├── worker.go      # Горутины-поставщики (продукты/магазины)
│   │   └── history.go     # Сохраняет в history.json
│   ├── models/
│   │   └── models.go      # Структуры данных (Config, ProductConfig, PriceUpdate)
│   └── analytics/
│       └── analytics.go   # Подсчёт средней/мин/макс цены
├── logs/
│   └── monitoring.log     # Логи сервиса
├── history/
│   └── history.json       # Архив всех цен
├── tests/
│   ├── config_test.go     # Тесты загрузки конфигурации
│   ├── generator_test.go  # Тесты генерации цен
│   └── monitoring_test.go # Тесты агрегации и порогов
├── Dockerfile             # Для запуска в контейнере (по желанию)
├── Makefile               # Удобные команды (run, test, build)
└── go.mod
```

---

# 📌 Что где лежит

* **cmd/app/main.go**
  Точка входа → инициализация логов, загрузка конфига, запуск мониторинга, graceful shutdown через `context.WithTimeout`.

* **config/config.json**
  Реалистичные данные (магазины, продукты, цены).

* **internal/config/loader.go**
  Функция `LoadConfig(path string) (Config, error)` → парсит JSON.

* **internal/generator/generator.go**
  Функция `GeneratePrice(min, max float64) float64` → возвращает случайную цену в диапазоне.

* **internal/monitoring/worker.go**
  Каждая горутина генерирует цену для продукта и пишет в канал.

* **internal/monitoring/aggregator.go**
  Слушает канал, пишет логи, вызывает тревогу если цена ниже `alertThreshold`.

* **internal/monitoring/history.go**
  Сохраняет все события в JSON-файл.

* **internal/analytics/analytics.go**
  После завершения работы считает среднее, min, max для каждого продукта в каждом магазине.

* **internal/models/models.go**
  Все структуры:

  ```go
  type ProductConfig struct {
      MinPrice       float64 `json:"minPrice"`
      MaxPrice       float64 `json:"maxPrice"`
      AlertThreshold float64 `json:"alertThreshold"`
  }

  type Config struct {
      Stores map[string]map[string]ProductConfig `json:"stores"`
  }

  type PriceUpdate struct {
      Store   string  `json:"store"`
      Product string  `json:"product"`
      Price   float64 `json:"price"`
      Time    string  `json:"time"`
  }
  ```

* **tests/**
  Покрываешь тестами генерацию, загрузку конфигов, обработку тревог.

* **logs/** и **history/**
  Отдельные папки для логов и истории, чтобы красиво смотрелось в портфолио.

---

# ✨ Как показать в портфолио

* **GitHub README.md**

  * Краткое описание проекта.
  * Как запустить (`make run`, `docker build`).
  * Скриншоты логов и `history.json`.
  * Вывод аналитики (например, таблица со средними ценами).
* **Dockerfile** — сразу даёт +крутость (можно сказать: "готов к деплою").
* **Makefile** — удобные команды (`make test`, `make run`).
* **Tests** — обязательно, это всегда ценят.

---
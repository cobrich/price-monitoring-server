
---

# 📌 ТЗ: «Сервис мониторинга цен в магазинах»

## 🎯 Цель проекта

Реализовать мини-сервис, который имитирует мониторинг цен на продукты в разных магазинах.
Система должна уметь:

* читать конфигурацию магазинов и продуктов из JSON,
* запускать параллельных «поставщиков цен» для каждого продукта в каждом магазине,
* собирать цены в общий канал и логировать,
* выявлять аномалии (низкие цены),
* работать ограниченное время и корректно завершаться,
* быть покрытой тестами,
* выдавать как логи, так и итоговую аналитику.

---

## 🔨 Основные требования

1. **Чтение конфигурации из JSON**

   * Файл `config.json` хранит список магазинов.
   * У каждого магазина есть список продуктов.
   * У каждого продукта:

     * `minPrice`, `maxPrice` — диапазон генерации цены,
     * `alertThreshold` — порог для предупреждения.

2. **Генерация данных (горутины)**

   * Для каждого продукта каждого магазина создаётся отдельная горутина.
   * Она генерирует случайную цену каждые 200–1000 мс и пишет в канал.

3. **Агрегатор (каналы)**

   * Все данные собираются в общий буферизованный канал.
   * Агрегатор логирует события:

     * `[INFO] Store: StoreA, Product: milk, Price: 270.00, Time: 2025-09-04 11:00:01`
     * `[ALERT] Store: StoreB, Product: water, Price too low: 60.00`

4. **Управление временем работы (context)**

   * Использовать `context.WithTimeout`, например 15 секунд.
   * После таймаута сервис должен корректно завершить все горутины.

5. **Логирование**

   * Использовать `log` с разными уровнями (`INFO`, `ALERT`).
   * Логи писать в консоль и файл `monitoring.log`.

6. **Сохранение истории**

   * Все события также сохранять в JSON-файл `history.json`.

7. **Аналитика после завершения**

   * Посчитать и вывести:

     * среднюю цену каждого продукта в каждом магазине,
     * минимальную и максимальную цену, зафиксированную за время работы.

8. **Тестирование**

   * Протестировать функцию генерации цены (цена всегда в диапазоне).
   * Протестировать корректность загрузки JSON.
   * Протестировать, что тревога срабатывает, если цена ниже порога.

---

## ✨ Дополнительные возможности (если захочешь «звёздочку» в портфолио)

* Добавить REST API (например, `GET /average-prices`) с выводом итоговой аналитики в JSON.
* Вынести конфигурацию в отдельный пакет.
* Сделать Dockerfile и запустить проект в контейнере.
* Добавить Makefile (`make run`, `make test`).
* Визуализировать логи или средние цены через Grafana/Prometheus (если хочешь «вау»-эффект).

---

## 📂 Пример большого JSON с реальными магазинами и продуктами

Файл `config.json`:

```json
{
  "stores": {
    "StoreA": {
      "milk": { "minPrice": 260, "maxPrice": 400, "alertThreshold": 280 },
      "bread": { "minPrice": 85, "maxPrice": 150, "alertThreshold": 100 },
      "butter": { "minPrice": 600, "maxPrice": 950, "alertThreshold": 700 },
      "cheese": { "minPrice": 950, "maxPrice": 1400, "alertThreshold": 1000 }
    },
    "StoreB": {
      "milk": { "minPrice": 240, "maxPrice": 420, "alertThreshold": 270 },
      "water": { "minPrice": 55, "maxPrice": 110, "alertThreshold": 65 },
      "eggs": { "minPrice": 310, "maxPrice": 480, "alertThreshold": 340 },
      "bread": { "minPrice": 90, "maxPrice": 160, "alertThreshold": 105 }
    },
    "StoreC": {
      "eggs": { "minPrice": 300, "maxPrice": 500, "alertThreshold": 350 },
      "cheese": { "minPrice": 1000, "maxPrice": 1600, "alertThreshold": 1100 },
      "yogurt": { "minPrice": 150, "maxPrice": 300, "alertThreshold": 180 },
      "juice": { "minPrice": 200, "maxPrice": 450, "alertThreshold": 220 }
    },
    "StoreD": {
      "chocolate": { "minPrice": 300, "maxPrice": 700, "alertThreshold": 350 },
      "coffee": { "minPrice": 1000, "maxPrice": 2000, "alertThreshold": 1200 },
      "tea": { "minPrice": 400, "maxPrice": 900, "alertThreshold": 450 },
      "sugar": { "minPrice": 200, "maxPrice": 500, "alertThreshold": 250 }
    },
    "StoreE": {
      "meat": { "minPrice": 2000, "maxPrice": 3500, "alertThreshold": 2200 },
      "fish": { "minPrice": 1500, "maxPrice": 2800, "alertThreshold": 1800 },
      "rice": { "minPrice": 300, "maxPrice": 700, "alertThreshold": 400 },
      "pasta": { "minPrice": 250, "maxPrice": 600, "alertThreshold": 300 }
    }
  }
}
```

---

## 📊 Как это будет смотреться в портфолио

* Название: **Price Monitoring System (Go, Concurrency, JSON, Context)**
* Краткое описание: «Сервис для мониторинга цен на продукты в разных магазинах с использованием горутин, каналов, контекста, логирования, JSON и тестирования».
* Технологии: Go, goroutines, channels, context, logging, JSON, testing.
* Опционально: Docker, REST API, Prometheus/Grafana.

---
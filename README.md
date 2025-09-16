# Price Monitoring Server

This is a powerful backend service written in Go that monitors real-time price data for various products from multiple sources, stores the history in a database, and serves the data via a REST API.

The application is built with a modular, provider-based architecture, allowing it to fetch data from any source, such as public APIs (like CoinGecko for cryptocurrencies) or internal random generators for simulation.

## Features

- **Real-time Monitoring**: Concurrently fetches data from multiple sources at configurable intervals.
- **REST API**: Exposes endpoints to get the latest prices, real-time statistics, and historical data. Built with [Gin](https://github.com/gin-gonic/gin).
- **Database Persistence**: Stores all price history in a PostgreSQL database for later analysis.
- **Provider Architecture**: Easily extensible to support new data sources (e.g., other APIs, web scrapers).
- **Containerized**: Fully containerized with Docker and Docker Compose for easy setup and deployment.

## How It Works

1.  **Configuration**: The service reads a `config.json` file to determine which products to monitor. For each product, a `provider` is specified (e.g., `coingecko` or `random`).
2.  **Provider Workers**: The application starts one worker goroutine for each unique data provider (not for each product).
3.  **Batched Fetching**: Each provider worker fetches data for all its assigned products in a single batch request at a set interval (e.g., the CoinGecko worker fetches all crypto prices once per minute).
4.  **Data Processing**: The collected data is sent through a channel to a central service that updates both the in-memory cache (for latest prices/stats) and saves a record in the PostgreSQL database.
5.  **API Layer**: The Gin-based API serves data from the in-memory cache for instant access and queries the database for historical data.

## How to Run

The easiest way to get the entire stack (the application and the PostgreSQL database) running is with Docker Compose.

1.  **Prerequisites**: Make sure you have Docker and Docker Compose installed.
2.  **Clone the repository**:
    ```bash
    git clone <your-repo-url>
    cd price-monitoring-server
    ```
3.  **Run the application**:
    ```bash
    docker-compose up --build
    ```
The `--build` flag is only necessary the first time or after you've made code changes. The server will be available at `http://localhost:8080`.

## API Endpoints

- `GET /health`: Health check endpoint. Returns `{"status": "UP"}`.
- `GET /prices`: Get the latest price for all monitored products.
- `GET /stats`: Get real-time analytics (min, max, count, sum) for all products.
- `GET /history/:productName`: Get the price history for a specific product (e.g., `/history/bitcoin`).

## Project Structure

```
/
├── cmd/app/main.go         # Main application entry point
├── config/config.json      # Configuration file for stores and products
├── internal/
│   ├── models/             # Go structs for data models
│   ├── monitoring/         # Core logic for monitoring workers
│   ├── provider/           # Price provider implementations (random, coingecko)
│   ├── repository/         # Database interaction logic (PostgreSQL)
│   └── service/            # Business logic and in-memory state management
├── Dockerfile              # Dockerfile for building the Go application
├── docker-compose.yml      # Docker Compose file for running the app and DB
└── go.mod
```

## Future Improvements

- **WebSockets**: Implement a WebSocket endpoint to push real-time price updates to clients.
- **Frontend Dashboard**: Build a simple frontend dashboard to visualize the price data and statistics.
- **More Providers**: Add more providers for different data sources (e.g., stock market APIs, other crypto exchanges).
- **Alerting**: Implement a system to send alerts (e.g., via email or Slack) when a price crosses a predefined threshold.

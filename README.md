# Quanta Conduit (Core API Service Monorepo)

This monorepo contains the core API services for our trading platform, including both Go and Python components. It provides a centralized location for all application APIs, leveraging the strengths of both Go and Python for different aspects of the system.

## Project Structure

```
/core-api-monorepo
├── /go-services
│   ├── /cmd
│   │   └── /api
│   │       └── main.go
│   ├── /internal
│   │   ├── /auth
│   │   ├── /handlers
│   │   ├── /models
│   │   └── /database
│   ├── /pkg
│   │   ├── /logger
│   │   └── /config
│   ├── go.mod
│   └── go.sum
├── /python-services
│   ├── /trading_strategies
│   │   ├── __init__.py
│   │   └── volatility_breakout.py
│   ├── /api
│   │   ├── __init__.py
│   │   └── main.py
│   ├── requirements.txt
│   └── setup.py
├── /shared
│   ├── /config
│   │   ├── default.yaml
│   │   └── development.yaml
│   │   └── production.yaml
│   ├── /utils
│   │   ├── logger.py
│   │   └── logger.go
│   ├── /models
│   │   ├── user.go
│   │   └── user.py
│   ├── /constants
│   │   ├── error_codes.py
│   │   └── error_codes.go
├── /scripts
│   ├── build.sh
│   └── test.sh
├── /docs
│   ├── api_spec.md
│   └── architecture_overview.md
├── docker-compose.yml
├── .gitignore
└── README.md
```

## Components

### Go Services

The Go services form the backbone of our core API, handling high-performance, low-latency operations.

- **cmd/api**: Main entry point for the Go API service.
- **internal**: Internal packages for the Go service, including authentication, request handlers, data models, and database interactions.
- **pkg**: Shared packages that can be used by other services, including our custom logger and configuration management.

### Python Services

The Python services handle our trading strategies and provide a FastAPI interface for strategy-related operations.

- **trading_strategies**: Implementation of various trading strategies.
- **api**: FastAPI application for exposing trading strategy operations.

### Shared Components

- **proto**: Protocol Buffer definitions for any gRPC services (if needed in the future).
- **scripts**: Build and test scripts for the entire monorepo.
- **docs**: Project documentation, including API specifications and architecture overviews.

## Key Features

1. **Unified API Layer**: Combines Go and Python services under a single API gateway.
2. **High-Performance Core**: Utilizes Go for core services requiring high performance and low latency.
3. **Flexible Trading Strategies**: Leverages Python for implementing and managing trading strategies.
4. **Custom Logging**: Integrates with our custom logging service for enhanced observability across all microservices.
5. **Scalable Architecture**: Designed to scale horizontally to handle increasing load.
6. **Secure Authentication**: Implements JWT-based authentication for secure inter-service communication.
7. **Efficient Data Management**: Utilizes PostgreSQL with connection pooling and Redis caching for optimal data access.

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/your-org/core-api-monorepo.git
   ```

2. Set up Go services:
   ```
   cd go-services
   go mod download
   ```

3. Set up Python services:
   ```
   cd ../python-services
   pip install -r requirements.txt
   ```

4. Start the services:
   ```
   docker-compose up
   ```

## Development

- For Go service development, navigate to the `go-services` directory and run `go run cmd/api/main.go`.
- For Python service development, navigate to the `python-services` directory and run `uvicorn api.main:app --reload`.

## Testing

Run the test script from the root of the monorepo:
```
./scripts/test.sh
```

## Deployment

The `docker-compose.yml` file provides a basic setup for local development and testing. For production deployment, consider using Kubernetes for orchestration and management of the microservices.

## Contributing

Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
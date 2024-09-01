## API Specification

### Core API Service (Go)

Base URL: `https://api.ourtradingplatform.com/v1`

#### Authentication

```
POST /auth/login
Request:
{
  "username": string,
  "password": string
}
Response:
{
  "token": string,
  "expires_at": timestamp
}

POST /auth/refresh
Request:
{
  "refresh_token": string
}
Response:
{
  "token": string,
  "expires_at": timestamp
}

POST /auth/logout
Request:
{
  "token": string
}
Response:
{
  "message": "Logged out successfully"
}
```

#### Account Management

```
GET /account
Response:
{
  "id": string,
  "balance": number,
  "created_at": timestamp,
  "updated_at": timestamp
}

PUT /account/deposit
Request:
{
  "amount": number
}
Response:
{
  "new_balance": number,
  "transaction_id": string
}

PUT /account/withdraw
Request:
{
  "amount": number
}
Response:
{
  "new_balance": number,
  "transaction_id": string
}
```

#### Portfolio Management

```
GET /portfolio
Response:
{
  "positions": [
    {
      "symbol": string,
      "quantity": number,
      "average_price": number
    }
  ],
  "total_value": number
}

POST /portfolio/position
Request:
{
  "symbol": string,
  "quantity": number,
  "price": number
}
Response:
{
  "transaction_id": string,
  "new_position": {
    "symbol": string,
    "quantity": number,
    "average_price": number
  }
}
```

### Trading Strategy Service (Python)

Base URL: `https://api.ourtradingplatform.com/v1/strategy`

```
POST /create
Request:
{
  "name": string,
  "type": string,
  "parameters": {
    "param1": value1,
    "param2": value2
  }
}
Response:
{
  "strategy_id": string,
  "name": string,
  "type": string,
  "parameters": object
}

GET /list
Response:
{
  "strategies": [
    {
      "id": string,
      "name": string,
      "type": string
    }
  ]
}

GET /{strategy_id}
Response:
{
  "id": string,
  "name": string,
  "type": string,
  "parameters": object,
  "performance": {
    "total_return": number,
    "sharpe_ratio": number
  }
}

POST /{strategy_id}/backtest
Request:
{
  "start_date": date,
  "end_date": date,
  "initial_capital": number
}
Response:
{
  "backtest_id": string,
  "summary": {
    "total_return": number,
    "sharpe_ratio": number,
    "max_drawdown": number
  }
}

POST /{strategy_id}/activate
Response:
{
  "message": "Strategy activated successfully",
  "active_since": timestamp
}

POST /{strategy_id}/deactivate
Response:
{
  "message": "Strategy deactivated successfully",
  "active_until": timestamp
}
```

### Websocket Endpoints

Base URL: `wss://api.ourtradingplatform.com/v1/ws`

- `/account`: Real-time account balance updates
- `/portfolio`: Real-time portfolio position updates
- `/market`: Real-time market data updates (if applicable)

This architecture and API specification provides a solid foundation for your trading platform, combining high-performance Go services with flexible Python-based trading strategies, all integrated with your existing Kafka and Rust components. The API design allows for comprehensive management of user accounts, portfolios, and trading strategies, while the websocket endpoints enable real-time updates for a responsive user experience.
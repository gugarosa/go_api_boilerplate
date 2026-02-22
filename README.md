# Go-API Boilerplate

A production-ready RESTful API boilerplate built with **Go**, **Gin**, **MongoDB**, and **Redis**. Ships fully containerized with Docker, JWT authentication (access + refresh tokens), hot-reloading for development, and a comprehensive test suite.

**Go 1.23+** · **MIT License**

---

## Features

- **JWT Authentication** — Dual-token system (short-lived access + long-lived refresh) with Redis-backed session management and token revocation
- **5 Example CRUD Resources** — Category, Product, Question, Survey, and Tag with full Create/Read/Update/Delete endpoints
- **MongoDB Aggregation** — `$lookup` pipelines for resolving entity relationships (e.g., products → tags + categories)
- **Dockerized Stack** — One-command setup with Docker Compose (API + MongoDB + Redis)
- **Hot-Reloading** — Automatic recompilation on file changes via Reflex
- **Graceful Shutdown** — Signal-aware server with read/write/idle timeouts
- **Test Suite** — Unit, integration/E2E, and stress/benchmark tests

---

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### 1. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your desired settings:

```env
PORT=8080
MODE=debug
DB_USER=admin
DB_PASS=password
DB_NAME=boilerplate
DB_HOST=db
DB_PORT=27017
REDIS_PASS=password
REDIS_HOST=cache
REDIS_PORT=6379
ACCESS_SECRET=your-access-secret
REFRESH_SECRET=your-refresh-secret
```

### 2. Build and run

```bash
docker-compose build
docker-compose up -d
```

The API will be available at `http://localhost:8080`.

### 3. Verify it works

```bash
# Register a user
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'

# Login
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

---

## API Endpoints

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/register` | No | Create a new user |
| POST | `/v1/login` | No | Login and receive tokens |
| POST | `/v1/refresh` | No | Refresh access token |
| POST | `/v1/logout` | Yes | Invalidate current session |

### Resources (Category, Product, Question, Survey, Tag)

Each resource follows the same pattern (shown for `/v1/category`):

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/v1/category/` | No | List all |
| GET | `/v1/category/:id` | No | Find by ID |
| POST | `/v1/category/` | Yes | Create |
| PATCH | `/v1/category/:id` | Yes | Update |
| DELETE | `/v1/category/:id` | Yes | Delete |

Replace `category` with `product`, `question`, `survey`, or `tag`.

---

## Project Structure

```
go_api_boilerplate/
├── docker-compose.yml          # Development stack
├── docker-compose.prod.yml     # Production stack
├── docker-compose.test.yml     # Test stack
├── Dockerfile                  # Dev image (Go 1.23 + Reflex)
├── Dockerfile.prod             # Production multi-stage image
├── .env.example                # Environment variable template
├── requests.json               # Postman collection
└── src/
    ├── api.go                  # Entry point
    ├── go.mod / go.sum         # Go modules
    ├── controllers/            # Route handlers (auth, category, product, ...)
    ├── database/               # MongoDB + Redis operations
    ├── middleware/              # JWT authentication guard
    ├── models/                 # Data structures
    ├── server/                 # Router and server initialization
    ├── utils/                  # Constants, logger, responses, validators
    └── tests/                  # Integration, E2E, and stress tests
```

For a detailed architecture overview, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Testing

The project includes a three-layer test suite:

- **Unit tests** — Logger, responses, validators, JWT middleware, BSON helpers (no external deps)
- **Integration/E2E tests** — Full auth flow, CRUD operations, edge cases (requires Docker)
- **Stress tests** — Concurrent reads/writes, high-volume creation, rapid login/logout cycles

Run all tests:

```bash
docker-compose -f docker-compose.test.yml build
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
docker-compose -f docker-compose.test.yml down
```

---

## Production Deployment

```bash
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d
```

The production image uses a multi-stage build that compiles to a minimal `scratch` container. Make sure to set strong secrets for `ACCESS_SECRET` and `REFRESH_SECRET` in your `.env` file.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.23 |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) v1.11 |
| Database | [MongoDB](https://www.mongodb.com/) (mongo-driver v1.17) |
| Cache / Sessions | [Redis](https://redis.io/) (go-redis v9.17) |
| Authentication | [golang-jwt](https://github.com/golang-jwt/jwt) v5 |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) v10 |
| Hot-Reload | [Reflex](https://github.com/cespare/reflex) |
| Containers | Docker + Docker Compose |

---

## License

[MIT](LICENSE)

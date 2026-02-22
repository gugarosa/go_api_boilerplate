# Go-API Boilerplate — Architecture Guide

> **Go:** 1.14+ · **License:** MIT · **Framework:** Gin
> A containerized RESTful API boilerplate with JWT authentication, MongoDB, and Redis.

---

## Table of Contents

1. [Overview](#1-overview)
2. [High-Level Architecture](#2-high-level-architecture)
3. [Infrastructure & Deployment](#3-infrastructure--deployment)
4. [Entry Point](#4-entry-point-apigo)
5. [Server Module](#5-server-module-server)
6. [Middleware Module](#6-middleware-module-middleware)
7. [Controllers Module](#7-controllers-module-controllers)
8. [Database Module](#8-database-module-database)
9. [Models Module](#9-models-module-models)
10. [Utils Module](#10-utils-module-utils)
11. [API Endpoints Reference](#11-api-endpoints-reference)
12. [Key Dependencies](#12-key-dependencies)
13. [Project Structure](#13-project-structure)
14. [Design Decisions](#14-design-decisions)
15. [Security Considerations](#15-security-considerations)
16. [Extending the Boilerplate](#16-extending-the-boilerplate)

---

## 1. Overview

Go-API Boilerplate is a ready-to-use foundation for building RESTful APIs in Go. It ships a fully containerized stack (Docker Compose) with **MongoDB** for persistent storage, **Redis** for token caching, and the **Gin** web framework for HTTP routing. The project includes a pre-built JWT-based authentication system (access + refresh tokens), hot-reloading for development via **Reflex**, and five example CRUD resources (**Category**, **Product**, **Question**, **Survey**, **Tag**) that demonstrate the intended patterns for extension.

The architecture follows a **layered MVC-like pattern**: incoming HTTP requests flow through the **router** → **middleware** (authentication guards) → **controllers** (business logic) → **database** (data access), with **models** defining the data structures and **utils** providing cross-cutting concerns.

---

## 2. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          api.go (entry point)                           │
│  Reads env config → Initializes MongoDB → Initializes Redis → Starts   │
│  the Gin HTTP server                                                    │
├───────────────┬───────────────┬───────────────┬─────────────────────────┤
│    server     │  controllers  │   database    │   utils                 │
│  (router +   │  (route       │  (MongoDB +   │  (constants, logger,    │
│   server     │   handlers)   │   Redis ops)  │   responses, validators)│
│   init)      │               │               │                         │
├───────────────┤               ├───────────────┤                         │
│  middleware   │               │    models     │                         │
│  (JWT auth   │               │  (data        │                         │
│   guards)    │               │   structs)    │                         │
└───────────────┴───────────────┴───────────────┴─────────────────────────┘
```

### Request Lifecycle

```
Client Request
      │
      ▼
┌─────────────┐
│   Gin       │
│   Router    │
│   /v1/*     │
└──────┬──────┘
       │
       ├──── Public routes (GET list, GET find) ─────────────────┐
       │                                                         │
       └──── Protected routes (POST, PATCH, DELETE) ──┐          │
                                                      ▼          │
                                               ┌──────────────┐  │
                                               │  AuthGuard   │  │
                                               │  (JWT verify)│  │
                                               └──────┬───────┘  │
                                                      │          │
                                                      ▼          ▼
                                               ┌──────────────────────┐
                                               │  Controller Handler  │
                                               │  (+ AuthRequest for  │
                                               │   Redis session check│
                                               │   on protected ops)  │
                                               └──────────┬───────────┘
                                                          │
                                                          ▼
                                               ┌──────────────────────┐
                                               │  Database Layer      │
                                               │  (MongoDB / Redis)   │
                                               └──────────┬───────────┘
                                                          │
                                                          ▼
                                                    JSON Response
```

---

## 3. Infrastructure & Deployment

### 3.1 Docker Compose Stack

The application runs as a three-service Docker Compose stack (Compose file version `3.8`):

| Service | Image | Purpose | Default Port |
|---------|-------|---------|-------------|
| `api` | Custom (Golang Alpine) | Go API application | `8080` |
| `db` | `mongo` | MongoDB document database | `27017` |
| `cache` | `redis` | Redis in-memory cache for JWT tokens | `6379` |

All services are configured with `restart: always` in production, ensuring automatic recovery after crashes or host reboots.

### 3.2 Data Persistence

Both stateful services use Docker volumes for data durability across container restarts:

| Service | Volume Mount | Purpose |
|---------|-------------|---------|
| `db` (MongoDB) | `./storage/db:/data/db` | Persists all database documents |
| `cache` (Redis) | `./storage/redis:/data` | Persists AOF (append-only file) for Redis data recovery |

Redis is configured with `--appendonly yes`, ensuring write operations are logged to disk.

### 3.3 Development vs. Production

| Aspect | Development (`Dockerfile`) | Production (`Dockerfile.prod`) |
|--------|---------------------------|-------------------------------|
| Base image | `golang:1.14-alpine` | Multi-stage: `golang:1.14-alpine` → `scratch` |
| Hot-reload | Yes, via Reflex (`reflex -c ./reflex.conf`) | No |
| Source volumes | Mounted (`./src:/src`) for live editing | Not mounted (binary copied at build) |
| Binary | Compiled on-the-fly via `go run .` (triggered by Reflex) | Compiled via `go build -o api` |
| Final image size | Full Go toolchain | Minimal (scratch + static binary) |

### 3.4 Environment Configuration

All configuration is driven by environment variables (`.env` file), passed to containers via Docker Compose:

| Variable | Purpose |
|----------|---------|
| `PORT` | HTTP server port |
| `MODE` | Gin mode (`debug` / `release` / `test`) |
| `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_HOST`, `DB_PORT` | MongoDB connection parameters |
| `REDIS_PASS`, `REDIS_HOST`, `REDIS_PORT` | Redis connection parameters |
| `ACCESS_SECRET` | HMAC signing key for JWT access tokens |
| `REFRESH_SECRET` | HMAC signing key for JWT refresh tokens |

---

## 4. Entry Point (`api.go`)

The `main()` function orchestrates the application bootstrap in four sequential steps:

1. **`getConfig()`** — reads all environment variables into a `map[string]string`
2. **`database.InitMongo()`** — connects to MongoDB, pings to verify, and initializes all collection references
3. **`database.InitRedis()`** — connects to Redis, pings to verify
4. **`server.InitServer()`** — sets the Gin mode, creates the engine, registers all routes, and starts the HTTP listener

Both database initializations use `utils.LogFatalError()` — if either MongoDB or Redis is unreachable, the application exits immediately rather than starting in a degraded state.

---

## 5. Server Module (`server/`)

### 5.1 Server Initialization (`server.go`)

Sets the Gin operation mode (e.g., `debug`, `release`) via `gin.SetMode()`, creates a default Gin engine via `gin.Default()` — which automatically includes the **Logger** (request logging) and **Recovery** (panic recovery) middleware — delegates to `InitRouter()`, and then calls `r.Run()` to start listening on the configured `PORT`.

### 5.2 Router (`router.go`)

All API routes are grouped under the `/v1` prefix. Each resource controller registers its own routes via a `CreateRoutes(r *gin.RouterGroup)` function:

| Controller | Route Group | Description |
|-----------|-------------|-------------|
| `auth` | `/v1/` | Authentication (login, register, refresh, logout) |
| `category` | `/v1/category` | Category CRUD |
| `product` | `/v1/product` | Product CRUD |
| `question` | `/v1/question` | Question CRUD |
| `survey` | `/v1/survey` | Survey CRUD |
| `tag` | `/v1/tag` | Tag CRUD |

A catch-all `NoRoute` handler returns `404` for undefined endpoints.

---

## 6. Middleware Module (`middleware/`)

The middleware module handles all JWT token operations.

### 6.1 Authentication (`authentication.go`)

Provides the full JWT lifecycle:

| Function | Purpose |
|----------|---------|
| `CreateToken(id)` | Generates an access token (15 min TTL) and refresh token (7 day TTL), each with a UUID. Signs both with HMAC-SHA256. **Note:** These TTLs are hardcoded in the middleware, not configurable via environment variables |
| `GetToken(r)` | Extracts the Bearer token from the `Authorization` header |
| `VerifyToken(r)` | Parses and validates the JWT signature against `ACCESS_SECRET` |
| `GetTokenData(r)` | Extracts `access_uuid` and `user_id` claims from a verified token, returns a `RedisAccess` struct |
| `ValidateToken(r)` | Lightweight check that a token is structurally valid |
| `AuthGuard()` | Gin middleware handler — validates the token and aborts with `401` on failure |

### 6.2 Refresh (`refresh.go`)

Handles the refresh token flow:

| Function | Purpose |
|----------|---------|
| `GetRefreshToken(c)` | Reads `refresh_token` from the JSON request body |
| `VerifyRefreshToken(c)` | Parses and validates the refresh JWT against `REFRESH_SECRET` |
| `GetRefreshTokenData(c)` | Extracts `refresh_uuid` and `user_id` (as `ObjectID`) from refresh token claims |

### 6.3 Authentication Flow

```
Register:  POST /v1/register  →  Hash password (bcrypt)  →  Store user in MongoDB
                               →  Returns success message (no auto-login; client must call /login)

Login:     POST /v1/login     →  Verify email/password   →  Create access + refresh tokens
                               →  Cache both UUIDs in Redis (with TTL)
                               →  Return tokens to client

Request:   Any protected route →  AuthGuard extracts Bearer token
                               →  Verify JWT signature
                               →  Check access_uuid exists in Redis
                               →  Allow or reject (401)

Refresh:   POST /v1/refresh   →  Verify refresh token    →  Delete old refresh UUID from Redis
                               →  Create new token pair   →  Cache new UUIDs in Redis

Logout:    POST /v1/logout    →  Delete access_uuid from Redis (invalidates token)
```

---

## 7. Controllers Module (`controllers/`)

### 7.1 Common Utilities (`common.go`)

Shared controller logic used across all resource controllers:

| Function | Purpose |
|----------|---------|
| `AuthRequest(c)` | Authenticates a request by extracting token metadata and verifying it against Redis |
| `BindAndValidateRequest(c, model)` | Binds JSON body to a Go struct and runs validation rules |
| `DecodeStruct(s)` | Marshals a Go struct into a `bson.M` document (used for MongoDB `$set` updates) |
| `EncodeStruct(m, s)` | Unmarshals a `bson.M` document back into a Go struct |

### 7.2 Resource Controllers

Each resource (category, product, question, survey, tag) follows an identical pattern with two files:

- **`<resource>.go`** — handler functions (`create`, `list`, `find`, `delete`, `update`)
- **`<resource>.routes.go`** — route registration via `CreateRoutes()`

#### Standard CRUD Pattern

| Operation | Method | Path | Auth Required | Description |
|-----------|--------|------|:------------:|-------------|
| Create | `POST` | `/<resource>/` | ✅ | Bind + validate → set timestamps → insert |
| List | `GET` | `/<resource>/` | ❌ | Find all (with aggregation for relations) |
| Find | `GET` | `/<resource>/:id` | ❌ | Find by ObjectID (with aggregation for relations) |
| Delete | `DELETE` | `/<resource>/:id` | ✅ | Authenticate → delete by ObjectID |
| Update | `PATCH` | `/<resource>/:id` | ✅ | Authenticate → bind + validate → decode to BSON → update |

#### Auth Controller (`auth/`)

The auth controller is distinct — it handles `login`, `register`, `refresh`, and `logout` rather than standard CRUD. Only `logout` requires authentication via `AuthGuard()`.

#### Dual-Layer Authentication

Protected resource routes use a **two-layer** authentication strategy:

1. **`AuthGuard()` middleware** — registered in the route definition, validates the JWT signature before the handler runs. Returns `401` and aborts the request chain on failure.
2. **`controllers.AuthRequest(c)`** — called inside the handler, extracts the token's `access_uuid` and verifies it still exists in Redis (i.e., the session hasn't been revoked via logout).

This means a valid JWT alone is insufficient — the corresponding Redis session must also be active. This enables server-side session revocation while retaining the stateless benefits of JWT for signature verification.

### 7.3 Aggregation Pipelines

Resources with relationships use MongoDB aggregation `$lookup` to resolve references at query time:

| Resource | Resolved Relations |
|----------|--------------------|
| Product | `tags` → Tags collection, `categories` → Categories collection |
| Question | `tags` → Tags collection |
| Survey | `questions` → Questions collection, `questions.tags` → Tags collection |

---

## 8. Database Module (`database/`)

### 8.1 MongoDB (`mongo.go`)

`InitMongo(url, database)` creates a client, connects with a 10-second timeout, pings to verify, and retrieves the database object. Exits fatally if the connection fails.

### 8.2 Collections (`collections.go`)

Six global `*mongo.Collection` variables are initialized via `SetCollections()`:

| Variable | MongoDB Collection |
|----------|--------------------|
| `UserCollection` | `users` |
| `CategoryCollection` | `categories` |
| `ProductCollection` | `products` |
| `QuestionCollection` | `questions` |
| `SurveyCollection` | `surveys` |
| `TagCollection` | `tags` |

### 8.3 Operators (`operators.go`)

A generic data access layer providing reusable MongoDB operations:

| Function | Operation | Notes |
|----------|-----------|-------|
| `InsertOne(collection, model)` | `insertOne` | Inserts any struct |
| `FindOne(collection, filter)` | `findOne` | Returns `bson.M` |
| `FindAll(collection)` | `find({})` | Returns `[]bson.M` |
| `FindOneWithAggregate(collection, pipeline)` | `aggregate` | Single result from pipeline |
| `FindAllWithAggregate(collection, pipeline)` | `aggregate` | Multiple results from pipeline |
| `UpdateOne(collection, id, update)` | `updateOne` with `$set` | By `_id` filter |
| `DeleteOne(collection, id)` | `deleteOne` | By `_id` filter |

### 8.4 Redis Cache (`cache.go`)

Manages JWT token sessions in Redis:

| Function | Purpose |
|----------|---------|
| `InitRedis(host, port, password)` | Creates and verifies the Redis connection |
| `CreateRedisAccess(id, token)` | Stores access and refresh UUIDs as keys, user ID as value, with TTL matching token expiration |
| `GetRedisAccess(access)` | Checks if an access UUID exists in Redis (validates the session) |
| `DeleteRedisAccess(uuid)` | Removes a UUID key from Redis (logout/refresh invalidation) |

---

## 9. Models Module (`models/`)

All models use MongoDB's `primitive.ObjectID` as the primary key and carry `bson`/`json`/`validate` struct tags.

### 9.1 Domain Models

| Model | Fields | Validation | Relationships |
|-------|--------|------------|---------------|
| **User** | `email`, `password`, `created_at`, `updated_at` | `email`: required + email format; `password`: required, 8-64 chars | — |
| **Category** | `name`, `created_at`, `updated_at` | `name`: required | Referenced by Product |
| **Tag** | `name`, `created_at`, `updated_at` | `name`: required | Referenced by Product, Question |
| **Product** | `name`, `brand`, `categories`, `summary`, `description`, `image`, `tags`, `active`, `created_at`, `updated_at` | `name`: required; `brand`: required; `categories`: required | References Category[], Tag[] |
| **Question** | `description`, `tags`, `active`, `created_at`, `updated_at` | `description`: required | References Tag[] |
| **Survey** | `name`, `questions`, `active`, `created_at`, `updated_at` | `name`: required | References Question[] |

> **Note:** Product, Question, and Survey include an `active` boolean field (defaults to `true` on creation), enabling soft-disable semantics. Category and Tag omit this field as they serve as lightweight reference entities.

### 9.2 Auth Models

| Model | Fields | Purpose |
|-------|--------|---------|
| **Token** | `access_token`, `access_uuid`, `access_expires`, `refresh_token`, `refresh_uuid`, `refresh_expires` | Carries JWT token pair metadata |
| **RedisAccess** | `access_uuid`, `user_id` | Maps a Redis-cached session to its owner |

### 9.3 Entity Relationship Diagram

Relationships are stored as **embedded arrays of ObjectIDs** within the referencing document (not as separate join collections). They are resolved at query time via MongoDB `$lookup` aggregation pipelines.

```
┌──────────┐            ┌──────────────┐            ┌──────────┐
│   Tag    │◄───────────│   Product    │───────────►│ Category │
│          │ ObjectID[] │              │ ObjectID[] │          │
└──────────┘            └──────────────┘            └──────────┘
      ▲
      │ ObjectID[]
┌──────────┐            ┌──────────────┐
│ Question │◄───────────│    Survey    │
│          │ ObjectID[] │              │
└──────────┘            └──────────────┘

┌──────────┐            ┌──────────────┐
│   User   │ · · · · · ·│  Token/Redis │
│          │  1:N (TTL) │  (sessions)  │
└──────────┘            └──────────────┘
```

---

## 10. Utils Module (`utils/`)

### 10.1 Constants (`constants.go`)

Centralized string constants for all API response messages (auth errors, CRUD success/failure messages, route errors). This ensures consistent messaging across all controllers.

### 10.2 Logger (`logger.go`)

Two logging utilities:

| Function | Behavior |
|----------|----------|
| `LogError(errs...)` | Logs and returns the first non-nil error (variadic) |
| `LogFatalError(errs...)` | Calls `log.Fatal()` on the first non-nil error, terminating the process |

### 10.3 Responses (`responses.go`)

Standardized JSON response helpers:

| Function | Output Format |
|----------|--------------|
| `ConstantResponse(c, status, message)` | `{"response": "<message>"}` |
| `DynamicResponse(c, status, map)` | Any arbitrary JSON map |

**Response examples:**

Success message (e.g., after create, update, delete):
```json
{ "response": "Document succesfully inserted." }
```

Error message:
```json
{ "response": "Unauthorized access." }
```

Login success (token pair):
```json
{ "access_token": "eyJhbG...", "refresh_token": "eyJhbG..." }
```

Resource list:
```json
{ "response": [ { "_id": "5f51...", "name": "Category", "created_at": "..." } ] }
```

Single resource (with resolved relations):
```json
{ "response": { "_id": "5f51...", "name": "Product", "tags": [...], "categories": [...] } }
```

### 10.4 Validators (`validators/model.go`)

| Function | Purpose |
|----------|---------|
| `BindModel(c, model)` | Binds the Gin request body to a model struct via `ShouldBind` |
| `ValidateModel(model)` | Runs `go-playground/validator` rules defined in struct tags |

---

## 11. API Endpoints Reference

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|:----:|-------------|
| `POST` | `/v1/login` | ❌ | Authenticate with email/password, receive token pair |
| `POST` | `/v1/register` | ❌ | Create a new user account |
| `POST` | `/v1/refresh` | ❌ | Exchange a refresh token for a new token pair |
| `POST` | `/v1/logout` | ✅ | Invalidate the current access token |

### Resources (Category, Product, Question, Survey, Tag)

| Method | Endpoint | Auth | Description |
|--------|----------|:----:|-------------|
| `POST` | `/v1/<resource>/` | ✅ | Create a new resource |
| `GET` | `/v1/<resource>/` | ❌ | List all resources |
| `GET` | `/v1/<resource>/:id` | ❌ | Find a single resource by ID |
| `PATCH` | `/v1/<resource>/:id` | ✅ | Update a resource by ID |
| `DELETE` | `/v1/<resource>/:id` | ✅ | Delete a resource by ID |

### Postman Collection

A ready-to-use Postman collection is provided in `requests.json`. It includes all endpoints with sample payloads and a **pre-request script** that automatically authenticates before each request by calling `/v1/login` and setting the `access_token` environment variable.

---

## 12. Key Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `gin-gonic/gin` | v1.9.1 | HTTP web framework and router |
| `go.mongodb.org/mongo-driver` | v1.5.1 | Official MongoDB Go driver |
| `go-redis/redis` | v7.2.0 | Redis client for Go |
| `dgrijalva/jwt-go` | v3.2.0 | JWT creation and validation |
| `go-playground/validator` | v10.14.0 | Struct validation via tags |
| `twinj/uuid` | v1.0.0 | UUID generation for token identifiers |
| `golang.org/x/crypto` | v0.9.0 | Bcrypt password hashing |
| `cespare/reflex` | (dev tool) | File watcher for hot-reloading in development |

---

## 13. Project Structure

```
go_api_boilerplate/
├── docker-compose.yml          # Development stack (API + Mongo + Redis)
├── docker-compose.prod.yml     # Production stack
├── Dockerfile                  # Development image (with Reflex hot-reload)
├── Dockerfile.prod             # Production multi-stage image (scratch)
├── .env.example                # Environment variable template
├── requests.json               # Postman collection for API testing
├── LICENSE                     # MIT License
├── README.md                   # Project documentation
├── ARCHITECTURE.md             # This file
└── src/
    ├── api.go                  # Application entry point
    ├── go.mod                  # Go module definition
    ├── go.sum                  # Dependency checksums
    ├── reflex.conf             # Hot-reload configuration
    ├── controllers/
    │   ├── common.go           # Shared controller utilities
    │   ├── auth/
    │   │   ├── auth.go         # Login, register, refresh, logout handlers
    │   │   └── auth.routes.go  # Auth route registration
    │   ├── category/
    │   │   ├── category.go         # Category CRUD handlers
    │   │   └── category.routes.go  # Category route registration
    │   ├── product/
    │   │   ├── product.go          # Product CRUD handlers
    │   │   └── product.routes.go   # Product route registration
    │   ├── question/
    │   │   ├── question.go         # Question CRUD handlers
    │   │   └── question.routes.go  # Question route registration
    │   ├── survey/
    │   │   ├── survey.go           # Survey CRUD handlers
    │   │   └── survey.routes.go    # Survey route registration
    │   └── tag/
    │       ├── tag.go              # Tag CRUD handlers
    │       └── tag.routes.go       # Tag route registration
    ├── database/
    │   ├── cache.go            # Redis initialization and token caching
    │   ├── collections.go      # MongoDB collection references
    │   ├── mongo.go            # MongoDB initialization
    │   └── operators.go        # Generic MongoDB CRUD operations
    ├── middleware/
    │   ├── authentication.go   # JWT creation, verification, and AuthGuard
    │   └── refresh.go          # Refresh token handling
    ├── models/
    │   ├── auth.go             # Token and RedisAccess models
    │   ├── user.go             # User model
    │   ├── category.go         # Category model
    │   ├── product.go          # Product model
    │   ├── question.go         # Question model
    │   ├── survey.go           # Survey model
    │   └── tag.go              # Tag model
    └── utils/
        ├── constants.go        # Response message constants
        ├── logger.go           # Error logging utilities
        ├── responses.go        # Standardized JSON response helpers
        └── validators/
            └── model.go        # Request binding and validation
```

---

## 14. Design Decisions

| Decision | Rationale |
|----------|-----------|
| **Gin framework** | Lightweight, high-performance HTTP framework with built-in middleware support, JSON binding, and route grouping |
| **MongoDB** | Schema-flexible document database well-suited for the varied resource models (products with nested references, surveys with questions) |
| **Redis for JWT sessions** | Enables token revocation (logout) without token blacklisting; TTL-based expiration aligns naturally with JWT expiry |
| **Dual-token JWT (access + refresh)** | Short-lived access tokens (15 min) minimize exposure; long-lived refresh tokens (7 days) reduce re-authentication friction |
| **Bcrypt password hashing** | Industry-standard adaptive hashing algorithm resistant to brute-force attacks |
| **Global collection variables** | Simple singleton pattern avoids dependency injection complexity in a boilerplate; collections initialized once at startup |
| **Generic database operators** | Reusable `InsertOne`, `FindAll`, `DeleteOne`, etc. avoid code duplication across controllers |
| **MongoDB aggregation for relations** | `$lookup` pipelines resolve ObjectID references at query time, avoiding manual join logic in application code |
| **Versioned API routes (`/v1/`)** | Enables future API evolution without breaking existing clients |
| **Multi-stage Docker build** | Production image uses `scratch` base for minimal size and attack surface (~5 MB vs ~300 MB) |
| **Reflex hot-reloading** | Watches `.go` and `go.mod` files, re-runs `go run .` on changes — eliminates manual restart during development |
| **Environment-driven config** | All secrets and connection strings injected via env vars, following the Twelve-Factor App methodology |
| **Controller-per-resource packages** | Each resource is a self-contained Go package with its own handlers and route registration, promoting separation of concerns |

---

## 15. Security Considerations

The boilerplate implements several security measures out of the box:

| Aspect | Implementation |
|--------|----------------|
| **Password storage** | Passwords are hashed with bcrypt (`DefaultCost`) before storage — plaintext passwords are never persisted |
| **Token signing** | JWTs are signed with HMAC-SHA256 using separate secrets for access and refresh tokens |
| **Session revocation** | Deleting a token's UUID from Redis immediately invalidates it, even if the JWT hasn't expired |
| **Short-lived access tokens** | 15-minute TTL limits the window of exposure if a token is compromised |
| **Panic recovery** | `gin.Default()` includes recovery middleware that catches panics and returns `500` instead of crashing |
| **Minimal production image** | `scratch`-based Docker image contains only the static binary — no shell, no package manager, reduced attack surface |
| **Fail-fast startup** | Application exits immediately (`log.Fatal`) if MongoDB or Redis is unreachable, preventing operation in a degraded state |

---

## 16. Extending the Boilerplate

### Adding a New Resource

To add a new CRUD resource (e.g., `order`), follow these steps:

1. **Define the model** — create `src/models/order.go` with a struct using `bson`/`json`/`validate` tags:

    ```go
    type Order struct {
        ID        primitive.ObjectID `bson:"_id,omitempty"`
        Name      string             `bson:"name" json:"name" validate:"required"`
        Active    bool               `bson:"active,omitempty"`
        CreatedAt time.Time          `bson:"created_at,omitempty"`
        UpdatedAt time.Time          `bson:"updated_at,omitempty"`
    }
    ```

2. **Register the collection** — add a global variable in `src/database/collections.go` and initialize it in `SetCollections()`:

    ```go
    var OrderCollection *mongo.Collection
    // inside SetCollections():
    OrderCollection = c.Collection("orders")
    ```

3. **Create the controller** — add `src/controllers/order/order.go` with `create`, `list`, `find`, `delete`, `update` handlers following the existing pattern (e.g., `category.go`).

4. **Register the routes** — add `src/controllers/order/order.routes.go` with a `CreateRoutes()` function:

    ```go
    func CreateRoutes(r *gin.RouterGroup) {
        order := r.Group("/order")
        {
            order.POST("/", middleware.AuthGuard(), create)
            order.GET("/", list)
            order.GET("/:id", find)
            order.DELETE("/:id", middleware.AuthGuard(), delete)
            order.PATCH("/:id", middleware.AuthGuard(), update)
        }
    }
    ```

5. **Wire the router** — import and call `order.CreateRoutes(v1)` in `src/server/router.go`.

### Adding a New Aggregation Relationship

To resolve ObjectID references in queries (e.g., an `order` referencing `products`), add a `$lookup` stage to the aggregation pipeline in your controller's `list` and `find` handlers:

```go
pipeline := []bson.M{
    {"$lookup": bson.M{"from": "products", "localField": "products", "foreignField": "_id", "as": "products"}},
}
```

# Snippetbox 📦

A robust, idiomatic web application built in Go for creating, storing, and sharing text and code snippets with automatic expiration dates (inspired by Alex Edwards' *"Let's Go"*).

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## ✨ Features

- **Snippet Management**: Create, read, and list text and code snippets.
- **Auto-Expiration**: Configure expiration windows (1 Day, 1 Week, 1 Year). Expired snippets are automatically filtered out.
- **Form Validation**: Strict validation for inputs (title length, content presence, valid expiration periods) with inline error highlights.
- **PostgreSQL Persistence**: Backed by PostgreSQL with connection pooling (`database/sql` + `lib/pq`).
- **Zero-Config Local Demo**: Includes an in-memory fallback store if a database connection is not detected, enabling immediate testing and development out of the box.
- **Security Middleware**:
  - Secure HTTP headers (CSP, HSTS/Nosniff, X-Frame-Options, XSS Protection).
  - Panic recovery with clean 500 error responses and connection termination.
  - Structured request logging (Client IP, Protocol, Method, Request URI).
- **Comprehensive Test Suite**: End-to-end integration and unit tests for handlers, models, validator, and middleware.

---

## 🛠 Tech Stack

- **Language**: Go (1.22+)
- **Router**: [`httprouter`](https://github.com/julienschmidt/httprouter)
- **Middleware Chaining**: [`alice`](https://github.com/justinas/alice)
- **Database Driver**: [`lib/pq`](https://github.com/lib/pq) (PostgreSQL)
- **Environment Management**: [`godotenv`](https://github.com/joho/godotenv)
- **Templating**: Go Standard Library `html/template` with template caching and custom template funcs

---

## 📁 Project Structure

```text
snippetbox/
├── cmd/
│   └── web/                   # Application entrypoint, routing, handlers, middleware
│       ├── main.go            # Entrypoint, server initialization, DB connection
│       ├── handlers.go        # HTTP request handlers (home, view, create)
│       ├── routes.go          # Route definitions and middleware chaining
│       ├── middleware.go      # Secure headers, panic recovery, request logger
│       ├── templates.go       # HTML template cache & custom template functions
│       ├── helpers.go         # Common HTTP helpers (render, errors, paths)
│       ├── handlers_test.go   # Integration tests for route handlers
│       ├── middleware_test.go # Unit tests for HTTP middlewares
│       └── main_test.go       # Server configuration and routing tests
├── internal/
│   ├── models/                # Data access layer
│   │   ├── snippets.go        # SnippetModel (PostgreSQL queries & in-memory fallback)
│   │   ├── errors.go          # Domain errors (ErrNoRecord)
│   │   └── snippets_test.go   # Model CRUD and expiration tests
│   └── validator/             # Form validation library
│       ├── validator.go       # Validation rules and helper functions
│       └── validator_test.go  # Validator unit tests
├── ui/                        # Presentation layer
│   ├── html/                  # HTML templates
│   │   ├── base.tmpl          # Master layout
│   │   ├── partials/          # Reusable components (e.g. nav)
│   │   └── pages/             # View, Home, and Create pages
│   └── static/                # Static assets (CSS, JavaScript, images, favicon)
├── schema.sql                 # PostgreSQL table definition, index, and seed data
├── .env.example               # Example environment variable configuration
├── go.mod                     # Go module definition
└── go.sum                     # Checksums for dependencies
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or higher installed.
- (Optional) [PostgreSQL](https://www.postgresql.org/) or [Docker](https://www.docker.com/) for database persistence.

---

### Option A: Quick Run (In-Memory Demo, No DB Required)

If you don't have PostgreSQL running, Snippetbox will automatically start in in-memory mode with sample seed data:

```bash
# Clone the repository
git clone https://github.com/Mkz-Prog/snippetbox.git
cd snippetbox

# Run the web server
go run ./cmd/web
```

Open your browser at **[http://localhost:4000](http://localhost:4000)**.

---

### Option B: Full Setup with PostgreSQL

#### 1. Start PostgreSQL

Using **Docker**:
```bash
docker run --name snippetbox-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=snippetbox \
  -p 5432:5432 \
  -d postgres:17
```

Or using an existing local PostgreSQL installation:
```bash
sudo -u postgres psql -c "CREATE DATABASE snippetbox;"
```

#### 2. Apply Database Schema & Seed Data

Run `schema.sql` against your PostgreSQL database:

```bash
psql -h localhost -U postgres -d snippetbox -f schema.sql
```

#### 3. Configure Environment

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` if your credentials differ from the defaults:
```env
ADDR=:4000
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=snippetbox
```

#### 4. Run the Application

```bash
go run ./cmd/web
```

The application will connect to PostgreSQL and listen on `http://localhost:4000`.

---

## ⚙️ Configuration & Flags

Snippetbox can be configured via environment variables or command-line flags. Command-line flags take precedence:

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:4000` | HTTP network address to listen on |
| `-dsn` | `postgres://postgres:postgres@localhost:5432/snippetbox?sslmode=disable` | PostgreSQL connection string |

Example using flags:
```bash
go run ./cmd/web -addr=":8080" -dsn="postgres://user:secret@localhost:5432/snippetbox?sslmode=disable"
```

---

## 🌐 Routes & Endpoints

| Method | Pattern | Handler | Description |
|--------|---------|---------|-------------|
| `GET` | `/` | `app.home` | Displays the 10 most recent unexpired snippets |
| `GET` | `/snippet/view/:id` | `app.snippetView` | Shows a specific snippet by ID |
| `GET` | `/snippet/create` | `app.snippetCreate` | Displays the snippet creation form |
| `POST` | `/snippet/create` | `app.snippetCreatePost` | Validates and creates a new snippet |
| `GET` | `/static/*filepath` | Static File Server | Serves CSS, JS, images, and favicon |

---

## 🧪 Testing & Verification

Run the complete automated test suite:

```bash
# Run all tests with verbose output
go test -v ./...

# Run tests with code coverage summary
go test -cover ./...

# Format all Go source files
go fmt ./...
```

---

## 📦 Building for Production

Compile a standalone binary:

```bash
# Build the binary
go build -o snippetbox ./cmd/web

# Run the compiled binary
./snippetbox -addr=":80"
```

---

## 📄 License

This project is open-source software licensed under the [MIT License](LICENSE).

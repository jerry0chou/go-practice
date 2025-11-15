# Go Practice

A comprehensive collection of Go programming examples and practice code covering various topics and common patterns.

## Overview

This repository contains practical examples and implementations for different aspects of Go development, organized by topic. Each module includes working code examples that demonstrate best practices and common patterns.

## Modules

- **Concurrency**: Goroutines, channels, mutexes, worker pools, context, select statements, and fan patterns
- **Config**: Environment variables, file-based configuration (JSON, YAML, TOML), hot reload, and validation
- **Database**: SQL basics, ORM (GORM), connection pooling, migrations, transactions, and multi-driver support (MySQL, PostgreSQL, SQLite)
- **HTTP**: Client/server implementations, middleware, GitHub API client, and utilities
- **Security**: JWT authentication, OAuth, RBAC authorization, password hashing, HTTPS/TLS, and input validation
- **Networking**: TCP/UDP examples, network utilities, and URL operations
- **Reflection**: Basic reflection, struct/interface/function reflection, and practical examples
- **File Operations**: File I/O operations and utilities
- **JSON**: JSON encoding/decoding and operations
- **String Operations**: String manipulation utilities
- **Format**: Formatting examples
- **Server**: HTTP server with handlers, middleware, and routing

## Getting Started

### Prerequisites

- Go 1.24.4 or later

### Installation

```bash
git clone https://github.com/jerrychou/go-practice.git
cd go-practice
go mod download
```

### Running Examples

Each module has a corresponding main file in the `run/` directory. For example:

```bash
# Run concurrency examples
go run run/concurrency_main.go

# Run security examples
go run run/security_main.go

# Run HTTP examples
go run run/http_main.go
```

## Project Structure

```
go-practice/
├── concurrency/     # Concurrency patterns and examples
├── config/          # Configuration management
├── database/        # Database operations and ORM
├── http/            # HTTP client and server
├── security/        # Security implementations
├── net/             # Network programming
├── reflect/         # Reflection examples
├── run/             # Main entry points for each module
└── ...
```

## Dependencies

Key dependencies include:
- GORM for database ORM
- JWT for authentication
- Various database drivers (MySQL, PostgreSQL, SQLite)

See `go.mod` for the complete list of dependencies.

## License

This is a practice repository for learning and reference purposes.


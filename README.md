# MCP Postgres

This project is a Go-based MCP (Machine Context Protocol) server that connects to a PostgreSQL database.

## Features
- Exposes endpoints for interacting with a PostgreSQL database
- Supports SELECT queries (default)

## Getting Started

### Prerequisites
- Go 1.18+
- PostgreSQL

### Setup
1. Clone this repository.
2. Configure your database connection in the environment or config file.
3. Run the server:
   ```sh
   go run main.go
   ```

### Example Usage
- List users in the database
- Insert, update, or delete users (if enabled in the server)

## Security
- By default, only SELECT queries are allowed for safety.
- To enable write operations, update the server code accordingly.

## License
MIT

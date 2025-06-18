package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" // Make sure to import PostgreSQL driver
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var db *sql.DB // Global DB variable

func main() {
	connStr := "user=root password=password dbname=mydb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Force a connection to verify credentials
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	tool := mcp.NewTool("get_table_data",
		mcp.WithDescription("Get the table contents from the db"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Postgres SQL query to select data")),
	)

	s.AddTool(tool, queryHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func queryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	queryTrimmed := strings.TrimSpace(strings.ToLower(query))
	if !strings.HasPrefix(queryTrimmed, "select") {
		return mcp.NewToolResultError("Only SELECT queries are allowed"), nil
	}

	rows, dbErr := db.Query(query)
	if dbErr != nil {
		return mcp.NewToolResultError(dbErr.Error()), nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return mcp.NewToolResultError("Failed to get columns"), nil
	}

	var results []map[string]any

	for rows.Next() {
		// Create a slice of interfaces to hold column values
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to scan row: %v", err)), nil
		}

		// Map the values to column names
		rowMap := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			// Handle []byte conversion
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Row error: %v", err)), nil
	}

	// Marshal results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError("Failed to marshal result to JSON"), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

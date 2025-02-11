package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

// ConnectDB initializes and returns a database connection
func ConnectDB() (*sql.DB, error) {
	// SQL Server connection details
	server := "localhost"              // Change this if your SQL Server is on another machine
	port := 1433                       // Default SQL Server port
	user := "sa"                       // Change to your actual SQL Server username
	password := "Yonatan1234$"         // Change to your actual password
	database := "rise_home_assignment" // Database name

	// Connection string
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, database)

	// Open a database connection
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to database: %v", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to SQL Server successfully!")
	return db, nil
}

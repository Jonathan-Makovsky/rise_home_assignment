package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

// ConnectDB initializes and returns a database connection
func ConnectDB() (*sql.DB, error) {
	// SQL Server connection details
	server := "db"                     // Docker service name
	port := 1433                       // Default SQL Server port
	user := "sa"                       // SQL Server username
	password := "Yonatan1234$"         // SQL Server password
	database := "rise_home_assignment" // Database name

	// Connection string
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, database)

	// Retry connection to the database (with backoff)
	var db *sql.DB
	var err error
	for retries := 0; retries < 10; retries++ {
		db, err = sql.Open("sqlserver", dsn)
		if err != nil {
			log.Printf("❌ Failed to connect to database, retrying... (%d/10)", retries+1)
			time.Sleep(2 * time.Second)
			continue
		}

		// Check if the connection is successful
		err = db.Ping()
		if err != nil {
			log.Printf("❌ Failed to ping database, retrying... (%d/10)", retries+1)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("✅ Connected to SQL Server successfully!")
		return db, nil
	}

	return nil, fmt.Errorf("❌ Failed to connect to database after 10 retries: %v", err)
}

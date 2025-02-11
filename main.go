package main

import (
	"fmt"
	"log"
	"rise_home_assignment/internal/config"
	"rise_home_assignment/internal/handlers"
	"rise_home_assignment/internal/repository"

	"github.com/gin-contrib/cors" // Import CORS package
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	// Initialize repository and handler
	contactRepo := repository.NewContactRepository(db)
	contactHandler := handlers.NewContactHandler(contactRepo)

	// Create a Gin router
	router := gin.Default()

	// Enable CORS with default settings (allow all origins)
	router.Use(cors.Default())

	// Register API routes
	router.GET("/contacts", contactHandler.GetContacts)           // Get contacts with pagination
	router.GET("/contacts/search", contactHandler.SearchContacts) // Search contacts by query
	router.POST("/contacts", contactHandler.AddContact)           // Add a new contact
	router.PUT("/contacts", contactHandler.EditContact)           // Edit an existing contact
	router.DELETE("/contacts", contactHandler.DeleteContact)      // Delete a contact

	// Start the server
	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}

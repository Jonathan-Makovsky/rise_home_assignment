package handlers

import (
	"net/http"
	"rise_home_assignment/internal/models"
	"rise_home_assignment/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ContactHandler handles API requests for contacts
type ContactHandler struct {
	Repo *repository.ContactRepository
}

// NewContactHandler creates a new instance of ContactHandler
func NewContactHandler(repo *repository.ContactRepository) *ContactHandler {
	return &ContactHandler{Repo: repo}
}

// GetContacts retrieves a paginated list of contacts
func (h *ContactHandler) GetContacts(c *gin.Context) {
	// Get the page number from the query parameters (default to page 1)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	// Fetch contacts from the database
	contacts, err := h.Repo.GetContacts(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	// Return contacts as JSON response
	c.JSON(http.StatusOK, contacts)
}

// SearchContacts searches for contacts by first name, last name, or phone number
func (h *ContactHandler) SearchContacts(c *gin.Context) {
	// Get the search query from URL parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Fetch search results from the database
	contacts, err := h.Repo.SearchContacts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search contacts"})
		return
	}

	// Return the search results as JSON
	c.JSON(http.StatusOK, contacts)
}

// AddContact adds a new contact to the database
func (h *ContactHandler) AddContact(c *gin.Context) {
	var contact models.Contact

	// Bind JSON request to the contact model
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Insert contact into the database
	err := h.Repo.AddContact(contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add contact"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Contact added successfully"})
}

// EditContact modifies an existing contact in the database
func (h *ContactHandler) EditContact(c *gin.Context) {
	var contact models.Contact

	// Bind JSON request to the contact model
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Edit the contact in the database
	err := h.Repo.EditContact(contact)
	if err != nil {
		// Return the error message if the contact doesn't exist or update fails
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Contact edited successfully"})
}

// DeleteContact removes a contact by phone number
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	// Get the phone number from URL parameters
	phoneNumber := c.DefaultQuery("phone", "")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	// Delete the contact from the database
	err := h.Repo.DeleteContact(phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

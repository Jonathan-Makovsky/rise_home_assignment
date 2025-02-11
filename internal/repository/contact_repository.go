package repository

import (
	"database/sql"
	"fmt"
	"rise_home_assignment/internal/models"

	_ "github.com/microsoft/go-mssqldb" // MSSQL driver
)

// ContactRepository handles database operations
type ContactRepository struct {
	DB *sql.DB
}

// NewContactRepository creates a new instance of ContactRepository
func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

// GetContacts retrieves a maximum of 10 contacts with pagination
func (r *ContactRepository) GetContacts(page int) ([]models.Contact, error) {
	limit := 10
	offset := (page - 1) * limit
	query := "SELECT first_name, last_name, phone_number, address, created_at FROM contacts ORDER BY first_name OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY;"

	rows, err := r.DB.Query(query,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contacts: %v", err)
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(&contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address, &contact.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %v", err)
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// SearchContacts searches for contacts by first name, last name, phone number, or address
func (r *ContactRepository) SearchContacts(query string) ([]models.Contact, error) {
	// Create a query to search by first name, last name, or address
	// We will use `LIKE` to match substrings for first name, last name, phone number, and address
	query = "%" + query + "%"

	sqlQuery := `
		SELECT first_name, last_name, phone_number, address, created_at
		FROM contacts
		WHERE first_name LIKE @query OR last_name LIKE @query OR phone_number LIKE @query OR address LIKE @query`

	rows, err := r.DB.Query(sqlQuery, sql.Named("query", query))
	if err != nil {
		return nil, fmt.Errorf("failed to search contacts: %v", err)
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(&contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address, &contact.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %v", err)
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// AddContact inserts a new contact into the database
func (r *ContactRepository) AddContact(contact models.Contact) error {
	query := "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES (@first_name, @last_name, @phone_number, @address)"
	_, err := r.DB.Exec(query,
		sql.Named("first_name", contact.FirstName),
		sql.Named("last_name", contact.LastName),
		sql.Named("phone_number", contact.PhoneNumber),
		sql.Named("address", contact.Address),
	)
	if err != nil {
		return fmt.Errorf("failed to insert contact: %v", err)
	}
	return nil
}

func (r *ContactRepository) EditContact(contact models.Contact) error {
	// First, check if the contact exists based on phone number
	var exists bool
	err := r.DB.QueryRow("SELECT 1 FROM contacts WHERE phone_number = @phone_number",
		sql.Named("phone_number", contact.PhoneNumber)).Scan(&exists)

	if err == sql.ErrNoRows {
		// If no contact is found, return a formatted error message
		return fmt.Errorf("contact with phone number %s not found", contact.PhoneNumber)
	} else if err != nil {
		// If any other error occurs while querying, return it
		return fmt.Errorf("failed to check contact existence: %v", err)
	}

	// Now update the contact since it exists
	query := "UPDATE contacts SET first_name = @first_name, last_name = @last_name, phone_number = @phone_number, address = @address WHERE phone_number = @phone_number"
	_, err = r.DB.Exec(query,
		sql.Named("first_name", contact.FirstName),
		sql.Named("last_name", contact.LastName),
		sql.Named("phone_number", contact.PhoneNumber),
		sql.Named("address", contact.Address),
	)

	if err != nil {
		return fmt.Errorf("failed to edit contact: %v", err)
	}
	return nil
}

// DeleteContact removes a contact by phone number
func (r *ContactRepository) DeleteContact(phoneNumber string) error {
	// First, check if the contact exists based on phone number
	var exists bool
	err := r.DB.QueryRow("SELECT 1 FROM contacts WHERE phone_number = @phone_number",
		sql.Named("phone_number", phoneNumber)).Scan(&exists)

	if err == sql.ErrNoRows {
		// If no contact is found, return an error
		return fmt.Errorf("contact with phone number %s not found", phoneNumber)
	} else if err != nil {
		// If any other error occurs while querying, return it
		return fmt.Errorf("failed to check contact existence: %v", err)
	}

	// Now delete the contact since it exists
	query := "DELETE FROM contacts WHERE phone_number = @phone_number"
	_, err = r.DB.Exec(query, sql.Named("phone_number", phoneNumber))
	if err != nil {
		return fmt.Errorf("failed to delete contact: %v", err)
	}
	return nil
}

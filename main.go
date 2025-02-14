// main.go
package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Contact struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Phone   string `json:"phone"`
    Email   string `json:"email"`
}

var db *sql.DB

func main() {
    // Database connection string from environment variables
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:postgres@db:5432/phonebook?sslmode=disable"
    }

    // Connect to database
    var err error
    db, err = sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create router
    r := mux.NewRouter()

    // Define routes
    r.HandleFunc("/contacts", getContacts).Methods("GET")
    r.HandleFunc("/contacts", createContact).Methods("POST")
    r.HandleFunc("/contacts/{id}", getContact).Methods("GET")
    r.HandleFunc("/contacts/{id}", updateContact).Methods("PUT")
    r.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")

    // Start server
    log.Printf("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func getContacts(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, phone, email FROM contacts")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var contacts []Contact
    for rows.Next() {
        var c Contact
        if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        contacts = append(contacts, c)
    }

    json.NewEncoder(w).Encode(contacts)
}

func createContact(w http.ResponseWriter, r *http.Request) {
    var contact Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := db.QueryRow(
        "INSERT INTO contacts (name, phone, email) VALUES ($1, $2, $3) RETURNING id",
        contact.Name, contact.Phone, contact.Email,
    ).Scan(&contact.ID)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(contact)
}

func getContact(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var contact Contact

    err := db.QueryRow(
        "SELECT id, name, phone, email FROM contacts WHERE id = $1",
        vars["id"],
    ).Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.Email)

    if err == sql.ErrNoRows {
        http.Error(w, "Contact not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(contact)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var contact Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.Exec(
        "UPDATE contacts SET name = $1, phone = $2, email = $3 WHERE id = $4",
        contact.Name, contact.Phone, contact.Email, vars["id"],
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if rowsAffected == 0 {
        http.Error(w, "Contact not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    result, err := db.Exec("DELETE FROM contacts WHERE id = $1", vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if rowsAffected == 0 {
        http.Error(w, "Contact not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

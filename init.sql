CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100)
);

-- Add some sample data
INSERT INTO contacts (name, phone, email) VALUES
    ('John Doe', '+1-555-555-5555', 'john@example.com'),
    ('Jane Smith', '+1-555-555-5556', 'jane@example.com');

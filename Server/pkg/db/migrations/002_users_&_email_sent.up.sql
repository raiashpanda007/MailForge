CREATE TABLE clients (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    api_key_id UUID NOT NULL REFERENCES apikeys(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE emails_sent (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id),
    body VARCHAR(2000) NOT NULL,
    subject VARCHAR(300) NOT NULL,
    client_email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

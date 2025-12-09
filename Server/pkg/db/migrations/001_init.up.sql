CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE apikeys (
    id UUID PRIMARY KEY,
    organization VARCHAR(200) UNIQUE,
    apikey TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    email_app_password TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


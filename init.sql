-- Initialize database schema for GalleryBlue

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for email lookups during login
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index for display_name uniqueness checks
CREATE INDEX IF NOT EXISTS idx_users_display_name ON users(display_name);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,

    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    alt_email VARCHAR(255),
    phone VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,

    password_hash TEXT,

    profile_pic TEXT,
    status BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    last_login_at TIMESTAMPTZ,
    date_of_birth DATE,

    address_line_1 TEXT,
    address_line_2 TEXT,
    city VARCHAR(255),
    state VARCHAR(255),
    pin_code VARCHAR(20),

    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_secret TEXT,
    two_factor_pending_secret TEXT,

    backup_code_hashes JSONB DEFAULT '[]'::jsonb
);

-- Unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username
    ON users(username);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
    ON users(email);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone
    ON users(phone);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_alt_email
    ON users(alt_email);
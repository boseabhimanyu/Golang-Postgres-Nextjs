
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,

    first_name VARCHAR(255),
    last_name VARCHAR(255),
    username VARCHAR(255),
    email VARCHAR(255),
    alt_email VARCHAR(255),
    phone VARCHAR(50),

    role VARCHAR(50),
    profile_pic TEXT,
    status BOOLEAN NOT NULL DEFAULT TRUE,

    password_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    last_login_at TIMESTAMPTZ,
    date_of_birth DATE,

    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(255),
    state VARCHAR(255),
    pin_code VARCHAR(20),

    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_secret TEXT,
    two_factor_pending_secret TEXT,

    backup_code_hashes JSONB DEFAULT '[]'::jsonb
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

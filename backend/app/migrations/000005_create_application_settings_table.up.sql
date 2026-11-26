
CREATE TABLE IF NOT EXISTS application_settings (
    id VARCHAR(50) PRIMARY KEY,

    registration_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,

    login_with_primary_email BOOLEAN NOT NULL DEFAULT TRUE,
    login_with_username BOOLEAN NOT NULL DEFAULT TRUE,
    login_with_phone BOOLEAN NOT NULL DEFAULT FALSE,
    login_with_alt_email BOOLEAN NOT NULL DEFAULT FALSE,

    menu_max_depth JSONB,

    updated_at TIMESTAMPTZ NOT NULL,
    updated_by VARCHAR(255)
);

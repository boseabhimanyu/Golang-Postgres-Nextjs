
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY,

    name VARCHAR(255),
    location VARCHAR(50),

    status BOOLEAN NOT NULL DEFAULT TRUE,

    items JSONB,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

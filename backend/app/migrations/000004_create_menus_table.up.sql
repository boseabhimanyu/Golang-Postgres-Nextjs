
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    location VARCHAR(50) NOT NULL
        CHECK (location IN ('top', 'left', 'bottom')),

    status BOOLEAN NOT NULL DEFAULT TRUE,

    items JSONB,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT uq_menus_location UNIQUE (location)
);
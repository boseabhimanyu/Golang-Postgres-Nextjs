
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY,

    title VARCHAR(255),
    content TEXT,

    author_id UUID,
    slug VARCHAR(255),
    visibility VARCHAR(50),

    is_homepage BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    css TEXT
);

-- GORM-defined indexes
CREATE INDEX IF NOT EXISTS idx_pages_author_id
    ON pages(author_id);

CREATE INDEX IF NOT EXISTS idx_pages_slug
    ON pages(slug);

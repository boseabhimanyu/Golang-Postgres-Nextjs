
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY,

    title VARCHAR(255) NOT NULL,
    content TEXT,

    author_id UUID NOT NULL,
    slug VARCHAR(255) NOT NULL,
    visibility VARCHAR(50) NOT NULL
        CHECK (visibility IN ('public', 'registered')),

    is_homepage BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    css TEXT,

    CONSTRAINT fk_pages_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT uq_pages_slug UNIQUE (slug)
);

-- Author lookup index
CREATE INDEX IF NOT EXISTS idx_pages_author_id
    ON pages(author_id);

-- Only one page can be the homepage
CREATE UNIQUE INDEX IF NOT EXISTS idx_pages_single_homepage
    ON pages(is_homepage)
    WHERE is_homepage = TRUE;
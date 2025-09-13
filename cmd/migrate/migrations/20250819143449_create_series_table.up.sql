CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title TEXT NOT NULL,
    description TEXT,
    cover_url TEXT,
    banner_url TEXT,
    status TEXT,
    release_date DATE,
    studio TEXT,
    genre TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
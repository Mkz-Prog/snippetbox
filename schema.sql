-- =============================================================================
-- Snippetbox PostgreSQL Database Schema & Setup
-- =============================================================================

-- 1. Create database (run as postgres superuser if needed)
-- CREATE DATABASE snippetbox WITH ENCODING 'UTF8';

-- 2. Create snippets table
CREATE TABLE IF NOT EXISTS snippets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    expires TIMESTAMPTZ NOT NULL
);

-- 3. Create index on created column for efficient order-by queries
CREATE INDEX IF NOT EXISTS idx_snippets_created ON snippets(created);

-- 4. Seed data (sample snippets)
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    (NOW() AT TIME ZONE 'UTC'),
    (NOW() AT TIME ZONE 'UTC') + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Sōseki',
    (NOW() AT TIME ZONE 'UTC'),
    (NOW() AT TIME ZONE 'UTC') + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning:\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijō',
    (NOW() AT TIME ZONE 'UTC'),
    (NOW() AT TIME ZONE 'UTC') + INTERVAL '7 days'
);

-- 5. User permissions (run as postgres superuser)
-- CREATE USER web WITH PASSWORD 'pass';
-- GRANT ALL PRIVILEGES ON DATABASE snippetbox TO web;
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO web;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO web;

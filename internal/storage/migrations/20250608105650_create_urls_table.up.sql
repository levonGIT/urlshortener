CREATE TABLE urls (
    id integer PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    count integer NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
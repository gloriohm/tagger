-- +goose Up
CREATE TABLE artists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brainz_artist_id TEXT UNIQUE,
    created_at DATETIME CURRENT_TIMESTAMP,
    updated_at DATETIME CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE artists;
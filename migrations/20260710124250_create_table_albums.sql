-- +goose Up
CREATE TABLE albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    brainz_release_id TEXT UNIQUE,
    brainz_releasegroup_id TEXT UNIQUE,
    created_at DATETIME CURRENT_TIMESTAMP,
    updated_at DATETIME CURRENT_TIMESTAMP,

    FOREIGN KEY (path) REFERENCES directory(id)
);
-- +goose Down
DROP TABLE albums;

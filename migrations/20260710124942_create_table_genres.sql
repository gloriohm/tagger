-- +goose Up
CREATE TABLE genres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME CURRENT_TIMESTAMP,
    updated_at DATETIME CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE genres;
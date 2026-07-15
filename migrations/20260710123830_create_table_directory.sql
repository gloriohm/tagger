-- +goose Up
CREATE TABLE directory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    created_at DATETIME CURRENT_TIMESTAMP,
    updated_at DATETIME CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE directory;
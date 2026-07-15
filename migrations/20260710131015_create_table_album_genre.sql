-- +goose Up
CREATE TABLE album_genre (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    album_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    genre_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME CURRENT_TIMESTAMP,

    UNIQUE (album_id, genre_id),
    UNIQUE (album_id, genre_order),

    FOREIGN KEY (album_id) REFERENCES albums(id),
    FOREIGN KEY (genre_id) REFERENCES genres(id)
);
-- +goose Down
DROP TABLE album_genre;

-- +goose Up
CREATE TABLE album_artist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    album_id INTEGER NOT NULL,
    artist_id INTEGER NOT NULL,
    artist_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME CURRENT_TIMESTAMP,

    UNIQUE (album_id, artist_id),
    UNIQUE (album_id, aritst_order),

    FOREIGN KEY (album_id) REFERENCES albums(id),
    FOREIGN KEY (artist_id) REFERENCES artists(id)
);
-- +goose Down
DROP TABLE album_artist;

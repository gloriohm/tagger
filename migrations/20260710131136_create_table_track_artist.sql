-- +goose Up
CREATE TABLE track_artist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    track_id INTEGER NOT NULL,
    artist_id INTEGER NOT NULL,
    artist_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME CURRENT_TIMESTAMP,

    UNIQUE (track_id, artist_id),
    UNIQUE (track_id, aritst_order),

    FOREIGN KEY (track_id) REFERENCES tracks(id),
    FOREIGN KEY (artist_id) REFERENCES artists(id)
);
-- +goose Down
DROP TABLE track_artist;

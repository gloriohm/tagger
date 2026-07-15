-- +goose Up
CREATE TABLE album_track (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    album_id INTEGER NOT NULL,
    track_id INTEGER NOT NULL,
    created_at DATETIME CURRENT_TIMESTAMP,

    UNIQUE (album_id, track_id),

    FOREIGN KEY (album_id) REFERENCES albums(id),
    FOREIGN KEY (track_id) REFERENCES tracks(id)
);
-- +goose Down
DROP TABLE album_track;

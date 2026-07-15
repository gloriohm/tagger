-- +goose Up
CREATE TABLE tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brainz_track_id TEXT UNIQUE,
    brainz_recording_id TEXT UNIQUE,
    artist_credit TEXT NOT NULL,
    track_number INTEGER NOT NULL,
    disc_number INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME CURRENT_TIMESTAMP,
    updated_at DATETIME CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE tracks;
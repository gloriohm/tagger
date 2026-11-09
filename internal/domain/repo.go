package domain

import (
	"database/sql"
	"fmt"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(db *sql.DB) *SQLiteRepo {
	return &SQLiteRepo{db: db}
}

// albums
func (r *SQLiteRepo) InsertAlbum(album Album) (Album, error) {
	err := r.db.QueryRow(
		"INSERT INTO albums (name, slug, mb_release_id, mb_releasegroup_id, release_date, total_discs) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
		album.Name, album.Slug, album.MBReleaseID, album.MBReleaseGroupID, album.ReleaseDate, album.TotalDiscs,
	).Scan(&album.ID)

	if err != nil {
		return Album{}, fmt.Errorf("insert album: %w", err)
	}

	return album, nil
}

func (r *SQLiteRepo) InsertAlbumArtistLinks(album Album) error {
	for _, artist := range album.AlbumArtists {
		_, err := r.db.Exec(
			"INSERT INTO album_artist (album_id, artist_id) VALUES (?, ?)",
			album.ID, artist.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// artists
func (r *SQLiteRepo) SelectArtistByID(id int64) (Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE id = ?`,
			id,
		),
	)
}

func (r *SQLiteRepo) SelectArtistByName(name string) (Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE name = ?`,
			name,
		),
	)
}

func (r *SQLiteRepo) SelectArtistByMBID(mbid string) (Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE mbid = ?`,
			mbid,
		),
	)
}

func (r *SQLiteRepo) InsertArtist(artist Artist) (Artist, error) {
	err := r.db.QueryRow(
		"INSERT INTO artists (name, slug, mbid) VALUES (?, ?, ?) RETURNING id",
		artist.Name, artist.Slug, artist.MBID,
	).Scan(&artist.ID)

	if err != nil {
		return Artist{}, fmt.Errorf("insert genre: %w", err)
	}

	return artist, nil
}

// genres
func (r *SQLiteRepo) InsertGenre(genre Genre) (Genre, error) {
	err := r.db.QueryRow(
		"INSERT INTO genres (name, slug) VALUES (?, ?) RETURNING id",
		genre.Name, genre.Slug,
	).Scan(&genre.ID)

	if err != nil {
		return Genre{}, fmt.Errorf("insert genre: %w", err)
	}

	return genre, nil
}

// labels
func (r *SQLiteRepo) InsertLabel(label Label) (Label, error) {
	err := r.db.QueryRow(
		"INSERT INTO labels (name, slug, mbid) VALUES (?, ?, ?) RETURNING id",
		label.Name, label.Slug, label.MBID,
	).Scan(&label.ID)

	if err != nil {
		return Label{}, fmt.Errorf("insert genre: %w", err)
	}

	return label, nil
}

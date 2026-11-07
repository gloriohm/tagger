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

func (r *SQLiteRepo) SelectArtistByID(id int64) (*Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE id = ?`,
			id,
		),
	)
}

func (r *SQLiteRepo) SelectArtistByName(name string) (*Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE name = ?`,
			name,
		),
	)
}

func (r *SQLiteRepo) SelectArtistByMBID(mbid string) (*Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name
			 FROM artists
			 WHERE mbid = ?`,
			mbid,
		),
	)
}

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

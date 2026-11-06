package domain

import (
	"database/sql"
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

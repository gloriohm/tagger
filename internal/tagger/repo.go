package tagger

import (
	"context"
	"database/sql"
	"errors"
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

func (r *SQLiteRepo) InsertPaths(paths []string) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO directory (path)
		VALUES (?, ?)
		ON CONFLICT (path) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, path := range paths {
		if _, err := stmt.Exec(path, ImportStatusInitialized); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepo) SelectPath(status ImportStatus) (*string, error) {
	q := `SELECT path FROM directory WHERE status = ? ORDER BY path ASC LIMIT 1;`

	var path string
	err := r.db.QueryRow(q, status).Scan(&path)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNoRows
	case err != nil:
		return nil, fmt.Errorf("[SELECT PATH] %w", err)
	default:
		return &path, nil
	}
}

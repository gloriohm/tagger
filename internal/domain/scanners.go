package domain

import "database/sql"

func scanArtist(row *sql.Row) (*Artist, error) {
	var artist Artist

	err := row.Scan(
		&artist.ID,
		&artist.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func scanLabel(row *sql.Row) (*Label, error) {
	var label Label

	err := row.Scan(
		&label.ID,
		&label.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &label, nil
}

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
		"INSERT INTO albums (name, album_artist_credit, mb_release_id, mb_releasegroup_id, release_date, total_discs) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
		album.Name, album.AlbumArtistCredit, album.MBReleaseID, album.MBReleaseGroupID, album.ReleaseDate, album.TotalDiscs,
	).Scan(&album.ID)

	if err != nil {
		return Album{}, fmt.Errorf("insert album: %w", err)
	}

	return album, nil
}

func (r *SQLiteRepo) InsertAlbumArtistLinks(album Album) error {
	for i, artist := range album.AlbumArtists {
		_, err := r.db.Exec(
			"INSERT INTO album_artist (album_id, artist_id, position) VALUES (?, ?, ?)",
			album.ID, artist.ID, i,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLiteRepo) InsertAlbumGenreLinks(album Album) error {
	for _, genre := range album.Genres {
		_, err := r.db.Exec(
			"INSERT INTO album_genre (album_id, genre_id) VALUES (?, ?)",
			album.ID, genre.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLiteRepo) InsertAlbumLabelLinks(album Album) error {
	for _, label := range album.Labels {
		_, err := r.db.Exec(
			"INSERT INTO album_label (album_id, label_id) VALUES (?, ?)",
			album.ID, label.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLiteRepo) SelectAlbumByID(id int64) (Album, error) {
	q := `
		SELECT
			id,
			mb_release_id,
			mb_release_group_id,
			name,
			album_artist_credit,
			release_date,
			total_discs
		FROM albums
		WHERE id = ?;
	`
	var album Album
	err := r.db.QueryRow(q, id).Scan(
		&album.ID,
		&album.MBReleaseID,
		&album.MBReleaseGroupID,
		&album.Name,
		&album.AlbumArtistCredit,
		&album.ReleaseDate,
		&album.TotalDiscs,
	)

	return album, err
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
			`SELECT id, name, mbid
			 FROM artists
			 WHERE name = ?`,
			name,
		),
	)
}

func (r *SQLiteRepo) SelectArtistByMBID(mbid string) (Artist, error) {
	return scanArtist(
		r.db.QueryRow(
			`SELECT id, name, mbid
			 FROM artists
			 WHERE mbid = ?`,
			mbid,
		),
	)
}

func (r *SQLiteRepo) PopulateAlbumArtists(album *Album) error {
	q := `
		SELECT 
			id, mbid, name
		FROM artists a
		JOIN album_artists aa ON aa.artist_id = a.id
		WHERE aa.album_id = ?
		ORDER BY aa.position
	`

	rows, err := r.db.Query(q, album.ID)
	if err != nil {
		return fmt.Errorf("query album artists: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var artist Artist
		if err := rows.Scan(
			&artist.ID,
			&artist.MBID,
			&artist.Name,
		); err != nil {
			return fmt.Errorf("scan album artist: %w", err)
		}

		album.AlbumArtists = append(album.AlbumArtists, artist)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate album artists: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) PopulateTrackArtists(album *Album) error {
	q := `
		SELECT
			ta.track_id,
			a.id,
			a.mbid,
			a.name
		FROM track_artists ta
		JOIN artists a ON a.id = ta.artist_id
		JOIN tracks t ON t.id = ta.track_id
		WHERE t.album_id = ?
		ORDER BY t.disc_number, t.track_number, ta.position;
	`
	rows, err := r.db.Query(q, album.ID)
	if err != nil {
		return fmt.Errorf("query track artists: %w", err)
	}
	defer rows.Close()

	trackIdx := make(map[int64]int, len(album.Tracks))

	for i := range album.Tracks {
		trackIdx[album.Tracks[i].ID] = i
	}

	for rows.Next() {
		var trackID int64
		var artist Artist

		if err := rows.Scan(
			&trackID,
			&artist.ID,
			&artist.MBID,
			&artist.Name,
		); err != nil {
			return fmt.Errorf("scan track artist: %w", err)
		}

		i, ok := trackIdx[trackID]
		if !ok {
			return fmt.Errorf("track %d not found in album", trackID)
		}

		album.Tracks[i].Artists = append(
			album.Tracks[i].Artists,
			artist,
		)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate track artists: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) InsertArtist(artist Artist) (Artist, error) {
	err := r.db.QueryRow(
		"INSERT INTO artists (name, mbid) VALUES (?, ?) RETURNING id",
		artist.Name, artist.MBID,
	).Scan(&artist.ID)

	if err != nil {
		return Artist{}, fmt.Errorf("insert artist: %w", err)
	}

	return artist, nil
}

// tracks
func (r *SQLiteRepo) InsertTracks(tracks []Track) ([]Track, error) {
	for i, t := range tracks {
		err := r.db.QueryRow(
			"INSERT INTO tracks (mb_track_id, mb_recording_id, album_id, name, artist_credit, track_number, disc_number) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id",
			t.MBRecordingID, t.MBTrackID, t.AlbumID, t.Name, t.ArtistCredit, t.TrackNumber, t.DiscNumber,
		).Scan(tracks[i].ID)

		if err != nil {
			return nil, fmt.Errorf("insert tracks: %W", err)
		}

	}

	return tracks, nil
}

func (r *SQLiteRepo) InsertTrackArtistLinks(tracks []Track) error {
	for _, track := range tracks {
		for i, artist := range track.Artists {
			_, err := r.db.Exec(
				"INSERT INTO track_artist (track_id, artist_id, position) VALUES (?, ?, ?)",
				track.ID, artist.ID, i,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *SQLiteRepo) PopulateTracks(album *Album) error {
	q := `
		SELECT 
			id, mb_track_id, mb_recording_id, album_id, name, artist_credit, track_number, disc_number
		FROM tracks
		WHERE album_id = ?
		ORDER BY disc_number, track_number
	`

	rows, err := r.db.Query(q, album.ID)
	if err != nil {
		return fmt.Errorf("query tracks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		if err := rows.Scan(
			&track.ID,
			&track.MBRecordingID,
			&track.MBTrackID,
			&track.AlbumID,
			&track.Name,
			&track.ArtistCredit,
			&track.TrackNumber,
			&track.DiscNumber,
		); err != nil {
			return fmt.Errorf("scan tracks: %w", err)
		}

		album.Tracks = append(album.Tracks, track)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate tracks: %w", err)
	}

	return nil
}

// genres
func (r *SQLiteRepo) InsertGenre(genre Genre) (Genre, error) {
	err := r.db.QueryRow(
		"INSERT INTO genres (name) VALUES (?) RETURNING id",
		genre.Name,
	).Scan(&genre.ID)

	if err != nil {
		return Genre{}, fmt.Errorf("insert genre: %w", err)
	}

	return genre, nil
}

func (r *SQLiteRepo) ListGenres() ([]Genre, error) {
	q := `
		SELECT
			id, name
		FROM genres
		ORDER BY name;
	`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("list genres: %w", err)
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var genre Genre
		if err := rows.Scan(
			&genre.ID,
			&genre.Name,
		); err != nil {
			return nil, fmt.Errorf("scan genres: %w", err)
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate genres: %w", err)
	}

	return genres, nil
}

func (r *SQLiteRepo) PopulateGenres(album *Album) error {
	q := `
		SELECT 
			id, name
		FROM genres g
		JOIN album_genres ag ON ag.genre_id = g.album_id
		WHERE ag.album_id = ?
		ORDER BY g.name;
	`

	rows, err := r.db.Query(q, album.ID)
	if err != nil {
		return fmt.Errorf("query genres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var genre Genre
		if err := rows.Scan(
			&genre.ID,
			&genre.Name,
		); err != nil {
			return fmt.Errorf("scan genres: %w", err)
		}

		album.Genres = append(album.Genres, genre)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate genres: %w", err)
	}

	return nil
}

// labels
func (r *SQLiteRepo) InsertLabel(label Label) (Label, error) {
	err := r.db.QueryRow(
		"INSERT INTO labels (name, mbid) VALUES (?, ?) RETURNING id",
		label.Name, label.MBID,
	).Scan(&label.ID)

	if err != nil {
		return Label{}, fmt.Errorf("insert genre: %w", err)
	}

	return label, nil
}

func (r *SQLiteRepo) PopulateLabels(album *Album) error {
	q := `
		SELECT 
			id, mbid, name
		FROM labels l
		JOIN album_labels al ON al.label_id = l.album_id
		WHERE al.album_id = ?
		ORDER BY l.name;
	`

	rows, err := r.db.Query(q, album.ID)
	if err != nil {
		return fmt.Errorf("query genres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var label Label
		if err := rows.Scan(
			&label.ID,
			&label.MBID,
			&label.Name,
		); err != nil {
			return fmt.Errorf("scan labels: %w", err)
		}

		album.Labels = append(album.Labels, label)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate labels: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) SelectLabelByMBID(mbid string) (Label, error) {
	q := `SELECT
			id, mbid, name
		FROM labels
		WHERE mbid = ?
	`
	var label Label
	err := r.db.QueryRow(q, mbid).Scan(&label.ID, &label.MBID, &label.Name)

	return label, err
}

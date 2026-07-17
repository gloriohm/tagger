package filer

import (
	"slices"
)

func (a Album) ValidateAlbum() error {
	if len(a) == 0 {
		return ErrNoTracks
	}

	if err := a.validateAlbumArtist(); err != nil {
		return err
	}

	if err := a.validateAlbumName(); err != nil {
		return err
	}

	if err := a.validateAlbumReleaseID(); err != nil {
		return err
	}

	if err := a.validateAlbumReleaseID(); err != nil {
		return err
	}

	return nil
}

func (a Album) validateAlbumArtist() error {
	// extract album artist from first track; return error if empty
	target := a[0].AlbumArtists
	if len(target) == 0 {
		return ErrNoValue
	}

	// check that album artist is the same across the entire album
	for _, t := range a[1:] {
		if !slices.Equal(target, t.AlbumArtists) {
			return ErrInconsistentValues
		}
	}

	return nil
}

func (a Album) validateConsistentValue(value func(Track) string) error {
	target := value(a[0])

	for _, track := range a[1:] {
		if value(track) != target {
			return ErrInconsistentValues
		}
	}

	return nil
}

func (a Album) validateAlbumName() error {
	return a.validateConsistentValue(func(track Track) string {
		return track.Album
	})
}

func (a Album) validateAlbumReleaseID() error {
	return a.validateConsistentValue(func(track Track) string {
		return track.MusicBrainzAlbumID
	})
}

func (a Album) validateAlbumReleaseGroupID() error {
	return a.validateConsistentValue(func(track Track) string {
		return track.MusicBrainzReleaseGroupID
	})
}

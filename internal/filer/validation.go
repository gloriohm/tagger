package filer

import (
	"slices"
)

func (a Album) validateAlbumArtist() error {
	if len(a) == 0 {
		return ErrNoTracks
	}

	// extract album artist from first track; return error if empty
	target := a[0].AlbumArtists
	if len(target) == 0 {
		return ErrNoValue
	}

	// check that album artist is the same across the entire album
	for _, t := range a {
		if !slices.Equal(target, t.AlbumArtists) {
			return ErrInconsistentValues
		}
	}

	return nil
}

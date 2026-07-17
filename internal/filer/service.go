package filer

import (
	"os"
	"path/filepath"
	"strings"

	"go.senan.xyz/taglib"
)

func LoadAlbum(path string) (Identifiers, error) {
	var album Album
	entries, err := os.ReadDir(path)
	if err != nil {
		return Identifiers{}, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".flac" && ext != ".opus" {
			continue
		}

		track, err := LoadTrack(entry.Name())
		if err != nil {
			return Identifiers{}, err
		}

		album = append(album, track)
	}

	if err := album.ValidateAlbum(); err != nil {
		return Identifiers{}, err
	}

	ids, err := parseIdentifiers(album[0])

	return ids, nil
}

func LoadTrack(path string) (Track, error) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		return Track{}, err
	}

	track := parseTrack(tags)

	return track, nil
}

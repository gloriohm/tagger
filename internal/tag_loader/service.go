package tag_loader

import (
	"os"
	"path/filepath"
	"strings"
	"tagger/internal/compare"

	"go.senan.xyz/taglib"
)

func LoadAlbum(path string) (AlbumIdentifiers, error) {
	var album Album
	entries, err := os.ReadDir(path)
	if err != nil {
		return AlbumIdentifiers{}, err
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
			return AlbumIdentifiers{}, err
		}

		album = append(album, track)
	}

	if err := album.ValidateAlbum(); err != nil {
		return AlbumIdentifiers{}, err
	}

	ids, err := parseAlbumIdentifiers(album[0])

	// append tracks so tags are available without loading again
	ids.AlbumTags = album

	// derive total discs for validation later
	ids.TotalDiscs = album.DeriveTotalDiscs()

	return ids, nil
}

func LoadTrack(path string) (Track, error) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		return Track{}, err
	}

	props, err := taglib.ReadProperties(path)

	track := parseTrack(path, tags, props)

	return track, nil
}

func TagAlbum(path string, album []map[string][]string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var musicFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".flac" || ext == ".opus" {
			musicFiles = append(musicFiles, entry.Name())
		}
	}

	if len(musicFiles) != len(album) {
		return ErrDifferentTrackCount
	}

	for i := range musicFiles {
		tags, err := taglib.ReadTags(musicFiles[i])
		if err != nil {
			return err
		}

		distance := compare.CompareTrackName(tags[taglib.Title][0], album[i][taglib.Title][0])
		if distance > 0.5 {
			return ErrCheckTrackOrder
		}

		if err := taglib.WriteTags(musicFiles[i], album[i], taglib.Clear); err != nil {
			return err
		}
	}

	return nil
}

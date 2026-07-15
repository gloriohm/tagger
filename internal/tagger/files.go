package tagger

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func findAllAlbumFolders(root string) ([]string, error) {
	var albumPaths []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".flac" || ext == ".opus" {
				albumPaths = append(albumPaths, path)

				// return after one file of the appropriate type is found
				return filepath.SkipDir
			}
		}

		return nil
	})

	return albumPaths, err
}

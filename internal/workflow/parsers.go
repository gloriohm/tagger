package workflow

import (
	"cmp"
	"slices"
	"strconv"
	"tagger/internal/compare"
	"tagger/internal/mbfetcher"
	"tagger/internal/tag_loader"
)

func localAlbumToCompare(in tag_loader.AlbumIdentifiers) compare.Album {
	out := compare.Album{
		Title:      in.AlbumName,
		Artist:     in.AlbumArtist,
		TotalDiscs: in.TotalDiscs,
	}

	// sort based on disc and track number to ensure order is chronological
	slices.SortFunc(in.AlbumTags, func(a, b tag_loader.Track) int {
		return cmp.Or(
			cmp.Compare(a.DiscNumber, b.DiscNumber),
			cmp.Compare(a.TrackNumber, b.TrackNumber),
		)
	})

	for _, t := range in.AlbumTags {
		tNumber, _ := strconv.Atoi(t.TrackNumber)
		dNumber, _ := strconv.Atoi(t.DiscNumber)

		track := compare.Track{
			Title:       t.Title,
			TrackNumber: tNumber,
			DiscNumber:  dNumber,
			Duration:    t.Duration,
		}
		out.Tracks = append(out.Tracks, track)
	}

	return out
}

func mbAlbumToCompare(in mbfetcher.Album) compare.Album {
	out := compare.Album{
		Title:      in.Name,
		TotalDiscs: in.TotalDiscs,
	}

	if len(in.AlbumArtists) > 0 {
		out.Artist = in.AlbumArtists[0].Name
	}

	// sort based on disc and track number to ensure order is chronological
	slices.SortFunc(in.Tracks, func(a, b mbfetcher.Track) int {
		return cmp.Or(
			cmp.Compare(a.DiscNumber, b.DiscNumber),
			cmp.Compare(a.TrackNumber, b.TrackNumber),
		)
	})

	for _, t := range in.Tracks {
		track := compare.Track{
			Title:       t.Name,
			TrackNumber: t.TrackNumber,
			DiscNumber:  t.DiscNumber,
			Duration:    t.Duration,
		}
		out.Tracks = append(out.Tracks, track)
	}

	return out
}

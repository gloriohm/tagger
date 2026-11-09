package domain

import (
	"tagger/internal/mbfetcher"
)

func ParseAlbum(in mbfetcher.Album) Album {
	var a Album
	a.MBReleaseID = in.MBReleaseID
	a.MBReleaseGroupID = in.MBReleaseGroupID
	a.Name = in.Name
	a.Slug = toSlug(in.Name)
	a.ReleaseDate = in.ReleaseDate
	a.TotalDiscs = in.TotalDiscs

	return a
}

func ParseLabels(in []mbfetcher.LabelInfo) []Label {
	labels := make([]Label, len(in))

	for _, l := range in {
		label := Label{
			MBID: l.MBID,
			Name: l.Name,
			Slug: toSlug(l.Name),
		}

		labels = append(labels, label)
	}

	return labels
}

func ParseArtists(in []mbfetcher.Artist) []Artist {
	artists := make([]Artist, len(in))

	for _, a := range in {
		artist := Artist{
			MBID: a.MBID,
			Name: a.Name,
			Slug: toSlug(a.Name),
		}

		artists = append(artists, artist)
	}

	return artists
}

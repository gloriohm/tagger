package filer

import "go.senan.xyz/taglib"

func parseTrack(tags map[string][]string) Track {
	track := Track{
		Title:                     firstTag(tags, taglib.Title),
		Artists:                   tags[taglib.Artist],
		AlbumArtists:              tags[taglib.AlbumArtist],
		Album:                     firstTag(tags, taglib.Album),
		Date:                      firstTag(tags, taglib.Date),
		TrackNumber:               firstTag(tags, taglib.TrackNumber),
		DiscNumber:                firstTag(tags, taglib.DiscNumber),
		Label:                     firstTag(tags, taglib.Label),
		BPM:                       firstTag(tags, taglib.BPM),
		CatalogNumber:             firstTag(tags, taglib.CatalogNumber),
		MusicBrainzTrackID:        firstTag(tags, taglib.MusicBrainzTrackID),
		MusicBrainzReleaseTrackID: firstTag(tags, taglib.MusicBrainzReleaseTrackID),
		MusicBrainzAlbumID:        firstTag(tags, taglib.MusicBrainzAlbumID),
		MusicBrainzReleaseGroupID: firstTag(tags, taglib.MusicBrainzReleaseGroupID),
		MusicBrainzArtistID:       firstTag(tags, taglib.MusicBrainzArtistID),
		MusicBrainzAlbumArtistID:  firstTag(tags, taglib.MusicBrainzAlbumArtistID),
	}
	return track
}

func firstTag(tags map[string][]string, key string) string {
	values := tags[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func parseIdentifiers(t Track) (Identifiers, error) {
	if len(t.AlbumArtists) == 0 {
		return Identifiers{}, ErrNoValue
	}

	return Identifiers{
		ReleaseID:      t.MusicBrainzAlbumID,
		ReleaseGroupID: t.MusicBrainzReleaseGroupID,
		AlbumName:      t.Album,
		AlbumArtist:    t.AlbumArtists[0],
	}, nil
}

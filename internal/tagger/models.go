package tagger

import "time"

type ComparableAlbum struct {
	MBReleaseID  string
	Name         string
	AlbumArtists []string
	Tracks       []ComparableTrack
}

type ComparableTrack struct {
	MBTrackID string
	Name      string
	Position  int
	Length    time.Duration
}

type ImportStatus string

var (
	ImportStatusInitialized  ImportStatus = "initialized"
	ImportStatusPendingInput ImportStatus = "pending_input"
	ImportStatusFinished     ImportStatus = "finished"
)

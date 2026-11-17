package tag_loader

import "errors"

var (
	ErrNoTracks            = errors.New("no tracks")
	ErrNoValue             = errors.New("track(s) missing essential metadata")
	ErrInconsistentValues  = errors.New("values that should be consistent across tracks differed")
	ErrDifferentTrackCount = errors.New("mismatch between music files in directory and files to tag")
	ErrCheckTrackOrder     = errors.New("suspected mismatch between files and album order")
)

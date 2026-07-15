package filer

import "errors"

var (
	ErrNoTracks           = errors.New("no tracks")
	ErrNoValue            = errors.New("track(s) missing essential metadata")
	ErrInconsistentValues = errors.New("values that should be consistent across tracks differed")
)

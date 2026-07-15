package tagger

import (
	"errors"
)

var (
	ErrDuplicateEntry = errors.New("entity already exists")
	ErrNoRows         = errors.New("no results")
)

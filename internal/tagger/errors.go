package tagger

import (
	"errors"
)

var (
	ErrDuplicateEntry = errors.New("entity already exists")
	ErrNoRows         = errors.New("no results")
	ErrEmptyValue     = errors.New("required value was empty")
	ErrStoreNotSet    = errors.New("store is nil") // if this error occurs, there is a bug in the workflow
)

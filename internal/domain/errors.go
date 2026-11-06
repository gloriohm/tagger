package domain

import (
	"errors"
)

var (
	ErrDuplicateEntry = errors.New("entity already exists")
	ErrNoRows         = errors.New("no results")
	ErrEmptyValue     = errors.New("required value was empty")
)

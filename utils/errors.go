package utils

import "errors"

// ErrEntityNotFound is returned when a requested entity could not be retrieved
var ErrEntityNotFound = errors.New("entity not found")

// ErrEntityExists is returned when an entity is already present
var ErrEntityExists = errors.New("entity exists")

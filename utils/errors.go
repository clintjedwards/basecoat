package utils

import "errors"

// ErrFormulaNotFound is thrown when a requested entity is not found
var ErrFormulaNotFound = errors.New("formula not found")

// ErrFormulaExists  is thrown when an entity is found before adding another
var ErrFormulaExists = errors.New("formula already exists")

// ErrJobNotFound is thrown when a requested entity is not found
var ErrJobNotFound = errors.New("job not found")

// ErrJobExists is thrown when an entity is found before adding another
var ErrJobExists = errors.New("job already exists")

// ErrUserExists is thrown when an entity is found before adding another
var ErrUserExists = errors.New("user already exists")

// ErrUserNotFound is thrown when a requested entity is not found
var ErrUserNotFound = errors.New("user not found")

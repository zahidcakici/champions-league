package handlers

import "errors"

// Custom validation errors
var (
	ErrInvalidHomeScore = errors.New("home score must be non-negative")
	ErrInvalidAwayScore = errors.New("away score must be non-negative")
)

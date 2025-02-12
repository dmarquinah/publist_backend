package errors

import "errors"

var (
	ErrItemNotFound     = errors.New("item not found")
	ErrPlaylistNotFound = errors.New("playlist not found")
	ErrTrackNotFound    = errors.New("track not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidName      = errors.New("invalid playlist name")
	ErrNameTooLong      = errors.New("playlist name too long")
	ErrInvalidPosition  = errors.New("invalid track position")
	// Add more custom errors as needed
)

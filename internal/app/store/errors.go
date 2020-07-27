package store

import "errors"

var (
	// ErrInvalidToken ...
	ErrInvalidToken = errors.New("Invalid token")
	// ErrDocumentNotFound ...
	ErrDocumentNotFound = errors.New("Document not found")
)

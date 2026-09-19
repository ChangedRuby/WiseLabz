// Package storeerr holds the sentinel errors returned by store repository
// methods. It is a dependency-free leaf so both the store and the HTTP layer
// (httputil) can reference them without an import cycle. The store package
// re-exports these values, so errors.Is(err, store.ErrNotFound) keeps working.
package storeerr

import "errors"

// Sentinel errors returned by repository methods.
var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource already exists")
	ErrVersionConflict = errors.New("version conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
)

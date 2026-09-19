package store

import "github.com/WiseLabz/wiselabz/internal/storeerr"

// Sentinel errors returned by repository methods. They are aliases of the
// values in storeerr so the HTTP layer can map them without importing store.
var (
	ErrNotFound        = storeerr.ErrNotFound
	ErrConflict        = storeerr.ErrConflict
	ErrVersionConflict = storeerr.ErrVersionConflict
	ErrUnauthorized    = storeerr.ErrUnauthorized
	ErrForbidden       = storeerr.ErrForbidden
)

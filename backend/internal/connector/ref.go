package connector

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// ValidateRefSegment checks that s is safe to interpolate as a single path
// segment of an upstream API URL. It rejects empty values, dot segments and
// any character that could alter the request path or query: '/', '\', '?',
// '#', '%', whitespace, control characters and non-ASCII.
func ValidateRefSegment(s string) error {
	if s == "" {
		return fmt.Errorf("reference must not be empty")
	}
	if s == "." || s == ".." || strings.Contains(s, "..") {
		return fmt.Errorf("reference %q contains a path traversal sequence", s)
	}
	for _, r := range s {
		if r <= 0x20 || r >= 0x7f || strings.ContainsRune(`/\?#%;`, r) {
			return fmt.Errorf("reference %q contains an invalid character", s)
		}
	}
	return nil
}

// ValidateCompositeRef validates a reference made of segments joined by '/'
// (e.g. "<zoneID>/<recordID>"); each segment is validated individually.
func ValidateCompositeRef(s string) error {
	if s == "" {
		return nil
	}
	for _, p := range strings.Split(s, "/") {
		if err := ValidateRefSegment(p); err != nil {
			return err
		}
	}
	return nil
}

// PathSegment validates s with ValidateRefSegment and returns it escaped for
// use as a URL path segment.
func PathSegment(s string) (string, error) {
	if err := ValidateRefSegment(s); err != nil {
		return "", err
	}
	return url.PathEscape(s), nil
}

// ValidateUnixSocketPath checks a docker unix:// socket path: it must be
// absolute, already clean (no '..' or redundant elements) and free of
// control characters.
func ValidateUnixSocketPath(p string) error {
	if p == "" {
		return fmt.Errorf("unix socket path must not be empty")
	}
	for _, r := range p {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("unix socket path contains control characters")
		}
	}
	if !filepath.IsAbs(p) {
		return fmt.Errorf("unix socket path %q must be absolute", p)
	}
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." {
			return fmt.Errorf("unix socket path %q must not contain '..'", p)
		}
	}
	if filepath.Clean(p) != p {
		return fmt.Errorf("unix socket path %q must be clean", p)
	}
	return nil
}

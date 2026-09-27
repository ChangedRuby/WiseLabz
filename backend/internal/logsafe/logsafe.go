// Package logsafe sanitizes user-controlled strings before they reach a logger.
package logsafe

import "strings"

// Sanitize strips CR/LF so a value can't forge or split log lines.
func Sanitize(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}

// Err returns err's message sanitized for logging. Error messages often wrap
// user input (IDs, names, upstream responses), so log them through this. A
// nil error yields "".
func Err(err error) string {
	if err == nil {
		return ""
	}
	return Sanitize(err.Error())
}

package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// verificationsFile is an append-only, newest-last JSON-lines log of every
// VerificationResult recorded in a backup directory, making verify results
// queryable (by `backup verify`/the scheduled job, and by operators via
// ListVerifications) without a database migration of its own.
const verificationsFile = "verifications.jsonl"

// RecordVerification appends res as one JSON line to {dir}/verifications.jsonl,
// creating dir (0o700) if needed. The file is append-only, so a concurrent
// manual `backup verify` and the scheduled job can both write to it safely.
func RecordVerification(dir string, res VerificationResult) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}
	data, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("marshal verification result: %w", err)
	}
	f, err := os.OpenFile(filepath.Join(dir, verificationsFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open verifications log: %w", err)
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write verification result: %w", err)
	}
	return nil
}

// ListVerifications reads back recorded verification results from
// {dir}/verifications.jsonl, newest first. A non-positive limit returns
// every recorded result. Returns (nil, nil) if no verification has run yet.
func ListVerifications(dir string, limit int) ([]VerificationResult, error) {
	data, err := os.ReadFile(filepath.Join(dir, verificationsFile))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read verifications log: %w", err)
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	results := make([]VerificationResult, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		var r VerificationResult
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			return nil, fmt.Errorf("parse verification result: %w", err)
		}
		results = append(results, r)
	}

	// Reverse to newest-first.
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

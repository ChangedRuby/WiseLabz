// Command backup verifies and restores WiseLabz backup bundles from the
// command line, without going through the HTTP API.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var code int
	switch os.Args[1] {
	case "verify":
		code = runVerify(os.Args[2:])
	case "restore":
		code = runRestore(os.Args[2:])
	default:
		usage()
		code = 1
	}
	os.Exit(code)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: backup <verify|restore> [flags]")
	fmt.Fprintln(os.Stderr, "  backup verify  [-dir DIR] [-file PATH]   checksum + row-count verify a bundle")
	fmt.Fprintln(os.Stderr, "  backup restore -file PATH [-yes]         verify, then import a bundle into the configured database")
}

// runVerify implements `backup verify`: checksums the bundle against its
// manifest sidecar and restores it into a throwaway in-memory database to
// confirm every row round-trips, without touching the real database.
func runVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	dir := fs.String("dir", "./data/backups", "backup directory to scan for the latest bundle (ignored if -file is set)")
	file := fs.String("file", "", "verify this specific bundle file instead of the latest one in -dir")
	_ = fs.Parse(args)

	path := *file
	if path == "" {
		p, err := backup.LatestBundle(*dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "verify: %v\n", err)
			return 1
		}
		path = p
	}

	res := backup.VerifyBundleFile(context.Background(), path)
	if err := backup.RecordVerification(filepath.Dir(path), res); err != nil {
		fmt.Fprintf(os.Stderr, "verify: record result: %v\n", err)
	}

	if res.Status != "pass" {
		fmt.Printf("FAIL  %s\n  error: %s\n", path, res.Error)
		return 1
	}
	fmt.Printf("PASS  %s\n  checksum: %s\n  counts:   %s\n", path, res.Checksum, formatCounts(res.ActualCounts))
	return 0
}

// runRestore implements `backup restore`: verifies the bundle (same check as
// `backup verify`), then imports it into the database from the loaded
// config, the same additive/idempotent Import used by the API and the
// scheduled backup job — records whose ID already exists are left untouched.
func runRestore(args []string) int {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	file := fs.String("file", "", "bundle file to restore (required)")
	yes := fs.Bool("yes", false, "skip the confirmation prompt")
	_ = fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, "restore: -file is required")
		return 1
	}

	fmt.Printf("Verifying %s before restore...\n", *file)
	res := backup.VerifyBundleFile(context.Background(), *file)
	if res.Status != "pass" {
		fmt.Fprintf(os.Stderr, "restore: verification failed, aborting: %s\n", res.Error)
		return 1
	}
	fmt.Printf("Verification passed (checksum %s, %s).\n", res.Checksum, formatCounts(res.ActualCounts))

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "restore: load config: %v\n", err)
		return 1
	}

	if !*yes && !confirm(fmt.Sprintf("This will import %s into %s (%s). Existing records with matching IDs are left untouched. Continue?", *file, cfg.DB.Driver, cfg.DB.DSN)) {
		fmt.Println("Aborted.")
		return 1
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	db, err := store.OpenDB(cfg.DB.Driver, cfg.DB.DSN, store.PoolConfig{
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime(),
		ConnMaxIdleTime: cfg.DB.ConnMaxIdleTime(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "restore: open database: %v\n", err)
		return 1
	}
	defer func() { _ = db.Close() }()

	if err := store.RunMigrations(db, cfg.DB.Driver, logger); err != nil {
		fmt.Fprintf(os.Stderr, "restore: run migrations: %v\n", err)
		return 1
	}

	s := store.New(db, cfg.DB.Driver)
	result, err := backup.ImportFromFile(context.Background(), s, *file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "restore: import: %v\n", err)
		return 1
	}

	fmt.Println("Restore complete.")
	fmt.Printf("  connectors:       imported=%d skipped=%d\n", result.Connectors.Imported, result.Connectors.Skipped)
	fmt.Printf("  docs:             imported=%d skipped=%d\n", result.Docs.Imported, result.Docs.Skipped)
	fmt.Printf("  docVersions:      imported=%d skipped=%d\n", result.DocVersions.Imported, result.DocVersions.Skipped)
	fmt.Printf("  templates:        imported=%d skipped=%d\n", result.Templates.Imported, result.Templates.Skipped)
	fmt.Printf("  templateSections: imported=%d skipped=%d\n", result.TemplateSections.Imported, result.TemplateSections.Skipped)
	fmt.Println("Connector/notification secrets are redacted from backups and must be re-entered; see docs/BACKUP_RECOVERY.md.")
	return 0
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N] ", prompt)
	answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

func formatCounts(counts map[string]int) string {
	if len(counts) == 0 {
		return "(no counts)"
	}
	parts := make([]string, 0, len(counts))
	for _, entity := range []string{"connectors", "docs", "docVersions", "templates", "templateSections"} {
		if v, ok := counts[entity]; ok {
			parts = append(parts, fmt.Sprintf("%s=%d", entity, v))
		}
	}
	return strings.Join(parts, " ")
}

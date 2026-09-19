package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/WiseLabz/wiselabz/internal/config"
)

const configUsage = "Usage: server config <validate|print --redacted>\n"

// runConfigCommand implements `server config ...` and returns the process exit code.
func runConfigCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(stderr, configUsage)
		return 2
	}
	switch args[0] {
	case "validate":
		cfg, err := config.Load()
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Config invalid: %v\n", err)
			return 1
		}
		if err := cfg.Validate(); err != nil {
			_, _ = fmt.Fprintf(stderr, "Config invalid:\n%v\n", err)
			return 1
		}
		_, _ = fmt.Fprintln(stdout, "Config OK")
		return 0
	case "print":
		// Only the redacted view exists: refusing to print without the flag
		// keeps the command from becoming a secret-dumping habit.
		if len(args) != 2 || args[1] != "--redacted" {
			_, _ = fmt.Fprint(stderr, "config print requires --redacted\n"+configUsage)
			return 2
		}
		cfg, err := config.Load()
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Failed to load config: %v\n", err)
			return 1
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(cfg.Redacted()); err != nil {
			_, _ = fmt.Fprintf(stderr, "Failed to print config: %v\n", err)
			return 1
		}
		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "Unknown config command: %s\n%s", args[0], configUsage)
		return 2
	}
}

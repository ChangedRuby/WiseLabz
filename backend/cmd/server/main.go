// Package main is the entry point for the WiseLabz server.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/quality"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/web"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		runHealthcheck()
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize structured logger
	logger := newLogger(cfg.Log)
	slog.SetDefault(logger)

	logger.Info("WiseLabz server starting",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"db_driver", cfg.DB.Driver,
	)

	if len(cfg.Auth.Secret) < 32 {
		logger.Error("WISELABZ_AUTH_SECRET is missing or too short: refusing to start without a strong JWT signing secret (min 32 chars)")
		os.Exit(1)
	}

	if _, err := crypto.DecodeKey(cfg.Encryption.Key); err != nil {
		logger.Error("WISELABZ_ENCRYPTION_KEY is missing or invalid: refusing to start without a valid base64-encoded 32-byte encryption key (e.g. `openssl rand -base64 32`)", "error", err)
		os.Exit(1)
	}

	if cfg.Server.Origin == "" {
		logger.Error("WISELABZ_SERVER_ORIGIN is missing: refusing to start without an explicit allowed CORS origin")
		os.Exit(1)
	}

	// Open database
	db, err := store.OpenDB(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close() //nolint:errcheck

	// Run migrations
	if err := store.RunMigrations(db, cfg.DB.Driver, logger); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize store
	s := store.New(db, cfg.DB.Driver)

	// Create root context that cancels on interrupt
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize singleton configs + seed admin if needed
	adminPassword := os.Getenv("WISELABZ_ADMIN_PASSWORD")
	if err := s.Init(ctx, adminPassword); err != nil {
		logger.Error("Failed to initialize store", "error", err)
		os.Exit(1)
	}

	logger.Info("Store initialized, database ready")

	// Initialize JWT service
	jwtSvc := auth.NewService(
		cfg.Auth.Secret,
		cfg.Auth.AccessTokenTTLDuration(),
		cfg.Auth.RefreshTokenTTLDuration(),
	)

	// Initialize WebSocket hub (must be created before sync engine so sync
	// can broadcast progress events).
	wsHub := ws.NewHub(cfg.Server.Origin)
	go wsHub.Run()
	logger.Info("WebSocket hub started")

	// Initialize notification dispatcher (must precede sync engine so it can
	// notify on alert creation)
	notifDispatcher := notifications.NewDispatcher(s, wsHub)

	// Initialize engines
	qualityChecker := quality.NewChecker(s, wsHub)
	syncEngine := sync.NewEngine(s, wsHub, notifDispatcher, qualityChecker, cfg.Encryption.Key)
	docEngine := doc.NewEngine(s)

	aiRegistry := ai.NewRegistry()
	ai.RegisterOpenAICompatible(aiRegistry)
	ai.RegisterClaude(aiRegistry)

	// Determine backup directory: use configured value, or compute from DB DSN
	backupDir := cfg.Backup.Dir
	if backupDir == "" {
		// Default: ./data/backups, or extract from SQLite path if configured
		backupDir = "./data/backups"
	}

	// Seed the backup schedule from config defaults if no row exists yet.
	// api.NewRouter's InitBackupJob reads it back and registers the cron job.
	if _, err := s.GetBackupSchedule(ctx); err != nil {
		logger.Info("Initializing backup schedule with defaults")
		defaultSched := store.BackupSchedule{
			CronExpr:    cfg.Backup.CronExpr,
			MaxBackups:  cfg.Backup.MaxBackups,
			MaxAgeHours: cfg.Backup.MaxAgeHours,
			Enabled:     cfg.Backup.Enabled,
		}
		if err := s.UpsertBackupSchedule(ctx, defaultSched); err != nil {
			logger.Error("Failed to initialize backup schedule", "error", err)
			os.Exit(1)
		}
	}

	// Seed the retention settings from config defaults if no row exists yet.
	// api.NewRouter's InitRetentionJob reads it back and registers the cron job.
	if _, err := s.GetRetentionSettings(ctx); err != nil {
		logger.Info("Initializing retention settings with defaults")
		defaultRetention := store.RetentionSettings{
			SnapshotDays:   cfg.Retention.SnapshotDays,
			DocVersionDays: cfg.Retention.DocVersionDays,
			AlertDays:      cfg.Retention.AlertDays,
			SyncRunDays:    cfg.Retention.SyncRunDays,
			AuditDays:      cfg.Retention.AuditDays,
			CronExpr:       cfg.Retention.CronExpr,
		}
		if err := s.UpsertRetentionSettings(ctx, defaultRetention); err != nil {
			logger.Error("Failed to initialize retention settings", "error", err)
			os.Exit(1)
		}
	}

	// Start scheduler for quality, sync, and backup jobs. The retention job
	// itself is registered by api.NewRouter (via the system handler's
	// InitRetentionJob), same reasoning as the backup job below.
	jobRunner := scheduler.New(logger)
	if _, err := jobRunner.AddJob("quality", cfg.Quality.CronExpr, func(jobCtx context.Context) {
		quality.RunStaleSweepOnce(jobCtx, s, wsHub, logger)
	}); err != nil {
		logger.Error("Failed to add quality job", "error", err)
		os.Exit(1)
	}
	if _, err := jobRunner.AddJob("sync", cfg.Sync.PollCronExpr, func(jobCtx context.Context) {
		syncEngine.RunDueSyncs(jobCtx, logger)
	}); err != nil {
		logger.Error("Failed to add sync job", "error", err)
		os.Exit(1)
	}

	// The backup job itself is registered by api.NewRouter (via the system
	// handler's InitBackupJob), not here — that keeps the handler's
	// BackupJobID bookkeeping in sync with what's actually scheduled, so a
	// later PUT /schedule can remove/replace it instead of stacking a
	// duplicate job alongside this startup registration.
	jobRunner.Start(ctx)

	// Build HTTP router
	routerCfg := api.Config{
		Store:      s,
		JWT:        jwtSvc,
		Config:     cfg,
		SyncEngine: syncEngine,
		DocEngine:  docEngine,
		WSHub:      wsHub,
		Scheduler:  jobRunner,
		BackupDir:  backupDir,
		AIRegistry: aiRegistry,
	}
	if cfg.Server.Embed {
		spaFiles, err := fs.Sub(web.DistFS, "dist")
		if err != nil {
			logger.Error("Failed to load embedded SPA files", "error", err)
			os.Exit(1)
		}
		routerCfg.SPAFiles = spaFiles
	}
	router := api.NewRouter(routerCfg)

	// Start HTTP server
	srv := &http.Server{
		Addr:              cfg.Server.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second, // slowloris mitigation on an unauthenticated listener
		ReadTimeout:       cfg.Server.ReadTimeoutDuration(),
		WriteTimeout:      cfg.Server.WriteTimeoutDuration(),
	}

	go func() {
		logger.Info("HTTP server listening", "addr", cfg.Server.Addr())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// Start snoozed alert expiration goroutine
	go runAlertExpirer(ctx, s, notifDispatcher, logger)
	go notifications.RunDeliveryRetries(ctx, notifDispatcher, logger)
	go store.RunDocLockSweep(ctx, s, wsHub, store.DocLockHeartbeat, logger)

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("Shutting down gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeoutDuration())
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}

	if err := s.Close(); err != nil {
		logger.Error("Failed to close store", "error", err)
	}

	logger.Info("Shutdown complete")
}

// runHealthcheck queries this server's own /api/health endpoint and exits 0 if
// healthy, 1 otherwise. Used as the Docker HEALTHCHECK command since distroless
// images ship no shell/curl to do this externally. Always terminates the process.
func runHealthcheck() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: failed to load config: %v\n", err)
		os.Exit(1)
	}

	host := cfg.Server.Host
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	url := fmt.Sprintf("http://%s:%d/api/health", host, cfg.Server.Port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "healthcheck: server returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var body struct {
		Healthy bool `json:"healthy"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: failed to parse response: %v\n", err)
		os.Exit(1)
	}

	if !body.Healthy {
		fmt.Fprintln(os.Stderr, "healthcheck: server reported unhealthy")
		os.Exit(1)
	}

	os.Exit(0)
}

// runAlertExpirer periodically un-snoozes expired alerts.
func runAlertExpirer(ctx context.Context, s *store.Store, _ *notifications.Dispatcher, logger *slog.Logger) {
	logger.Info("Alert expirer started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			alerts, err := s.GetExpiredSnoozedAlerts(ctx)
			if err != nil {
				logger.Error("Failed to get expired snoozed alerts", "error", err)
			}
			for _, a := range alerts {
				if err := s.UpdateAlertStatus(ctx, a.ID, "pending", ""); err != nil {
					logger.Error("Failed to un-snooze alert", "error", err, "alert_id", a.ID)
				}
			}
			if len(alerts) > 0 {
				logger.Info("Un-snoozed expired alerts", "count", len(alerts))
			}
		}
		// Check every 60 seconds
		select {
		case <-ctx.Done():
			return
		case <-time.After(60 * time.Second):
		}
	}
}

func newLogger(cfg config.LogSettings) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

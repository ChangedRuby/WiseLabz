// Package main is the entry point for the WiseLabz server.
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api"
	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/config"
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
	if len(os.Args) > 1 && os.Args[1] == "config" {
		os.Exit(runConfigCommand(os.Args[2:], os.Stdout, os.Stderr))
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

	if err := cfg.Validate(); err != nil {
		logger.Error("Invalid configuration (run `server config validate` to check)", "error", err)
		os.Exit(1)
	}

	// Open database
	db, err := store.OpenDB(cfg.DB.Driver, cfg.DB.DSN, store.PoolConfig{
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime(),
		ConnMaxIdleTime: cfg.DB.ConnMaxIdleTime(),
	})
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	// Run migrations
	if err := store.RunMigrations(db, cfg.DB.Driver, logger); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize store
	s := store.New(db, cfg.DB.Driver)

	// Create root context that cancels on interrupt. This only signals that
	// shutdown should begin; the lifecycle manager below owns the separate
	// context that actually stops the background goroutines, so it can do so
	// in order instead of everything canceling out at once.
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
	// can broadcast progress events). Started by the lifecycle manager below.
	wsHub := ws.NewHub(splitOrigins(cfg.Server.Origin)...)

	// Initialize notification dispatcher (must precede sync engine so it can
	// notify on alert creation)
	notifDispatcher := notifications.NewDispatcher(s, wsHub)
	notifDispatcher.SetEncryptionKey(cfg.Encryption.Key)

	// Initialize engines
	qualityChecker := quality.NewChecker(s, wsHub, notifDispatcher,
		quality.RotationConfig{MaxAgeDays: cfg.Rotation.MaxAgeDays, WarnDays: cfg.Rotation.WarnDays})
	syncEngine := sync.NewEngine(s, wsHub, notifDispatcher, qualityChecker, cfg.Encryption.Key)
	docEngine := doc.NewEngine(s)
	syncEngine.SetDocRegenerator(docEngine)

	aiRegistry := ai.NewRegistry()
	ai.RegisterOpenAICompatible(aiRegistry)
	ai.RegisterClaude(aiRegistry)

	embedRegistry := ai.NewEmbedRegistry()
	ai.RegisterOllamaEmbedder(embedRegistry)
	ai.RegisterOpenAIEmbedder(embedRegistry)

	// Determine backup directory: use configured value, or compute from DB DSN
	backupDir := cfg.Backup.Dir
	if backupDir == "" {
		// Default: ./data/backups, or extract from SQLite path if configured
		backupDir = "./data/backups"
	}

	syncEngine.SetBaseContext(ctx)

	// Start scheduler for quality, sync, digest, and alert expiry jobs. The
	// backup and retention jobs (and the schedule/settings seeding that used
	// to live here) are registered by api.NewRouter, via the system
	// handler's InitBackupJob/InitRetentionJob — see those for why.
	jobRunner := scheduler.New(logger)
	if _, err := jobRunner.AddJob("quality", cfg.Quality.CronExpr, func(jobCtx context.Context) {
		quality.RunStaleSweepOnce(jobCtx, s, wsHub, notifDispatcher, logger)
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
	if _, err := jobRunner.AddJob("digest", "0 * * * *", func(jobCtx context.Context) {
		notifDispatcher.RunDigestSweep(jobCtx, time.Now().UTC(), logger)
	}); err != nil {
		logger.Error("Failed to add digest job", "error", err)
		os.Exit(1)
	}
	if _, err := jobRunner.AddJob("alertExpirer", "0 * * * * *", func(jobCtx context.Context) {
		expireAlertsOnce(jobCtx, s, notifDispatcher, logger)
	}); err != nil {
		logger.Error("Failed to add alert expirer job", "error", err)
		os.Exit(1)
	}

	// Unlike the backup job itself (registered by api.NewRouter via
	// InitBackupJob, since its schedule is operator-configurable through
	// PUT /schedule), the verify job has no persisted schedule of its own —
	// it just registers here on a fixed cadence, same as "quality"/"digest"
	// above.
	if _, err := jobRunner.AddJob("backup-verify", backup.DefaultVerifyCronExpr, func(jobCtx context.Context) {
		backup.RunVerifyOnce(jobCtx, backupDir, logger)
	}); err != nil {
		logger.Error("Failed to add backup verify job", "error", err)
		os.Exit(1)
	}

	// Build HTTP router
	readyState := &syshandler.ReadyState{}
	routerCfg := api.Config{
		Store:          s,
		JWT:            jwtSvc,
		Config:         cfg,
		SyncEngine:     syncEngine,
		DocEngine:      docEngine,
		WSHub:          wsHub,
		Scheduler:      jobRunner,
		BackupDir:      backupDir,
		AIRegistry:     aiRegistry,
		EmbedRegistry:  embedRegistry,
		QualityChecker: qualityChecker,
		Ready:          readyState,
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

	// The lifecycle manager owns every long-running goroutine (HTTP server,
	// WS hub, scheduler, delivery retries, doc lock sweep) under one
	// errgroup, and runs the ordered stop on shutdown: mark not-ready ->
	// drain HTTP/WS -> stop scheduler -> wait for remaining goroutines ->
	// wait for in-flight dispatch goroutines -> close the DB last.
	lifecycle := newLifecycleManager(lifecycleDeps{
		Logger:          logger,
		HTTPServer:      srv,
		WSHub:           wsHub,
		Scheduler:       jobRunner,
		Dispatcher:      notifDispatcher,
		Store:           s,
		Ready:           readyState,
		ShutdownTimeout: cfg.Server.ShutdownTimeoutDuration(),
	})
	lifecycle.Start()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("Shutting down gracefully")

	if err := lifecycle.Shutdown(); err != nil {
		logger.Error("Failed to close store", "error", err)
	}

	logger.Info("Shutdown complete")
}

// runHealthcheck queries this server's own /readyz endpoint and exits 0 for a
// ready server. Used as the Docker HEALTHCHECK command since distroless images
// ship no shell/curl to do this externally. Always terminates the process.
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
	url := fmt.Sprintf("http://%s:%d/readyz", host, cfg.Server.Port)

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

	os.Exit(0)
}

// expireAlertsOnce un-snoozes every alert whose snooze has expired and
// notifies affected users via the dispatcher, same fanout as a newly created
// alert (it's actionable again). Runs as a scheduler job instead of its own
// hand-rolled select-then-sleep loop.
func expireAlertsOnce(ctx context.Context, s *store.Store, d *notifications.Dispatcher, logger *slog.Logger) {
	expired, err := s.GetExpiredSnoozedAlerts(ctx)
	if err != nil {
		logger.Error("get expired snoozed alerts", "error", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	n, err := s.UnsnoozeExpiredAlerts(ctx)
	if err != nil {
		logger.Error("unsnooze expired alerts", "error", err)
		return
	}
	if n > 0 {
		logger.Info("un-snoozed expired alerts", "count", n)
		d.NotifyAlertsCreated(ctx, expired)
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

// splitOrigins turns the comma-separated server.origin setting into a list.
func splitOrigins(raw string) []string {
	var out []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, o)
		}
	}
	return out
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/msrovani/myclaw/internal/config"
	"github.com/msrovani/myclaw/internal/db"
	"github.com/msrovani/myclaw/internal/memory"
	"github.com/msrovani/myclaw/internal/observability"
	"github.com/msrovani/myclaw/internal/providers"
	"github.com/msrovani/myclaw/internal/router"
	"github.com/msrovani/myclaw/internal/skills"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := observability.NewLogger(cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(logger)

	slog.Info("starting xxxclaw",
		"version", Version,
		"log_level", cfg.LogLevel,
		"http_addr", cfg.HTTPAddr,
	)

	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","version":"%s","time":"%s"}`, Version, time.Now().UTC().Format(time.RFC3339))
	})

	// Readiness endpoint
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready"}`)
	})

	// pprof endpoints (protected: only in dev/admin mode)
	if cfg.PprofEnabled {
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		go func() {
			pprofAddr := cfg.PprofAddr
			slog.Info("pprof server starting", "addr", pprofAddr)
			if err := http.ListenAndServe(pprofAddr, pprofMux); err != nil {
				slog.Error("pprof server failed", "error", err)
			}
		}()
	}

	// --- XXXCLAW Core Subsystems ---
	slog.Info("initializing core subsystems...")

	dbMgr := db.NewManager(db.Config{
		BaseDataDir: "./data",
		BusyTimeout: 5000,
		MaxReaders:  4,
	})

	memEngine := memory.NewEngine(dbMgr)
	_ = memEngine // Available for future controllers

	economy := router.NewEconomy(dbMgr)
	llmRouter := router.NewRouter([]providers.Provider{}, economy)
	_ = llmRouter // Available for future controllers

	skillRegistry := skills.NewRegistry()
	_ = skillRegistry // Available for future controllers
	// -------------------------------

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("http server starting", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received, draining...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown error", "error", err)
	}

	// Shutdown core subsystems safely
	slog.Info("closing database connections...")
	dbMgr.CloseAll()

	slog.Info("xxxclaw stopped")
}

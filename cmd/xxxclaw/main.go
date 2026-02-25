package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/msrovani/myclaw/internal/config"
	"github.com/msrovani/myclaw/internal/container"
	"github.com/msrovani/myclaw/internal/dashboard"
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
		"env", cfg.Env,
		"http_addr", cfg.HTTPAddr,
	)

	// --- Initialization via Type-Safe Container ---
	ct := container.New()

	// 1. Database Manager
	dbMgr := db.NewManager(db.Config{
		BaseDataDir: "./data",
		BusyTimeout: cfg.DBBusyTimeout,
		MaxReaders:  4,
		Env:         cfg.Env,
	})
	ct.Register("db", dbMgr)

	// 2. LLM Providers & Router
	ollama := providers.NewOllamaProvider(cfg.OllamaURL)
	ct.Register("provider.ollama", ollama)

	economy := router.NewEconomy(dbMgr)
	ct.Register("economy", economy)

	llmRouter := router.NewRouter([]providers.Provider{ollama}, economy)
	ct.Register("router", llmRouter)

	// 3. Memory Engine (Injected with Ollama for embeddings)
	memEngine := memory.NewEngine(dbMgr, ollama)
	ct.Register("memory", memEngine)

	// 4. Skills Registry
	skillRegistry := skills.NewRegistry()
	if err := skills.LoadAgentDir(".agent", skillRegistry); err != nil {
		slog.Warn("skills: error mounting .agent directory", "error", err)
	}
	ct.Register("skills", skillRegistry)

	// --- HTTP Server Setup ---
	mux := http.NewServeMux()

	// Health/Ready endpoints
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, Version)
	})
	dashboard.NewServer(mux, memEngine, economy, llmRouter, skillRegistry)

	// pprof (only if enabled)
	if cfg.PprofEnabled {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			http.Redirect(w, r, "/debug/pprof/", http.StatusFound)
		})
		go func() {
			slog.Info("pprof server starting", "addr", cfg.PprofAddr)
			_ = http.ListenAndServe(cfg.PprofAddr, nil) // uses DefaultServeMux
		}()
	}

	// --- Shutdown Logic ---
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	// Critical: Shutdown all services in LIFO order via container
	if err := ct.Shutdown(shutdownCtx); err != nil {
		slog.Error("container shutdown error", "error", err)
	}

	slog.Info("xxxclaw stopped")
}

package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any env vars that might interfere
	envVars := []string{
		"XXXCLAW_HTTP_ADDR", "XXXCLAW_LOG_LEVEL", "XXXCLAW_LOG_FORMAT",
		"XXXCLAW_PPROF_ENABLED", "XXXCLAW_PPROF_ADDR", "XXXCLAW_DB_PATH",
		"XXXCLAW_ENV", "XXXCLAW_VECTOR_ENABLED", "XXXCLAW_VECTOR_DIM",
	}
	for _, k := range envVars {
		os.Unsetenv(k)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.LogFormat != "json" {
		t.Errorf("LogFormat = %q, want %q", cfg.LogFormat, "json")
	}
	if !cfg.PprofEnabled {
		t.Error("PprofEnabled = false, want true")
	}
	if cfg.Env != "dev" {
		t.Errorf("Env = %q, want %q", cfg.Env, "dev")
	}
	if cfg.VectorDim != 384 {
		t.Errorf("VectorDim = %d, want %d", cfg.VectorDim, 384)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("XXXCLAW_HTTP_ADDR", ":9090")
	t.Setenv("XXXCLAW_LOG_LEVEL", "debug")
	t.Setenv("XXXCLAW_ENV", "prod")
	t.Setenv("XXXCLAW_VECTOR_DIM", "768")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.HTTPAddr != ":9090" {
		t.Errorf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":9090")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.Env != "prod" {
		t.Errorf("Env = %q, want %q", cfg.Env, "prod")
	}
	if cfg.VectorDim != 768 {
		t.Errorf("VectorDim = %d, want %d", cfg.VectorDim, 768)
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	t.Setenv("XXXCLAW_LOG_LEVEL", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid log level")
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	t.Setenv("XXXCLAW_LOG_LEVEL", "info")
	t.Setenv("XXXCLAW_ENV", "staging")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid env")
	}
}

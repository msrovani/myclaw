package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration.
type Config struct {
	// HTTP server address
	HTTPAddr string

	// Logging
	LogLevel  string
	LogFormat string // "json" or "text"

	// pprof
	PprofEnabled bool
	PprofAddr    string

	// SQLite
	DBPath        string
	DBBusyTimeout int // milliseconds

	// Vector search
	VectorEnabled bool
	VectorDim     int

	// Environment
	Env string // "dev", "prod", "edge"

	// LLM Providers
	OllamaURL    string
	GeminiAPIKey string
	ClaudeAPIKey string

	// Token Economy
	DefaultBudgetPerSession int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:                envOrDefault("XXXCLAW_HTTP_ADDR", ":8080"),
		LogLevel:                envOrDefault("XXXCLAW_LOG_LEVEL", "info"),
		LogFormat:               envOrDefault("XXXCLAW_LOG_FORMAT", "json"),
		PprofEnabled:            envBoolOrDefault("XXXCLAW_PPROF_ENABLED", true),
		PprofAddr:               envOrDefault("XXXCLAW_PPROF_ADDR", ":6060"),
		DBPath:                  envOrDefault("XXXCLAW_DB_PATH", "data/xxxclaw.db"),
		DBBusyTimeout:           envIntOrDefault("XXXCLAW_DB_BUSY_TIMEOUT", 5000),
		VectorEnabled:           envBoolOrDefault("XXXCLAW_VECTOR_ENABLED", true),
		VectorDim:               envIntOrDefault("XXXCLAW_VECTOR_DIM", 384),
		Env:                     envOrDefault("XXXCLAW_ENV", "dev"),
		OllamaURL:               envOrDefault("XXXCLAW_OLLAMA_URL", "http://localhost:11434"),
		GeminiAPIKey:            os.Getenv("XXXCLAW_GEMINI_API_KEY"),
		ClaudeAPIKey:            os.Getenv("XXXCLAW_CLAUDE_API_KEY"),
		DefaultBudgetPerSession: envIntOrDefault("XXXCLAW_DEFAULT_BUDGET", 10000),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[strings.ToLower(c.LogLevel)] {
		return fmt.Errorf("invalid log level: %s (must be debug/info/warn/error)", c.LogLevel)
	}

	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[strings.ToLower(c.LogFormat)] {
		return fmt.Errorf("invalid log format: %s (must be json/text)", c.LogFormat)
	}

	validEnvs := map[string]bool{"dev": true, "prod": true, "edge": true}
	if !validEnvs[strings.ToLower(c.Env)] {
		return fmt.Errorf("invalid env: %s (must be dev/prod/edge)", c.Env)
	}

	if c.VectorDim <= 0 {
		return fmt.Errorf("vector dimension must be positive, got %d", c.VectorDim)
	}

	return nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envBoolOrDefault(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

func envIntOrDefault(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}

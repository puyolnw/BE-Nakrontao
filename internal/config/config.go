// Package config reads runtime settings from the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Addr          string
	DatabaseURL   string
	StaticDir     string // built frontend; empty disables static serving
	SecureCookies bool
	IPHashSalt    string

	AnthropicKey   string
	AnthropicModel string
	EmbedURL       string // OpenAI-compatible /v1/embeddings endpoint
	EmbedKey       string
	EmbedModel     string

	TrustedProxy  string // CIDRs of the reverse proxy, comma separated
	RetentionDays int    // chat logs and page views older than this are deleted
	ChatPerMinute int    // per client
	ChatDailyCap  int    // global LLM calls per day, cost control

	BootstrapAdminEmail    string
	BootstrapAdminPassword string
}

func Load() (Config, error) {
	c := Config{
		Addr:                   listenAddr(),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		StaticDir:              os.Getenv("STATIC_DIR"),
		SecureCookies:          env("SECURE_COOKIES", "true") == "true",
		IPHashSalt:             os.Getenv("IP_HASH_SALT"),
		AnthropicKey:           os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicModel:         env("ANTHROPIC_MODEL", "claude-sonnet-5-5"),
		EmbedURL:               os.Getenv("EMBED_URL"),
		EmbedKey:               os.Getenv("EMBED_API_KEY"),
		EmbedModel:             env("EMBED_MODEL", "text-embedding-3-small"),
		TrustedProxy:           os.Getenv("TRUSTED_PROXY"),
		RetentionDays:          envInt("RETENTION_DAYS", 180),
		ChatPerMinute:          envInt("CHAT_PER_MINUTE", 10),
		ChatDailyCap:           envInt("CHAT_DAILY_CAP", 2000),
		BootstrapAdminEmail:    os.Getenv("BOOTSTRAP_ADMIN_EMAIL"),
		BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
	}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL is required")
	}
	if len(c.IPHashSalt) < 16 {
		return c, fmt.Errorf("IP_HASH_SALT must be at least 16 characters")
	}
	return c, nil
}

// listenAddr honours ADDR, then PORT (set by hosts such as Railway), then :8080.
func listenAddr() string {
	if a := os.Getenv("ADDR"); a != "" {
		return a
	}
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(k)); err == nil && v > 0 {
		return v
	}
	return def
}

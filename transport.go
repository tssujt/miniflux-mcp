package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

const (
	transportStdio          = "stdio"
	transportStreamableHTTP = "streamable-http"
	defaultHTTPAddr         = ":8080"
	defaultHTTPPath         = "/mcp"
)

type transportConfig struct {
	Transport string
	HTTPAddr  string
	HTTPPath  string
	AuthToken string
}

func loadTransportConfig() (transportConfig, error) {
	cfg := transportConfig{
		Transport: envOrDefault("MCP_TRANSPORT", transportStdio),
		HTTPAddr:  envOrDefault("MCP_HTTP_ADDR", defaultHTTPAddr),
		HTTPPath:  envOrDefault("MCP_HTTP_PATH", defaultHTTPPath),
		AuthToken: os.Getenv("MCP_AUTH_TOKEN"),
	}

	switch cfg.Transport {
	case transportStdio:
		return cfg, nil
	case transportStreamableHTTP:
		if cfg.AuthToken == "" {
			return transportConfig{}, fmt.Errorf("MCP_AUTH_TOKEN is required when MCP_TRANSPORT=%s", transportStreamableHTTP)
		}
		if !strings.HasPrefix(cfg.HTTPPath, "/") || cfg.HTTPPath == "/" {
			return transportConfig{}, fmt.Errorf("MCP_HTTP_PATH must start with / and cannot be /")
		}
		if cfg.HTTPPath == "/healthz" {
			return transportConfig{}, fmt.Errorf("MCP_HTTP_PATH cannot be /healthz")
		}
		return cfg, nil
	default:
		return transportConfig{}, fmt.Errorf("unsupported MCP_TRANSPORT %q (supported: %s, %s)", cfg.Transport, transportStdio, transportStreamableHTTP)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func serveMCP(mcpServer *server.MCPServer, cfg transportConfig) error {
	switch cfg.Transport {
	case transportStdio:
		return server.ServeStdio(mcpServer)
	case transportStreamableHTTP:
		return serveStreamableHTTP(mcpServer, cfg)
	default:
		return fmt.Errorf("unsupported MCP transport %q", cfg.Transport)
	}
}

func serveStreamableHTTP(mcpServer *server.MCPServer, cfg transportConfig) error {
	mcpHandler := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithStateLess(true),
	)

	mux := http.NewServeMux()
	mux.Handle(cfg.HTTPPath, requireBearerToken(cfg.AuthToken, mcpHandler))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	return httpServer.ListenAndServe()
}

func requireBearerToken(token string, next http.Handler) http.Handler {
	expectedTokenHash := sha256.Sum256([]byte(token))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme, providedToken, ok := strings.Cut(r.Header.Get("Authorization"), " ")
		providedTokenHash := sha256.Sum256([]byte(providedToken))
		if !ok || !strings.EqualFold(scheme, "Bearer") || subtle.ConstantTimeCompare(providedTokenHash[:], expectedTokenHash[:]) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

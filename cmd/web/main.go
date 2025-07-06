package main

import (
	"errors"
	"finance/internal/config"
	"finance/internal/web"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/guilhermebr/gox/logger"
)

// Injected on build time by ldflags.
var (
	BuildCommit = "undefined"
	BuildTime   = "undefined"
)

func main() {
	var cfg config.Config
	if err := cfg.Load(""); err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}

	// Logger
	log, err := logger.NewLogger("")
	if err != nil {
		panic(fmt.Errorf("creating logger: %w", err))
	}

	log = log.With(
		slog.String("environment", cfg.Environment),
		slog.String("build_commit", BuildCommit),
		slog.String("build_time", BuildTime),
		slog.Int("go_max_procs", runtime.GOMAXPROCS(0)),
		slog.Int("runtime_num_cpu", runtime.NumCPU()),
	)

	// API base URL configuration
	apiBaseURL := cfg.Web.ApiBaseURL

	// Web handlers - now only needs API base URL
	webHandlers := web.NewHandlers(apiBaseURL)

	// Server
	server := http.Server{
		Handler:           webHandlers.Router(),
		Addr:              cfg.Web.Address,
		ReadHeaderTimeout: 60 * time.Second,
	}

	log.Info("web server started",
		slog.String("address", server.Addr),
		slog.String("api_base_url", apiBaseURL),
	)

	if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		log.Error("failed to listen and serve web server",
			slog.String("error", serverErr.Error()),
		)
	}
}

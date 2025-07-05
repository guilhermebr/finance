package main

import (
	"context"
	"errors"
	"finance/domain/finance"
	"finance/internal/config"
	"finance/internal/repository/pg"
	"finance/internal/web"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/guilhermebr/gox/logger"
	"github.com/guilhermebr/gox/postgres"
)

// Injected on build time by ldflags.
var (
	BuildCommit = "undefined"
	BuildTime   = "undefined"
)

func main() {
	ctx := context.Background()

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

	// Database connection
	conn, err := postgres.New(ctx, "")
	if err != nil {
		log.Error("failed to setup postgres",
			slog.String("error", err.Error()),
		)
		return
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		log.Error("failed to reach postgres",
			slog.String("error", err.Error()),
		)
		return
	}

	// Finance repositories
	accountRepo := pg.NewAccountRepository(conn)
	categoryRepo := pg.NewCategoryRepository(conn)
	transactionRepo := pg.NewTransactionRepository(conn)
	balanceRepo := pg.NewBalanceRepository(conn)

	// Finance use cases
	accountUseCase := finance.NewAccountUseCase(accountRepo, balanceRepo)
	categoryUseCase := finance.NewCategoryUseCase(categoryRepo)
	transactionUseCase := finance.NewTransactionUseCase(transactionRepo, accountRepo, categoryRepo, balanceRepo)
	balanceUseCase := finance.NewBalanceUseCase(balanceRepo, accountRepo)

	// Web handlers
	webHandlers := web.NewHandlers(accountUseCase, categoryUseCase, transactionUseCase, balanceUseCase)

	// Server
	server := http.Server{
		Handler:           webHandlers.Router(),
		Addr:              ":8080", // Different port from API
		ReadHeaderTimeout: 60 * time.Second,
	}

	log.Info("web server started",
		slog.String("address", server.Addr),
	)

	if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		log.Error("failed to listen and serve web server",
			slog.String("error", serverErr.Error()),
		)
	}
}

package main

import (
	"context"
	"errors"
	"finance/domain/example"
	"finance/domain/finance"
	"finance/internal/api"
	v1 "finance/internal/api/v1"
	"finance/internal/config"
	"finance/internal/repository/pg"
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

	// Legacy repository for examples
	exampleRepo := pg.NewRepository(conn)

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

	// API Handlers V1
	// ------------------------------------------
	apiV1 := v1.ApiHandlers{
		ExampleUseCase:     example.New(exampleRepo),
		AccountUseCase:     accountUseCase,
		CategoryUseCase:    categoryUseCase,
		TransactionUseCase: transactionUseCase,
		BalanceUseCase:     balanceUseCase,
	}

	router := api.Router()
	apiV1.Routes(router)

	// SERVER
	// ------------------------------------------
	server := http.Server{
		Handler:           router,
		Addr:              cfg.Service.Address,
		ReadHeaderTimeout: 60 * time.Second,
	}
	log.Info("server started",
		slog.String("address", server.Addr),
	)

	if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		log.Error("failed to listen and serve server",
			slog.String("error", serverErr.Error()),
		)
	}
}

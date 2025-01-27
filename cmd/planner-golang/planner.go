package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"planner-golang/internal/api"
	"planner-golang/internal/api/spec"
	"planner-golang/internal/mailer/mailpit"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(" bye! ")
}

func run(ctx context.Context) error {
	// USED ONLY FOR LOCAL TESTS
	//needed to load dotenv definitions
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	//end dotenv

	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("planner")
	defer func() { _ = logger.Sync() }() //anonimized function to hide hints

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s  password=%s host=%s port=%s dbname=%s",
		os.Getenv("PLANNER_DATABASE_USER"),
		os.Getenv("PLANNER_DATABASE_PASSWORD"),
		os.Getenv("PLANNER_DATABASE_HOST"),
		os.Getenv("PLANNER_DATABASE_PORT"),
		os.Getenv("PLANNER_DATABASE_NAME"),
	),
	)

	if err != nil {
		return err
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	si := api.NewAPI(pool, logger, mailpit.NewMailPit(pool))
	r := chi.NewMux()
	r.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))
	r.Mount("/", spec.Handler(&si))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("falied to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

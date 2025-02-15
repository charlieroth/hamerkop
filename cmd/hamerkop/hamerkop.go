package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charlieroth/hamerkop/internal/app"
	"github.com/charlieroth/hamerkop/internal/config"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	ctx := context.Background()
	if err := run(ctx, &log); err != nil {
		log.Fatal().Err(err).Msg("failed to boot hamerkop")
	}
}

func run(ctx context.Context, log *zerolog.Logger) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	log.Info().Msg("⚡️ booting hamerkop...")
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	log.Info().Msg("✅ loaded relay configuration")

	log.Info().Msg("✅ initializing relay...")
	hamerkop := app.NewApp(cfg, log)
	err = hamerkop.Init()
	if err != nil {
		return err
	}
	log.Info().Msg("✅ relay initialized")

	addr := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)
	log.Info().Msgf("🚀 hamerkop listening at %s", addr)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// Start relay
	go func() {
		log.Info().Msg("🔄 starting relay...")
		serverErrors <- http.ListenAndServe(addr, hamerkop.Relay)
	}()

	// Handle shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("❌ server error: %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("👋 shutting down hamerkop... (signal: %s)", sig)
		defer log.Info().Msg("👋 hamerkop shutdown complete")

		ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
		defer cancel()
		if err := hamerkop.Shutdown(ctx); err != nil {
			hamerkop.Close()
			return fmt.Errorf("failed to shutdown relay: %w", err)
		}
	}

	return nil
}

package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charlieroth/hamerkop/internal/app"
	"github.com/charlieroth/hamerkop/internal/config"
	"github.com/charlieroth/hamerkop/internal/utils"
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
	log.Info().Msg("⚡️ booting hamerkop...")
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	log.Info().Msg("✅ loaded relay configuration")

	pubkey, err := utils.NpubToPubkey(cfg.Relay.Pubkey)
	if err != nil {
		return err
	}

	log.Info().Msg("✅ initializing relay...")
	hamerkop := app.NewApp(cfg)
	err = hamerkop.Init()
	if err != nil {
		return err
	}
	log.Info().Msg("✅ relay initialized")

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.Header.Get("Upgrade") == "websocket" {
			// If request is for websocket, serve the relay
			hamerkop.Relay.ServeHTTP(w, r)
		} else if r.Method == "GET" && r.Header.Get("Accept") == "application/nostr+json" {
			hamerkop.Relay.ServeHTTP(w, r)
		} else {
			// Otherwise, serve the index.html template
			relaySupportedNIPsStrings := []string{}
			for _, nip := range cfg.Relay.SupportedNIPs {
				relaySupportedNIPsStrings = append(relaySupportedNIPsStrings, fmt.Sprintf("%v", nip))
			}
			relaySupportedNIPs := strings.Join(relaySupportedNIPsStrings, ", ")

			relayAllowedKindsStrings := []string{}
			for _, kind := range cfg.AllowedKinds {
				relayAllowedKindsStrings = append(relayAllowedKindsStrings, fmt.Sprintf("%v", kind))
			}
			relayAllowedKinds := strings.Join(relayAllowedKindsStrings, ", ")

			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			data := struct {
				RelayName          string
				RelayPubkey        string
				RelayDescription   string
				RelayURL           string
				RelaySoftware      string
				RelayVersion       string
				RelaySupportedNIPs string
				RelayAllowedKinds  string
			}{
				RelayName:          cfg.Relay.Name,
				RelayPubkey:        pubkey,
				RelayDescription:   cfg.Relay.Description,
				RelayURL:           "https://hamerkop.charlieroth.me",
				RelaySoftware:      cfg.Relay.Software,
				RelayVersion:       cfg.Relay.Version,
				RelaySupportedNIPs: relaySupportedNIPs,
				RelayAllowedKinds:  relayAllowedKinds,
			}
			err := tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	addr := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)
	log.Info().Msgf("🚀 hamerkop listening at %s", addr)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msg("🔄 starting relay...")
		serverErrors <- http.ListenAndServe(addr, nil)
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
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

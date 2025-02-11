package app

import (
	"context"

	"github.com/charlieroth/hamerkop/internal/config"
	"github.com/charlieroth/hamerkop/internal/store"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/policies"
)

type App struct {
	cfg   *config.Config
	Relay *khatru.Relay
	Store store.Store
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg:   cfg,
		Relay: khatru.NewRelay(),
		Store: store.NewStore("hamerkop.db"),
	}
}

func (a *App) Init() error {
	// Set the relay information for NIP-11 support
	a.Relay.Info.Name = a.cfg.Relay.Name
	a.Relay.Info.Description = a.cfg.Relay.Description
	a.Relay.Info.Version = a.cfg.Relay.Version
	a.Relay.Info.Software = a.cfg.Relay.Software
	a.Relay.Info.Contact = a.cfg.Relay.Contact
	a.Relay.Info.Icon = a.cfg.Relay.Icon
	a.Relay.Info.PubKey = a.cfg.Relay.Pubkey
	a.Relay.Info.SupportedNIPs = a.cfg.Relay.SupportedNIPs
	a.Relay.ServiceURL = "https://hamerkop.charlieroth.me"

	// Set up policies
	if !a.cfg.Limits.AllowEmptyFilters {
		a.Relay.RejectFilter = append(a.Relay.RejectFilter, policies.NoEmptyFilters)
	}

	if !a.cfg.Limits.AllowComplexFilters {
		a.Relay.RejectFilter = append(a.Relay.RejectFilter, policies.NoComplexFilters)
	}

	a.Relay.RejectEvent = append(
		a.Relay.RejectEvent,
		policies.PreventLargeTags(100),
		policies.RestrictToSpecifiedKinds(true, a.cfg.AllowedKinds...),
	)

	// Set up event handlers
	a.Relay.StoreEvent = append(a.Relay.StoreEvent, a.Store.SaveEvent)
	a.Relay.QueryEvents = append(a.Relay.QueryEvents, a.Store.QueryEvents)
	a.Relay.DeleteEvent = append(a.Relay.DeleteEvent, a.Store.DeleteEvent)
	a.Relay.CountEvents = append(a.Relay.CountEvents, a.Store.CountEvents)
	a.Relay.ReplaceEvent = append(a.Relay.ReplaceEvent, a.Store.ReplaceEvent)

	// Initialize event store
	if err := a.Store.Init(); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return nil
}

func (a *App) Close() {
	a.Store.Close()
}

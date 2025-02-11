package store

import (
	"context"

	"github.com/fiatjaf/eventstore/badger"
	"github.com/nbd-wtf/go-nostr"
)

type Store interface {
	Init() error
	Close()
	CountEvents(ctx context.Context, filter nostr.Filter) (int64, error)
	DeleteEvent(ctx context.Context, event *nostr.Event) error
	QueryEvents(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error)
	SaveEvent(ctx context.Context, event *nostr.Event) error
	ReplaceEvent(ctx context.Context, event *nostr.Event) error
	Serial() []byte
}

func NewStore(path string) Store {
	return &badger.BadgerBackend{
		Path: path,
	}
}

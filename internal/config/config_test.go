package config

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadConfig(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Errorf("Failed to load environment variables: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}

	if config == nil {
		t.Errorf("Config is nil")
	}

	if config.Port < 0 || config.Port > 65535 {
		t.Errorf("Invalid port number: %d", config.Port)
	}

	if config.Hostname == "" {
		t.Errorf("Hostname is empty")
	}

	if len(config.AllowedKinds) == 0 {
		t.Errorf("AllowedKinds is empty")
	}

	if config.Relay.Name == "" {
		t.Errorf("Relay name is empty")
	}

	if config.Relay.Description == "" {
		t.Errorf("Relay description is empty")
	}

	if config.Relay.Icon == "" {
		t.Errorf("Relay icon is empty")
	}

	if config.Relay.Pubkey == "" {
		t.Errorf("Relay pubkey is empty")
	}

	if config.Relay.Contact == "" {
		t.Errorf("Relay contact is empty")
	}

	if len(config.Relay.SupportedNIPs) == 0 {
		t.Errorf("Relay supported NIPs is empty")
	}

	if config.Relay.Software == "" {
		t.Errorf("Relay software is empty")
	}

	if config.Relay.Version == "" {
		t.Errorf("Relay version is empty")
	}
}

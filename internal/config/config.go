package config

import (
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"gopkg.in/yaml.v3"
)

type RelayConfig struct {
	Name          string `json:"name" yaml:"name"`
	Description   string `json:"description" yaml:"description"`
	Banner        string `json:"banner" yaml:"banner"`
	Icon          string `json:"icon" yaml:"icon"`
	Pubkey        string `json:"pubkey" yaml:"pubkey"`
	Contact       string `json:"contact" yaml:"contact"`
	SupportedNIPs []any  `json:"supported_nips" yaml:"supported_nips"`
	Software      string `json:"software" yaml:"software"`
	Version       string `json:"version" yaml:"version"`
}

type LimitsConfig struct {
	AllowEmptyFilters   bool `json:"allow_empty_filters" yaml:"allow_empty_filters"`
	AllowComplexFilters bool `json:"allow_complex_filters" yaml:"allow_complex_filters"`
}

type Config struct {
	Port         int          `json:"port" yaml:"port"`
	Hostname     string       `json:"hostname" yaml:"hostname"`
	Relay        RelayConfig  `json:"relay" yaml:"relay"`
	Limits       LimitsConfig `json:"limits" yaml:"limits"`
	AllowedKinds []uint16
}

func NewConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	cfg.AllowedKinds = []uint16{
		uint16(nostr.KindCommunityDefinition),
		uint16(nostr.KindCommunityList),
		uint16(nostr.KindCommunityPostApproval),
		uint16(nostr.KindBadgeAward),
		uint16(nostr.KindBadgeDefinition),
		uint16(nostr.KindArticle),
		uint16(nostr.KindDraftArticle),
		uint16(nostr.KindComment),
		uint16(nostr.KindBookmarkList),
		uint16(nostr.KindHighlights),
		uint16(nostr.KindMuteList),
		uint16(nostr.KindDeletion),
	}

	return &cfg, nil
}

func (c *Config) RelayURL() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

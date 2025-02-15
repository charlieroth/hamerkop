package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

type RelayConfig struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Icon          string `json:"icon"`
	Pubkey        string `json:"pubkey"`
	Contact       string `json:"contact"`
	SupportedNIPs []int  `json:"supported_nips"`
	Software      string `json:"software"`
	Version       string `json:"version"`
}

type LimitsConfig struct {
	AllowEmptyFilters   bool `json:"allow_empty_filters"`
	AllowComplexFilters bool `json:"allow_complex_filters"`
}

type Config struct {
	Port         int          `json:"port"`
	Hostname     string       `json:"hostname"`
	Relay        RelayConfig  `json:"relay"`
	Limits       LimitsConfig `json:"limits"`
	AllowedKinds []uint16
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("PORT is not a valid integer: %w", err)
	}

	hostname := os.Getenv("HOSTNAME")
	relayName := os.Getenv("RELAY_NAME")
	relayDescription := os.Getenv("RELAY_DESCRIPTION")
	relayIcon := os.Getenv("RELAY_ICON")
	relayPubkey := os.Getenv("RELAY_PUBKEY")
	relayContact := os.Getenv("RELAY_CONTACT")
	relaySoftware := os.Getenv("RELAY_SOFTWARE")
	relayVersion := os.Getenv("RELAY_VERSION")

	limitsAllowEmptyFilters := os.Getenv("LIMITS_ALLOW_EMPTY_FILTERS")
	limitsAllowEmptyFiltersBool, err := strconv.ParseBool(limitsAllowEmptyFilters)
	if err != nil {
		return nil, fmt.Errorf("LIMITS_ALLOW_EMPTY_FILTERS is not a valid boolean: %w", err)
	}

	limitsAllowComplexFilters := os.Getenv("LIMITS_ALLOW_COMPLEX_FILTERS")
	limitsAllowComplexFiltersBool, err := strconv.ParseBool(limitsAllowComplexFilters)
	if err != nil {
		return nil, fmt.Errorf("LIMITS_ALLOW_COMPLEX_FILTERS is not a valid boolean: %w", err)
	}

	supportedNIPs := os.Getenv("RELAY_SUPPORTED_NIPS")
	if supportedNIPs == "" {
		return nil, fmt.Errorf("RELAY_SUPPORTED_NIPS is not set")
	}

	supportedNIPsSlice := strings.Split(supportedNIPs, ",")
	supportedNIPsSliceInt := make([]int, len(supportedNIPsSlice))
	for i, nip := range supportedNIPsSlice {
		nipInt, err := strconv.ParseInt(nip, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("RELAY_SUPPORTED_NIPS is not a valid integer: %w", err)
		}

		supportedNIPsSliceInt[i] = int(nipInt)
	}

	cfg := Config{
		Port:     portInt,
		Hostname: hostname,
		Relay: RelayConfig{
			Name:          relayName,
			Description:   relayDescription,
			Icon:          relayIcon,
			Pubkey:        relayPubkey,
			Contact:       relayContact,
			SupportedNIPs: supportedNIPsSliceInt,
			Software:      relaySoftware,
			Version:       relayVersion,
		},
		Limits: LimitsConfig{
			AllowEmptyFilters:   limitsAllowEmptyFiltersBool,
			AllowComplexFilters: limitsAllowComplexFiltersBool,
		},
		AllowedKinds: []uint16{
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
		},
	}

	return &cfg, nil
}

func (c *Config) RelayURL() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

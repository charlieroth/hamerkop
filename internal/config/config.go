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
	Banner        string `json:"banner"`
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

func NewConfig() (*Config, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("PORT is not set")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("PORT is not a valid integer: %w", err)
	}

	hostname, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		return nil, fmt.Errorf("HOSTNAME is not set")
	}

	relayName, ok := os.LookupEnv("RELAY_NAME")
	if !ok {
		return nil, fmt.Errorf("RELAY_NAME is not set")
	}

	relayDescription, ok := os.LookupEnv("RELAY_DESCRIPTION")
	if !ok {
		return nil, fmt.Errorf("RELAY_DESCRIPTION is not set")
	}

	relayBanner, ok := os.LookupEnv("RELAY_BANNER")
	if !ok {
		return nil, fmt.Errorf("RELAY_BANNER is not set")
	}

	relayIcon, ok := os.LookupEnv("RELAY_ICON")
	if !ok {
		return nil, fmt.Errorf("RELAY_ICON is not set")
	}

	relayPubkey, ok := os.LookupEnv("RELAY_PUBKEY")
	if !ok {
		return nil, fmt.Errorf("RELAY_PUBKEY is not set")
	}

	relayContact, ok := os.LookupEnv("RELAY_CONTACT")
	if !ok {
		return nil, fmt.Errorf("RELAY_CONTACT is not set")
	}

	relaySoftware, ok := os.LookupEnv("RELAY_SOFTWARE")
	if !ok {
		return nil, fmt.Errorf("RELAY_SOFTWARE is not set")
	}

	relayVersion, ok := os.LookupEnv("RELAY_VERSION")
	if !ok {
		return nil, fmt.Errorf("RELAY_VERSION is not set")
	}

	limitsAllowEmptyFilters, ok := os.LookupEnv("LIMITS_ALLOW_EMPTY_FILTERS")
	if !ok {
		return nil, fmt.Errorf("LIMITS_ALLOW_EMPTY_FILTERS is not set")
	}

	limitsAllowEmptyFiltersBool, err := strconv.ParseBool(limitsAllowEmptyFilters)
	if err != nil {
		return nil, fmt.Errorf("LIMITS_ALLOW_EMPTY_FILTERS is not a valid boolean: %w", err)
	}

	limitsAllowComplexFilters, ok := os.LookupEnv("LIMITS_ALLOW_COMPLEX_FILTERS")
	if !ok {
		return nil, fmt.Errorf("LIMITS_ALLOW_COMPLEX_FILTERS is not set")
	}

	limitsAllowComplexFiltersBool, err := strconv.ParseBool(limitsAllowComplexFilters)
	if err != nil {
		return nil, fmt.Errorf("LIMITS_ALLOW_COMPLEX_FILTERS is not a valid boolean: %w", err)
	}

	supportedNIPs, ok := os.LookupEnv("RELAY_SUPPORTED_NIPS")
	if !ok {
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
			Banner:        relayBanner,
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

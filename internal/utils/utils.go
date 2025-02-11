package utils

import "github.com/nbd-wtf/go-nostr/nip19"

func NpubToPubkey(npub string) (string, error) {
	_, v, err := nip19.Decode(npub)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

package oauth2keystore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/trichner/toolbox/pkg/keyring"

	"golang.org/x/oauth2"
)

func NewKeyringTokenStore(serviceName string) *KeyringTokenStore {
	return &KeyringTokenStore{serviceName: serviceName}
}

type KeyringTokenStore struct {
	serviceName string
}

func (k *KeyringTokenStore) Get(key string) (*oauth2.Token, error) {
	ring, err := keyring.Open(k.serviceName)
	if err != nil {
		return nil, err
	}

	item, err := ring.Get(key)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot read %s token for %q: %w", k.serviceName, key, err)
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(item.Secret), &token)
	if err != nil {
		return nil, fmt.Errorf("invalid %s token for %q: %w", k.serviceName, key, err)
	}

	return &token, nil
}

func (k *KeyringTokenStore) Put(key string, token *oauth2.Token) error {
	ring, err := keyring.Open(k.serviceName)
	if err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	item := &keyring.Item{Secret: string(data)}
	if err := ring.Put(key, item); err != nil {
		return fmt.Errorf("failed to store token for %s: %w", k.serviceName, err)
	}
	return nil
}

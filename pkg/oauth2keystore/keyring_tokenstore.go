package oauth2keystore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
)

func NewKeyringTokenStore(serviceName string, description string) *KeyringTokenStore {
	return &KeyringTokenStore{serviceName: serviceName, description: description}
}

type KeyringTokenStore struct {
	serviceName string
	description string
}

func (k *KeyringTokenStore) Get(key string) (*oauth2.Token, error) {
	ring, err := keyring.Open(keyring.Config{ServiceName: k.serviceName})
	if err != nil {
		return nil, err
	}

	item, err := ring.Get(key)
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot read %s token for %q: %w", k.serviceName, key, err)
	}

	var token oauth2.Token
	err = json.Unmarshal(item.Data, &token)
	if err != nil {
		return nil, fmt.Errorf("invalid %s token for %q: %w", k.serviceName, key, err)
	}

	return &token, nil
}

func (k *KeyringTokenStore) Put(key string, token *oauth2.Token) error {
	ring, err := keyring.Open(keyring.Config{ServiceName: k.serviceName})
	if err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	item := keyring.Item{
		Key:         key,
		Data:        data,
		Description: k.description,
	}
	if err := ring.Set(item); err != nil {
		return fmt.Errorf("failed to store token for %s: %w", k.serviceName, err)
	}
	return nil
}

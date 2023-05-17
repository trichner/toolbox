package sheets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
)

const serviceName = "toolbox googleapis.com"

type keyringTokenStore struct {
}

func (k *keyringTokenStore) Get(key string) (*oauth2.Token, error) {

	ring, err := keyring.Open(keyring.Config{ServiceName: serviceName})
	if err != nil {
		return nil, err
	}

	item, err := ring.Get(key)
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot read Google API token for %q: %w", key, err)
	}

	var token oauth2.Token
	err = json.Unmarshal(item.Data, &token)
	if err != nil {
		return nil, fmt.Errorf("invalid Google API token for %q: %w", key, err)
	}

	return &token, nil
}

func (k *keyringTokenStore) Put(key string, token *oauth2.Token) error {

	ring, err := keyring.Open(keyring.Config{ServiceName: serviceName})
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
		Description: "Jira API Token",
	}
	if err := ring.Set(item); err != nil {
		return fmt.Errorf("failed to store Jira API token: %w", err)
	}
	return nil
}

package vault

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/vault/v1"
)

type VaultService struct {
	service *vault.Service
}

func NewService(ctx context.Context, tokenSource oauth2.TokenSource) (*VaultService, error) {
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	service, err := vault.NewService(ctx, option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, err
	}

	return &VaultService{service: service}, nil
}

package cfg

import (
	"context"
	"fmt"
)

type ConfigProvider interface {
	// Getenv reads an environment variable
	Getenv(name string) string

	// ReadFile reads a configuration file from the configuration
	ReadFile(name string) ([]byte, error)
}

type kConfigProviderContextKey struct{}

func FromContext(ctx context.Context) ConfigProvider {
	v := ctx.Value(kConfigProviderContextKey{})
	if v == nil {
		return nil
	}
	cfg, ok := v.(ConfigProvider)
	if !ok {
		panic(fmt.Errorf("unexpected type for ConfigProvider context value: %v", v))
	}
	return cfg
}

func WithConfigProvider(ctx context.Context, cfg ConfigProvider) context.Context {
	return context.WithValue(ctx, kConfigProviderContextKey{}, cfg)
}

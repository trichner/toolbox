package vault

import (
	"context"
	"fmt"

	"google.golang.org/api/vault/v1"
)

//go:generate stringer -type=MatterState
type MatterState int

func (i MatterState) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

const (
	MatterStateUnknown MatterState = iota
	MatterStateUnspecified
	MatterStateOpen
	MatterStateClosed
	MatterStateDeleted
)

type Matter struct {
	Id    string
	Name  string
	State MatterState
}

func (v *VaultService) CreateMatter(ctx context.Context, email string) (*Matter, error) {
	m, err := v.service.Matters.Create(&vault.Matter{Name: email, Description: "-"}).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot create matter for %q: %w", email, err)
	}

	return mapMatter(m), nil
}

func (v *VaultService) FindMatter(ctx context.Context, id string) (*Matter, error) {
	m, err := v.service.Matters.Get(id).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot find matter %q: %w", id, err)
	}

	return mapMatter(m), nil
}

func mapMatter(m *vault.Matter) *Matter {
	if m == nil {
		return nil
	}
	return &Matter{
		Id:    m.MatterId,
		Name:  m.Name,
		State: valueOfState(m.State),
	}
}

func valueOfState(s string) MatterState {
	switch s {
	case "STATE_UNSPECIFIED":
		return MatterStateUnspecified
	case "OPEN":
		return MatterStateOpen
	case "CLOSED":
		return MatterStateClosed
	case "DELETED":
		return MatterStateDeleted
	}

	return MatterStateUnknown
}

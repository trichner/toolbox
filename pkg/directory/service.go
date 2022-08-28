package directory

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type Service struct {
	service *admin.Service
}

func NewService(ctx context.Context, tokenSource oauth2.TokenSource) (*Service, error) {
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	service, err := admin.NewService(ctx, option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, err
	}

	return &Service{service: service}, nil
}

type User struct {
	Id           string
	PrimaryEmail string
	FirstName    string
	LastName     string
}

func (d *Service) FindUser(ctx context.Context, userKey string) (*User, error) {
	u, err := d.service.Users.Get(userKey).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot find user by key %q: %w", userKey, err)
	}

	return mapUser(u), nil
}

func (d *Service) FindUserByPrimaryEmail(ctx context.Context, email string) (*User, error) {
	u, err := d.FindUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("cannot find user by primary email %q: %w", email, err)
	}

	if u.PrimaryEmail != email {
		return nil, fmt.Errorf("email %q is not the primary for %q", email, u.PrimaryEmail)
	}
	return u, nil
}

func (d *Service) DeleteUserByPrimaryEmail(ctx context.Context, primaryEmail string) (*User, error) {
	u, err := d.FindUserByPrimaryEmail(ctx, primaryEmail)
	if err != nil {
		return nil, err
	}

	err = d.service.Users.Delete(u.Id).Context(ctx).Do()
	return u, err
}

func (d *Service) SuspendUserByPrimaryEmail(ctx context.Context, primaryEmail string, suspended bool) (*User, error) {
	u, err := d.FindUserByPrimaryEmail(ctx, primaryEmail)
	if err != nil {
		return nil, err
	}

	newUser, err := d.service.Users.Update(u.PrimaryEmail, &admin.User{Suspended: suspended}).Context(ctx).Do()
	return mapUser(newUser), err
}

func mapUser(u *admin.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		Id:           u.Id,
		PrimaryEmail: u.PrimaryEmail,
		FirstName:    u.Name.GivenName,
		LastName:     u.Name.FamilyName,
	}
}

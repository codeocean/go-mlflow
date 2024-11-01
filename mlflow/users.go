package mlflow

import (
	"context"
	"net/url"
)

type UserService service

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty"`
}

func (s *UserService) Create(ctx context.Context, username, password string) (*User, error) {
	opts := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Username: username,
		Password: password,
	}

	var res struct {
		User *User `json:"user,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "users/create", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

func (s *UserService) Get(ctx context.Context, username string) (*User, error) {
	var res struct {
		User *User `json:"user,omitempty"`
	}

	params := url.Values{}
	params.Set("username", username)

	_, err := s.client.Do(ctx, "GET", "users/get", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, username, password string) error {
	opts := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Username: username,
		Password: password,
	}

	_, err := s.client.Do(ctx, "PATCH", "users/update-password", nil, &opts, nil)
	return err
}

func (s *UserService) UpdateAdmin(ctx context.Context, username string, isAdmin bool) error {
	opts := struct {
		Username string `json:"username,omitempty"`
		IsAdmin  bool   `json:"is_admin,omitempty"`
	}{
		Username: username,
		IsAdmin:  isAdmin,
	}

	_, err := s.client.Do(ctx, "PATCH", "users/update-admin", nil, &opts, nil)
	return err
}

func (s *UserService) Delete(ctx context.Context, username string) error {
	opts := struct {
		Username string `json:"username,omitempty"`
	}{
		Username: username,
	}

	_, err := s.client.Do(ctx, "DELETE", "users/delete", nil, &opts, nil)
	return err
}

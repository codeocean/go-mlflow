package mlflow

import (
	"context"
	"net/url"
)

// UserService handles communication with the user management related methods of the MLflow API.
//
// Note: User management APIs are only available in MLflow installations with authentication enabled.
type UserService service

// User represents an MLflow user account
type User struct {
	// ID is the unique identifier for the user
	ID int `json:"id,omitempty"`
	// Username is the user's login name
	Username string `json:"username,omitempty"`
	// IsAdmin indicates whether the user has administrator privileges
	IsAdmin bool `json:"is_admin,omitempty"`
}

// Create creates a new user account
//
// Requires administrator privileges.
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

// Get retrieves a user account by username
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

// UpdatePassword changes a user's password
//
// Requires administrator privileges.
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

// UpdateAdmin updates a user's administrator privileges
//
// Requires administrator privileges.
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

// Delete removes a user account
//
// Requires administrator privileges.
func (s *UserService) Delete(ctx context.Context, username string) error {
	opts := struct {
		Username string `json:"username,omitempty"`
	}{
		Username: username,
	}

	_, err := s.client.Do(ctx, "DELETE", "users/delete", nil, &opts, nil)
	return err
}

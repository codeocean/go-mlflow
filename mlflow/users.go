package mlflow

import (
	"context"
)

type UserService service

// User ...
type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty"`
}

// UserResponse ...
type UserResponse struct {
	User *User `json:"user,omitempty"`
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

	_, err := s.client.Do(ctx, "POST", "/users/create", &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

func (s *UserService) Get(ctx context.Context, username string) (*User, error) {
	var res struct {
		User *User `json:"user,omitempty"`
	}

	params := QueryParams{"username": username}
	url := "/users/get?" + params.String()

	_, err := s.client.Do(ctx, "GET", url, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

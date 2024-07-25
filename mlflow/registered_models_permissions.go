package mlflow

import (
	"context"
	"net/url"
)

type RegisteredModelService service

type RegisteredModelPermission struct {
	Name       string     `json:"name"`
	UserID     int        `json:"user_id"`
	Permission Permission `json:"permission"`
}

func (s *RegisteredModelService) CreatePermission(ctx context.Context, name, username string, permission Permission) (*RegisteredModelPermission, error) {
	opts := struct {
		Name       string     `json:"name"`
		Username   string     `json:"username"`
		Permission Permission `json:"permission"`
	}{
		Name:       name,
		Username:   username,
		Permission: permission,
	}

	var res struct {
		RegisteredModelPermission *RegisteredModelPermission `json:"registered_model_permission,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "registered-models/permissions/create", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.RegisteredModelPermission, nil
}

func (s *RegisteredModelService) GetPermission(ctx context.Context, name, username string) (*RegisteredModelPermission, error) {
	var res struct {
		RegisteredModelPermission *RegisteredModelPermission `json:"registered_model_permission,omitempty"`
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("username", username)

	_, err := s.client.Do(ctx, "GET", "registered-models/permissions/get", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.RegisteredModelPermission, nil
}

func (s *RegisteredModelService) UpdatePermission(ctx context.Context, name, username string, permission Permission) error {
	opts := struct {
		Name       string     `json:"name"`
		Username   string     `json:"username"`
		Permission Permission `json:"permission"`
	}{
		Name:       name,
		Username:   username,
		Permission: permission,
	}

	_, err := s.client.Do(ctx, "PATCH", "registered-models/permissions/update", nil, &opts, nil)
	return err
}

func (s *RegisteredModelService) DeletePermission(ctx context.Context, name, username string) error {
	opts := struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}{
		Name:     name,
		Username: username,
	}

	_, err := s.client.Do(ctx, "DELETE", "registered-models/permissions/delete", nil, &opts, nil)
	return err
}

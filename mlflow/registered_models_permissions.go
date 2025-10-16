package mlflow

import (
	"context"
	"net/url"
)

// RegisteredModelService handles communication with the registered model related methods of the MLflow API.
//
// Registered models are part of the MLflow Model Registry, allowing you to manage and version models.
type RegisteredModelService service

// RegisteredModelPermission represents a user's access level to a registered model
type RegisteredModelPermission struct {
	// Name is the name of the registered model
	Name string `json:"name"`
	// UserID is the unique identifier for the user
	UserID int `json:"user_id"`
	// Permission is the level of access granted
	Permission Permission `json:"permission"`
}

// CreatePermission grants a user access to a registered model
//
// Requires MANAGE permission on the registered model.
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

// GetPermission retrieves a user's permission level for a registered model
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

// UpdatePermission modifies a user's permission level for a registered model
//
// Requires MANAGE permission on the registered model.
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

// DeletePermission revokes a user's access to a registered model
//
// Requires MANAGE permission on the registered model.
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

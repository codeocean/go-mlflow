package mlflow

import (
	"context"
	"net/url"
)

// Permission represents the level of access a user has to a resource
type Permission string

const (
	// PermissionRead allows viewing the resource
	PermissionRead Permission = "READ"
	// PermissionEdit allows modifying the resource
	PermissionEdit Permission = "EDIT"
	// PermissionManage allows full control including permission management
	PermissionManage Permission = "MANAGE"
	// PermissionNoPermissions indicates no access
	PermissionNoPermissions Permission = "NO_PERMISSIONS"
)

// ExperimentPermission represents a user's access level to an experiment
type ExperimentPermission struct {
	// ExperimentID is the unique identifier for the experiment
	ExperimentID string `json:"experiment_id"`
	// UserID is the unique identifier for the user
	UserID int `json:"user_id"`
	// Permission is the level of access granted
	Permission Permission `json:"permission"`
}

// CreatePermission grants a user access to an experiment
//
// Requires MANAGE permission on the experiment.
func (s *ExperimentService) CreatePermission(ctx context.Context, id, username string, permission Permission) (*ExperimentPermission, error) {
	opts := struct {
		ExperimentID string     `json:"experiment_id"`
		Username     string     `json:"username"`
		Permission   Permission `json:"permission"`
	}{
		ExperimentID: id,
		Username:     username,
		Permission:   permission,
	}

	var res struct {
		ExperimentPermission *ExperimentPermission `json:"experiment_permission,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "experiments/permissions/create", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.ExperimentPermission, nil
}

// GetPermission retrieves a user's permission level for an experiment
func (s *ExperimentService) GetPermission(ctx context.Context, id, username string) (*ExperimentPermission, error) {
	var res struct {
		ExperimentPermission *ExperimentPermission `json:"experiment_permission,omitempty"`
	}

	params := url.Values{}
	params.Set("experiment_id", id)
	params.Set("username", username)

	_, err := s.client.Do(ctx, "GET", "experiments/permissions/get", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.ExperimentPermission, nil
}

// UpdatePermission modifies a user's permission level for an experiment
//
// Requires MANAGE permission on the experiment.
func (s *ExperimentService) UpdatePermission(ctx context.Context, id, username string, permission Permission) error {
	opts := struct {
		ExperimentID string     `json:"experiment_id"`
		Username     string     `json:"username"`
		Permission   Permission `json:"permission"`
	}{
		ExperimentID: id,
		Username:     username,
		Permission:   permission,
	}

	_, err := s.client.Do(ctx, "PATCH", "experiments/permissions/update", nil, &opts, nil)
	return err
}

// DeletePermission revokes a user's access to an experiment
//
// Requires MANAGE permission on the experiment.
func (s *ExperimentService) DeletePermission(ctx context.Context, id, username string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id"`
		Username     string `json:"username"`
	}{
		ExperimentID: id,
		Username:     username,
	}

	_, err := s.client.Do(ctx, "DELETE", "experiments/permissions/delete", nil, &opts, nil)
	return err
}

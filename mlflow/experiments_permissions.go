package mlflow

import (
	"context"
	"net/url"
)

type Permission string

const (
	PermissionRead          Permission = "READ"
	PermissionEdit          Permission = "EDIT"
	PermissionManage        Permission = "MANAGE"
	PermissionNoPermissions Permission = "NO_PERMISSIONS"
)

type ExperimentPermission struct {
	ExperimentID string     `json:"experiment_id"`
	UserID       int        `json:"user_id"`
	Permission   Permission `json:"permission"`
}

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

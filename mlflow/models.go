package mlflow

import (
	"context"
	"strconv"
)

type ModelService service

type ModelVersionStatus string

const (
	ModelVersionStatusPending ModelVersionStatus = "PENDING_REGISTRATION"
	ModelVersionStatusFailed  ModelVersionStatus = "FAILED_REGISTRATION"
	ModelVersionStatusReady   ModelVersionStatus = "READY"
)

type ModelVersion struct {
	Name                 string             `json:"name,omitempty"`
	Version              string             `json:"version,omitempty"`
	CreationTimestamp    int64              `json:"creation_timestamp,omitempty"`
	LastUpdatedTimestamp int64              `json:"last_updated_timestamp,omitempty"`
	UserID               string             `json:"user_id,omitempty"`
	CurrentStage         string             `json:"current_stage,omitempty"`
	Description          string             `json:"description,omitempty"`
	Source               string             `json:"source,omitempty"`
	RunID                string             `json:"run_id,omitempty"`
	Status               ModelVersionStatus `json:"status,omitempty"`
	StatusMessage        string             `json:"status_message,omitempty"`
	Tags                 []*ModelVersionTag `json:"tags,omitempty"`
	RunLink              string             `json:"run_link,omitempty"`
	Aliases              []string           `json:"aliases,omitempty"`
}

type ModelVersionTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *ModelService) SetTag(ctx context.Context, name string, version int, key, value string) error {
	opts := struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
		Key     string `json:"key,omitempty"`
		Value   string `json:"value,omitempty"`
	}{
		Name:    name,
		Version: strconv.Itoa(version),
		Key:     key,
		Value:   value,
	}

	_, err := s.client.Do(ctx, "POST", "model-versions/set-tag", nil, &opts, nil)
	return err
}

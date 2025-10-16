package mlflow

import (
	"context"
)

// ModelVersionService handles communication with the model version related methods of the MLflow API.
//
// Model versions are specific iterations of registered models in the MLflow Model Registry.
type ModelVersionService service

// ModelVersionStatus represents the current state of a model version registration
type ModelVersionStatus string

const (
	// ModelVersionStatusPending indicates the model version is being registered
	ModelVersionStatusPending ModelVersionStatus = "PENDING_REGISTRATION"
	// ModelVersionStatusFailed indicates the model version registration failed
	ModelVersionStatusFailed ModelVersionStatus = "FAILED_REGISTRATION"
	// ModelVersionStatusReady indicates the model version is successfully registered and ready
	ModelVersionStatusReady ModelVersionStatus = "READY"
)

// ModelVersion represents a specific version of a registered model
type ModelVersion struct {
	// Name is the name of the registered model
	Name string `json:"name,omitempty"`
	// Version is the version number as a string
	Version string `json:"version,omitempty"`
	// CreationTimestamp is when this version was created (Unix timestamp in milliseconds)
	CreationTimestamp int64 `json:"creation_timestamp,omitempty"`
	// LastUpdatedTimestamp is when this version was last modified (Unix timestamp in milliseconds)
	LastUpdatedTimestamp int64 `json:"last_updated_timestamp,omitempty"`
	// UserID is the ID of the user who created this version
	UserID string `json:"user_id,omitempty"`
	// CurrentStage is the current stage of the model version (e.g., "Staging", "Production")
	CurrentStage string `json:"current_stage,omitempty"`
	// Description is a human-readable description of the model version
	Description string `json:"description,omitempty"`
	// Source is the URI of the model artifacts
	Source string `json:"source,omitempty"`
	// RunID is the ID of the run that produced this model version
	RunID string `json:"run_id,omitempty"`
	// Status is the current registration status
	Status ModelVersionStatus `json:"status,omitempty"`
	// StatusMessage provides additional details about the status
	StatusMessage string `json:"status_message,omitempty"`
	// Tags is the list of tags associated with this model version
	Tags []*ModelVersionTag `json:"tags,omitempty"`
	// RunLink is a link to the run that produced this model version
	RunLink string `json:"run_link,omitempty"`
	// Aliases is the list of aliases pointing to this version
	Aliases []string `json:"aliases,omitempty"`
}

// ModelVersionTag represents a key-value tag on a model version
type ModelVersionTag struct {
	// Key is the tag name
	Key string `json:"key"`
	// Value is the tag value
	Value string `json:"value"`
}

// SetTag sets a tag on a model version
//
// If the tag already exists, its value will be updated.
func (s *ModelVersionService) SetTag(ctx context.Context, name, version, key, value string) error {
	opts := struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
		Key     string `json:"key,omitempty"`
		Value   string `json:"value,omitempty"`
	}{
		Name:    name,
		Version: version,
		Key:     key,
		Value:   value,
	}

	_, err := s.client.Do(ctx, "POST", "model-versions/set-tag", nil, &opts, nil)
	return err
}

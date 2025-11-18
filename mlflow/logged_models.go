package mlflow

import (
	"context"
	"fmt"
	"net/url"
)

// LoggedModelService provides methods for interacting with logged models in MLflow experiments.
// It enables operations such as retrieving, listing, and managing models that have been logged to the tracking server.
type LoggedModelService service

// LoggedModelStatus represents the status of a logged model
type LoggedModelStatus string

const (
	// LoggedModelStatusReady indicates the model is successfully logged and ready to use
	LoggedModelStatusReady LoggedModelStatus = "LOGGED_MODEL_READY"
	// LoggedModelStatusPending indicates the model is being registered
	LoggedModelStatusPending LoggedModelStatus = "LOGGED_MODEL_PENDING_REGISTRATION"
	// LoggedModelStatusFailed indicates the model registration failed
	LoggedModelStatusFailed LoggedModelStatus = "LOGGED_MODEL_FAILED_REGISTRATION"
	// LoggedModelStatusPendingDel indicates the model is pending deletion
	LoggedModelStatusPendingDel LoggedModelStatus = "LOGGED_MODEL_PENDING_DELETION"
	// LoggedModelStatusDeleting indicates the model is being deleted
	LoggedModelStatusDeleting LoggedModelStatus = "LOGGED_MODEL_DELETING"
	// LoggedModelStatusRegistration indicates the model registration is being finalized
	LoggedModelStatusRegistration LoggedModelStatus = "LOGGED_MODEL_PENDING_FINALIZATION"
)

// LoggedModel represents a model logged to an MLflow Experiment
type LoggedModel struct {
	// Info contains metadata about the logged model
	Info *LoggedModelInfo `json:"info,omitempty"`
	// Data contains parameters and metrics associated with the model
	Data *LoggedModelData `json:"data,omitempty"`
}

// LoggedModelInfo contains metadata about a logged model
type LoggedModelInfo struct {
	// ModelID is the unique identifier for the logged model
	ModelID string `json:"model_id,omitempty"`
	// ExperimentID is the ID of the experiment this model belongs to
	ExperimentID string `json:"experiment_id,omitempty"`
	// Name is the human-readable name of the model
	Name string `json:"name,omitempty"`
	// ArtifactUri is the URI where the model artifacts are stored
	ArtifactUri string `json:"artifact_uri,omitempty"`
	// CreationTimestampMs is when the model was created (Unix timestamp in milliseconds)
	CreationTimestampMs int64 `json:"creation_timestamp_ms,omitempty"`
	// LastUpdatedTimestampMs is when the model was last modified (Unix timestamp in milliseconds)
	LastUpdatedTimestampMs int64 `json:"last_updated_timestamp_ms,omitempty"`
	// Status is the current status of the model
	Status LoggedModelStatus `json:"status,omitempty"`
	// StatusMessage provides additional details about the status
	StatusMessage string `json:"status_message,omitempty"`
	// ModelType is the type/flavor of the model (e.g., "sklearn", "pytorch")
	ModelType string `json:"model_type,omitempty"`
	// SourceRunID is the ID of the run that produced this model
	SourceRunID string `json:"source_run_id,omitempty"`
	// CreatorID is the ID of the user who created the model
	CreatorID int64 `json:"creator_id,omitempty"`
	// Tags is the list of tags associated with this model
	Tags []*LoggedModelTag `json:"tags,omitempty"`
	// Registrations is the list of Model Registry registrations for this model
	Registrations []*LoggedModelRegistrationInfo `json:"registrations,omitempty"`
}

// LoggedModelData contains data associated with a logged model
type LoggedModelData struct {
	// Params is the list of parameters associated with this model
	Params []*LoggedModelParameter `json:"params,omitempty"`
	// Metrics is the list of metrics recorded for this model
	Metrics []*Metric `json:"metrics,omitempty"`
}

// LoggedModelTag represents a tag on a logged model
type LoggedModelTag struct {
	// Key is the tag name
	Key string `json:"key,omitempty"`
	// Value is the tag value
	Value string `json:"value,omitempty"`
}

// LoggedModelParameter represents a parameter of a logged model
type LoggedModelParameter struct {
	// Key is the parameter name
	Key string `json:"key,omitempty"`
	// Value is the parameter value
	Value string `json:"value,omitempty"`
}

// LoggedModelRegistrationInfo contains information about model registry registrations
type LoggedModelRegistrationInfo struct {
	// Name is the name of the registered model
	Name string `json:"name,omitempty"`
	// Version is the version number of the registered model
	Version string `json:"version,omitempty"`
}

// LoggedModelSearchOptions contains options for searching logged models
type LoggedModelSearchOptions struct {
	// ExperimentIDs is the list of experiment IDs to search within (required)
	ExperimentIDs []string `json:"experiment_ids"`
	// Filter is a search filter expression (e.g., "name = 'my-model'")
	Filter string `json:"filter,omitempty"`
	// MaxResults is the maximum number of models to return
	MaxResults int32 `json:"max_results,omitempty"`
	// OrderBy specifies how to sort the results
	OrderBy []*LoggedModelOrderBy `json:"order_by,omitempty"`
	// PageToken is used for pagination to fetch the next page of results
	PageToken string `json:"page_token,omitempty"`
	// Datasets specifies datasets for filtering metrics
	Datasets []*LoggedModelSearchDataset `json:"datasets,omitempty"`
}

// LoggedModelOrderBy specifies ordering for search results
type LoggedModelOrderBy struct {
	// FieldName is the name of the field to sort by
	FieldName string `json:"field_name"`
	// Ascending indicates whether to sort in ascending order (true) or descending (false)
	Ascending bool `json:"ascending,omitempty"`
	// DatasetName is the name of the dataset when sorting by dataset-specific metrics
	DatasetName string `json:"dataset_name,omitempty"`
	// DatasetDigest is the digest of the dataset when sorting by dataset-specific metrics
	DatasetDigest string `json:"dataset_digest,omitempty"`
}

// LoggedModelSearchDataset specifies a dataset for filtering metrics
type LoggedModelSearchDataset struct {
	// DatasetName is the name of the dataset
	DatasetName string `json:"dataset_name"`
	// DatasetDigest is the digest/hash of the dataset for versioning
	DatasetDigest string `json:"dataset_digest,omitempty"`
}

// LoggedModelSearchResults contains the results of a logged model search
type LoggedModelSearchResults struct {
	// Models is the list of models matching the search criteria
	Models []*LoggedModel `json:"models,omitempty"`
	// NextPageToken is used to retrieve the next page of results (empty if no more results)
	NextPageToken string `json:"next_page_token,omitempty"`
}

// LoggedModelArtifactList contains artifacts for a logged model
type LoggedModelArtifactList struct {
	// RootURI is the root location for storing artifacts
	RootURI string `json:"root_uri,omitempty"`
	// Files is the list of artifact files and directories
	Files []*FileInfo `json:"files,omitempty"`
	// NextPageToken is used to retrieve the next page of results (empty if no more results)
	NextPageToken string `json:"next_page_token,omitempty"`
}

// LoggedModelCreateOptions contains options for creating a logged model
type LoggedModelCreateOptions struct {
	// ExperimentID is the ID of the experiment to create the model in (required)
	ExperimentID string `json:"experiment_id"`
	// Name is the human-readable name for the model
	Name string `json:"name,omitempty"`
	// ModelType is the type/flavor of the model (e.g., "sklearn", "pytorch")
	ModelType string `json:"model_type,omitempty"`
	// SourceRunID is the ID of the run that produced this model
	SourceRunID string `json:"source_run_id,omitempty"`
	// Params is the list of initial parameters for the model
	Params []*LoggedModelParameter `json:"params,omitempty"`
	// Tags is the list of initial tags for the model
	Tags []*LoggedModelTag `json:"tags,omitempty"`
}

// Create creates a new logged model
func (s *LoggedModelService) Create(ctx context.Context, opts *LoggedModelCreateOptions) (*LoggedModel, error) {
	var res struct {
		Model *LoggedModel `json:"model,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "logged-models", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

// Get retrieves a logged model by ID
func (s *LoggedModelService) Get(ctx context.Context, modelID string) (*LoggedModel, error) {
	var res struct {
		Model *LoggedModel `json:"model,omitempty"`
	}

	_, err := s.client.Do(ctx, "GET", fmt.Sprintf("logged-models/%s", modelID), nil, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

// GetBatch retrieves multiple logged models by their IDs
func (s *LoggedModelService) GetBatch(ctx context.Context, modelIDs []string) ([]*LoggedModel, error) {
	params := url.Values{}
	for _, id := range modelIDs {
		params.Add("model_ids", id)
	}

	var res struct {
		Models []*LoggedModel `json:"models,omitempty"`
	}

	_, err := s.client.Do(ctx, "GET", "logged-models:batchGet", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Models, nil
}

// Delete deletes a logged model
func (s *LoggedModelService) Delete(ctx context.Context, modelID string) error {
	opts := struct {
		ModelID string `json:"model_id"`
	}{
		ModelID: modelID,
	}

	_, err := s.client.Do(ctx, "DELETE", fmt.Sprintf("logged-models/%s", modelID), nil, &opts, nil)
	return err
}

// Search searches for logged models matching the specified criteria
func (s *LoggedModelService) Search(ctx context.Context, opts *LoggedModelSearchOptions) (*LoggedModelSearchResults, error) {
	var res LoggedModelSearchResults

	_, err := s.client.Do(ctx, "POST", "logged-models/search", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SetTags sets tags on a logged model
func (s *LoggedModelService) SetTags(ctx context.Context, modelID string, tags map[string]string) (*LoggedModel, error) {
	tagList := make([]*LoggedModelTag, 0, len(tags))
	for key, value := range tags {
		tagList = append(tagList, &LoggedModelTag{Key: key, Value: value})
	}

	opts := struct {
		Tags []*LoggedModelTag `json:"tags"`
	}{
		Tags: tagList,
	}

	var res struct {
		Model *LoggedModel `json:"model,omitempty"`
	}

	_, err := s.client.Do(ctx, "PATCH", fmt.Sprintf("logged-models/%s/tags", modelID), nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

// DeleteTag deletes a tag from a logged model
func (s *LoggedModelService) DeleteTag(ctx context.Context, modelID, tagKey string) error {
	_, err := s.client.Do(ctx, "DELETE", fmt.Sprintf("logged-models/%s/tags/%s", modelID, tagKey), nil, nil, nil)
	return err
}

// LogParams logs parameters to a logged model
func (s *LoggedModelService) LogParams(ctx context.Context, modelID string, params map[string]string) error {
	paramList := make([]*LoggedModelParameter, 0, len(params))
	for key, value := range params {
		paramList = append(paramList, &LoggedModelParameter{Key: key, Value: value})
	}

	opts := struct {
		ModelID string                  `json:"model_id"`
		Params  []*LoggedModelParameter `json:"params"`
	}{
		ModelID: modelID,
		Params:  paramList,
	}

	_, err := s.client.Do(ctx, "POST", fmt.Sprintf("logged-models/%s/params", modelID), nil, &opts, nil)
	return err
}

// ListArtifacts lists artifacts for a logged model
func (s *LoggedModelService) ListArtifacts(ctx context.Context, modelID, path string) (*LoggedModelArtifactList, error) {
	params := url.Values{}
	if path != "" {
		params.Set("artifact_directory_path", path)
	}

	var res LoggedModelArtifactList

	_, err := s.client.Do(ctx, "GET", fmt.Sprintf("logged-models/%s/artifacts/directories", modelID), params, nil, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Finalize finalizes a logged model with the given status
func (s *LoggedModelService) Finalize(ctx context.Context, modelID string, status LoggedModelStatus) (*LoggedModel, error) {
	opts := struct {
		ModelID string            `json:"model_id"`
		Status  LoggedModelStatus `json:"status"`
	}{
		ModelID: modelID,
		Status:  status,
	}

	var res struct {
		Model *LoggedModel `json:"model,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "logged-models/finalize", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

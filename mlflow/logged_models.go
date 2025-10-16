package mlflow

import (
	"context"
	"fmt"
	"net/url"
)

type LoggedModelService service

// LoggedModelStatus represents the status of a logged model
type LoggedModelStatus string

const (
	LoggedModelStatusReady        LoggedModelStatus = "LOGGED_MODEL_READY"
	LoggedModelStatusPending      LoggedModelStatus = "LOGGED_MODEL_PENDING_REGISTRATION"
	LoggedModelStatusFailed       LoggedModelStatus = "LOGGED_MODEL_FAILED_REGISTRATION"
	LoggedModelStatusPendingDel   LoggedModelStatus = "LOGGED_MODEL_PENDING_DELETION"
	LoggedModelStatusDeleting     LoggedModelStatus = "LOGGED_MODEL_DELETING"
	LoggedModelStatusRegistration LoggedModelStatus = "LOGGED_MODEL_PENDING_FINALIZATION"
)

// LoggedModel represents a model logged to an MLflow Experiment
type LoggedModel struct {
	Info *LoggedModelInfo `json:"info,omitempty"`
	Data *LoggedModelData `json:"data,omitempty"`
}

// LoggedModelInfo contains metadata about a logged model
type LoggedModelInfo struct {
	ModelID                string                         `json:"model_id,omitempty"`
	ExperimentID           string                         `json:"experiment_id,omitempty"`
	Name                   string                         `json:"name,omitempty"`
	ArtifactUri            string                         `json:"artifact_uri,omitempty"`
	CreationTimestampMs    int64                          `json:"creation_timestamp_ms,omitempty"`
	LastUpdatedTimestampMs int64                          `json:"last_updated_timestamp_ms,omitempty"`
	Status                 LoggedModelStatus              `json:"status,omitempty"`
	StatusMessage          string                         `json:"status_message,omitempty"`
	ModelType              string                         `json:"model_type,omitempty"`
	SourceRunID            string                         `json:"source_run_id,omitempty"`
	CreatorID              int64                          `json:"creator_id,omitempty"`
	Tags                   []*LoggedModelTag              `json:"tags,omitempty"`
	Registrations          []*LoggedModelRegistrationInfo `json:"registrations,omitempty"`
}

// LoggedModelData contains data associated with a logged model
type LoggedModelData struct {
	Params  []*LoggedModelParameter `json:"params,omitempty"`
	Metrics []*Metric               `json:"metrics,omitempty"`
}

// LoggedModelTag represents a tag on a logged model
type LoggedModelTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// LoggedModelParameter represents a parameter of a logged model
type LoggedModelParameter struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// LoggedModelRegistrationInfo contains information about model registry registrations
type LoggedModelRegistrationInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// LoggedModelSearchOptions contains options for searching logged models
type LoggedModelSearchOptions struct {
	ExperimentIDs []string                    `json:"experiment_ids"`
	Filter        string                      `json:"filter,omitempty"`
	MaxResults    int32                       `json:"max_results,omitempty"`
	OrderBy       []*LoggedModelOrderBy       `json:"order_by,omitempty"`
	PageToken     string                      `json:"page_token,omitempty"`
	Datasets      []*LoggedModelSearchDataset `json:"datasets,omitempty"`
}

// LoggedModelOrderBy specifies ordering for search results
type LoggedModelOrderBy struct {
	FieldName     string `json:"field_name"`
	Ascending     bool   `json:"ascending,omitempty"`
	DatasetName   string `json:"dataset_name,omitempty"`
	DatasetDigest string `json:"dataset_digest,omitempty"`
}

// LoggedModelSearchDataset specifies a dataset for filtering metrics
type LoggedModelSearchDataset struct {
	DatasetName   string `json:"dataset_name"`
	DatasetDigest string `json:"dataset_digest,omitempty"`
}

// LoggedModelSearchResults contains the results of a logged model search
type LoggedModelSearchResults struct {
	Models        []*LoggedModel `json:"models,omitempty"`
	NextPageToken string         `json:"next_page_token,omitempty"`
}

// LoggedModelArtifactList contains artifacts for a logged model
type LoggedModelArtifactList struct {
	RootURI       string      `json:"root_uri,omitempty"`
	Files         []*FileInfo `json:"files,omitempty"`
	NextPageToken string      `json:"next_page_token,omitempty"`
}

// Create creates a new logged model
func (s *LoggedModelService) Create(ctx context.Context, experimentID string, opts *LoggedModelCreateOptions) (*LoggedModel, error) {
	if opts == nil {
		opts = &LoggedModelCreateOptions{}
	}
	opts.ExperimentID = experimentID

	var res struct {
		Model *LoggedModel `json:"model,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "logged-models/create", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Model, nil
}

// LoggedModelCreateOptions contains options for creating a logged model
type LoggedModelCreateOptions struct {
	ExperimentID string                  `json:"experiment_id"`
	Name         string                  `json:"name,omitempty"`
	ModelType    string                  `json:"model_type,omitempty"`
	SourceRunID  string                  `json:"source_run_id,omitempty"`
	Params       []*LoggedModelParameter `json:"params,omitempty"`
	Tags         []*LoggedModelTag       `json:"tags,omitempty"`
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
	if opts == nil {
		return nil, fmt.Errorf("search options are required")
	}

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
func (s *LoggedModelService) ListArtifacts(ctx context.Context, modelID string, path string) (*LoggedModelArtifactList, error) {
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

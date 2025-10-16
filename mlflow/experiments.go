package mlflow

import (
	"context"
	"net/url"
)

// ExperimentService handles operations related to MLflow experiments
type ExperimentService service

// Experiment represents an MLflow experiment, which groups related runs
type Experiment struct {
	ExperimentID     string           `json:"experiment_id,omitempty"`
	Name             string           `json:"name,omitempty"`
	ArtifactLocation string           `json:"artifact_location,omitempty"`
	LifecycleStage   string           `json:"lifecycle_stage,omitempty"`
	LastUpdateTime   int64            `json:"last_update_time,omitempty"`
	CreationTime     int64            `json:"creation_time,omitempty"`
	Tags             []*ExperimentTag `json:"tags,omitempty"`
}

// ExperimentTag represents a tag associated with an experiment
type ExperimentTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// ExperimentsSearchOptions contains options for searching experiments
type ExperimentsSearchOptions struct {
	Filter     string   `json:"filter,omitempty"`
	ViewType   ViewType `json:"view_type,omitempty"`
	MaxResults int64    `json:"max_results,omitempty"`
	OrderBy    []string `json:"order_by,omitempty"`
	PageToken  string   `json:"page_token,omitempty"`
}

// ExperimentsSearchResults contains the results of an experiment search
type ExperimentsSearchResults struct {
	Experiments   []*Experiment `json:"experiments,omitempty"`
	NextPageToken string        `json:"next_page_token,omitempty"`
}

// Create creates a new MLflow experiment
//
// Returns the ID of the created experiment or an error.
// Experiment names must be unique within an MLflow tracking server.
func (s *ExperimentService) Create(ctx context.Context, name string) (string, error) {
	opts := struct {
		Name string `json:"name,omitempty"`
	}{
		Name: name,
	}

	var res struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "experiments/create", nil, &opts, &res)
	if err != nil {
		return "", err
	}

	return res.ExperimentID, nil
}

// Update renames an existing MLflow experiment
//
// The experiment must not be deleted. Use Restore first if needed.
func (s *ExperimentService) Update(ctx context.Context, id, name string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
		NewName      string `json:"new_name,omitempty"`
	}{
		ExperimentID: id,
		NewName:      name,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/update", nil, &opts, nil)
	return err
}

// Delete marks an MLflow experiment as deleted
//
// The experiment and all its runs will be moved to deleted state.
// Use Restore to undelete if needed.
func (s *ExperimentService) Delete(ctx context.Context, id string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/delete", nil, &opts, nil)
	return err
}

// Restore restores a previously deleted experiment
//
// This also restores all runs that were part of the experiment.
func (s *ExperimentService) Restore(ctx context.Context, id string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/restore", nil, &opts, nil)
	return err
}

// SetTag sets a tag on an experiment
//
// If the tag already exists, its value will be updated.
func (s *ExperimentService) SetTag(ctx context.Context, id, key, value string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
		Key          string `json:"key,omitempty"`
		Value        string `json:"value,omitempty"`
	}{
		ExperimentID: id,
		Key:          key,
		Value:        value,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/set-experiment-tag", nil, &opts, nil)
	return err
}

// Get retrieves the metadata for a single experiment by ID
func (s *ExperimentService) Get(ctx context.Context, id string) (*Experiment, error) {
	var res struct {
		Experiment *Experiment `json:"experiment,omitempty"`
	}

	params := url.Values{}
	params.Set("experiment_id", id)

	_, err := s.client.Do(ctx, "GET", "experiments/get", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Experiment, nil
}

// GetByName retrieves the metadata for a single experiment by name
func (s *ExperimentService) GetByName(ctx context.Context, name string) (*Experiment, error) {
	var res struct {
		Experiment *Experiment `json:"experiment,omitempty"`
	}

	params := url.Values{}
	params.Set("experiment_name", name)

	_, err := s.client.Do(ctx, "GET", "experiments/get-by-name", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Experiment, nil
}

// Search queries experiments using a filter expression
//
// The filter string is a search expression like:
//   - "name = 'my-experiment'"
//   - "tags.team = 'ml-platform'"
//   - "name LIKE 'experiment-%'"
//
// Results can be ordered and paginated using the options parameter.
func (s *ExperimentService) Search(ctx context.Context, opts *ExperimentsSearchOptions) (*ExperimentsSearchResults, error) {
	var res ExperimentsSearchResults

	_, err := s.client.Do(ctx, "POST", "experiments/search", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

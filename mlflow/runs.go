package mlflow

import (
	"context"
	"net/url"
	"time"
)

// RunService handles operations related to MLflow runs
type RunService service

// RunStatus represents the execution status of an MLflow run
type RunStatus string

const (
	// RunStatusRunning indicates the run is currently executing
	RunStatusRunning RunStatus = "RUNNING"
	// RunStatusScheduled indicates the run is scheduled but not yet started
	RunStatusScheduled RunStatus = "SCHEDULED"
	// RunStatusFinished indicates the run completed successfully
	RunStatusFinished RunStatus = "FINISHED"
	// RunStatusFailed indicates the run completed with errors
	RunStatusFailed RunStatus = "FAILED"
	// RunStatusKilled indicates the run was terminated by user action
	RunStatusKilled RunStatus = "KILLED"
)

// ViewType specifies which runs to return in a search based on their lifecycle stage
type ViewType string

const (
	// ViewTypeActiveOnly returns only active runs
	ViewTypeActiveOnly ViewType = "ACTIVE_ONLY"
	// ViewTypeDeletedOnly returns only deleted runs
	ViewTypeDeletedOnly ViewType = "DELETED_ONLY"
	// ViewTypeAll returns all runs regardless of lifecycle stage
	ViewTypeAll ViewType = "ALL"
)

// Run represents a single execution of an MLflow experiment
type Run struct {
	// Info contains metadata about the run
	Info *RunInfo `json:"info,omitempty"`
	// Data contains metrics, parameters, and tags logged during the run
	Data *RunData `json:"data,omitempty"`
	// Inputs contains datasets and models used as inputs
	Inputs *RunInputs `json:"inputs,omitempty"`
	// Outputs contains models produced as outputs
	Outputs *RunOutputs `json:"outputs,omitempty"`
}

// RunInfo contains metadata about an MLflow run
type RunInfo struct {
	// RunID is the unique identifier for the run
	RunID string `json:"run_id,omitempty"`
	// RunName is the human-readable name of the run
	RunName string `json:"run_name,omitempty"`
	// ExperimentID is the ID of the experiment this run belongs to
	ExperimentID string `json:"experiment_id,omitempty"`
	// Status is the current execution status of the run
	Status RunStatus `json:"status,omitempty"`
	// StartTime is when the run started (Unix timestamp in milliseconds)
	StartTime int64 `json:"start_time,omitempty"`
	// EndTime is when the run ended (Unix timestamp in milliseconds)
	EndTime int64 `json:"end_time,omitempty"`
	// ArtifactUri is the URI where artifacts are stored
	ArtifactUri string `json:"artifact_uri,omitempty"`
	// LifecycleStage is the lifecycle stage (e.g., "active" or "deleted")
	LifecycleStage string `json:"lifecycle_stage,omitempty"`
}

// RunData contains the metrics, parameters, and tags logged during a run
type RunData struct {
	// Metrics is the list of metrics logged to the run
	Metrics []*Metric `json:"metrics,omitempty"`
	// Params is the list of parameters logged to the run
	Params []*Param `json:"params,omitempty"`
	// Tags is the list of tags associated with the run
	Tags []*RunTag `json:"tags,omitempty"`
}

// Metric represents a numeric value logged during a run at a specific step
type Metric struct {
	// Key is the metric name
	Key string `json:"key,omitempty"`
	// Value is the numeric value of the metric
	Value float64 `json:"value,omitempty"`
	// Timestamp is when the metric was logged (Unix timestamp in milliseconds)
	Timestamp int64 `json:"timestamp,omitempty"`
	// Step is the training step at which the metric was logged
	Step int64 `json:"step,omitempty"`
}

// Param represents a key-value parameter logged to a run
type Param struct {
	// Key is the parameter name
	Key string `json:"key,omitempty"`
	// Value is the parameter value
	Value string `json:"value,omitempty"`
}

// RunTag represents a tag associated with a run
type RunTag struct {
	// Key is the tag name
	Key string `json:"key,omitempty"`
	// Value is the tag value
	Value string `json:"value,omitempty"`
}

// RunInputs contains datasets and models used as inputs to a run
type RunInputs struct {
	// DatasetInputs is the list of datasets used as inputs
	DatasetInputs []*DatasetInput `json:"dataset_inputs,omitempty"`
	// ModelInputs is the list of models used as inputs
	ModelInputs []*ModelInput `json:"model_inputs,omitempty"`
}

// ModelInput represents a model used as input to a run
type ModelInput struct {
	// ModelID is the unique identifier of the input model
	ModelID string `json:"model_id,omitempty"`
}

// DatasetInput represents a dataset used as input to a run
type DatasetInput struct {
	// Tags is the list of tags associated with this dataset input
	Tags []*InputTag `json:"tags,omitempty"`
	// Dataset contains metadata about the dataset
	Dataset *Dataset `json:"dataset,omitempty"`
}

// InputTag represents a tag on a dataset input
type InputTag struct {
	// Key is the tag name
	Key string `json:"key,omitempty"`
	// Value is the tag value
	Value string `json:"value,omitempty"`
}

// Dataset represents metadata about a dataset used in MLflow
type Dataset struct {
	// Name is the name of the dataset
	Name string `json:"name,omitempty"`
	// Digest is the hash/digest of the dataset for versioning
	Digest string `json:"digest,omitempty"`
	// SourceType is the type of data source (e.g., "path", "delta", "s3")
	SourceType string `json:"source_type,omitempty"`
	// Source is the location of the dataset
	Source string `json:"source,omitempty"`
	// Schema is the schema definition of the dataset
	Schema string `json:"schema,omitempty"`
	// Profile contains profiling statistics about the dataset
	Profile string `json:"profile,omitempty"`
}

// RunOutputs contains models produced as outputs of a run
type RunOutputs struct {
	// ModelOutputs is the list of models produced by the run
	ModelOutputs []*ModelOutput `json:"model_outputs,omitempty"`
}

// ModelOutput represents a model logged as output of a run
type ModelOutput struct {
	// ModelID is the unique identifier of the output model
	ModelID string `json:"model_id,omitempty"`
	// Step is the training step at which the model was logged
	Step int64 `json:"step,omitempty"`
}

// RunSearchOptions contains options for searching MLflow runs
type RunSearchOptions struct {
	// ExperimentIDs is the list of experiment IDs to search within
	ExperimentIDs []string `json:"experiment_ids,omitempty"`
	// Filter is a search filter expression (e.g., "metrics.accuracy > 0.9")
	Filter string `json:"filter,omitempty"`
	// RunViewType specifies which runs to return based on lifecycle stage
	RunViewType ViewType `json:"run_view_type,omitempty"`
	// MaxResults is the maximum number of runs to return
	MaxResults int32 `json:"max_results,omitempty"`
	// OrderBy is the list of order-by clauses (e.g., "metrics.rmse DESC")
	OrderBy []string `json:"order_by,omitempty"`
	// PageToken is used for pagination to fetch the next page of results
	PageToken string `json:"page_token,omitempty"`
}

// RunSearchResults contains the results of a run search operation
type RunSearchResults struct {
	// Runs is the list of runs matching the search criteria
	Runs []*Run `json:"runs,omitempty"`
	// NextPageToken is used to retrieve the next page of results (empty if no more results)
	NextPageToken string `json:"next_page_token,omitempty"`
}

// Create creates a new MLflow run
//
// Parameters:
//   - experimentID: ID of the experiment to create the run in
//   - name: Name for the run
//   - startTime: Start time in milliseconds since epoch (0 uses current time)
//   - tags: Optional key-value tags to attach to the run
//
// Returns the created run or an error
func (s *RunService) Create(ctx context.Context, experimentID, name string, startTime int64, tags map[string]string) (*Run, error) {
	opts := struct {
		ExperimentID string    `json:"experiment_id,omitempty"`
		RunName      string    `json:"run_name,omitempty"`
		StartTime    int64     `json:"start_time,omitempty"`
		Tags         []*RunTag `json:"tags,omitempty"`
	}{
		ExperimentID: experimentID,
		RunName:      name,
		StartTime:    startTime,
	}

	if startTime == 0 {
		opts.StartTime = time.Now().UnixMilli()
	}

	for key, value := range tags {
		opts.Tags = append(opts.Tags, &RunTag{Key: key, Value: value})
	}

	var res struct {
		Run *Run `json:"run,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "runs/create", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Run, nil
}

// Update updates an existing MLflow run's metadata
//
// Parameters:
//   - id: ID of the run to update
//   - name: New name for the run (empty string keeps current name)
//   - status: New status for the run
//   - endTime: End time in milliseconds since epoch (0 uses current time for terminal states)
//
// Returns the updated run info or an error
func (s *RunService) Update(ctx context.Context, id, name string, status RunStatus, endTime int64) (*RunInfo, error) {
	opts := struct {
		RunID   string    `json:"run_id,omitempty"`
		RunName string    `json:"run_name,omitempty"`
		Status  RunStatus `json:"status,omitempty"`
		EndTime int64     `json:"end_time,omitempty"`
	}{
		RunID:   id,
		RunName: name,
		Status:  status,
		EndTime: endTime,
	}

	if endTime == 0 {
		if status == RunStatusFinished || status == RunStatusFailed || status == RunStatusKilled {
			opts.EndTime = time.Now().UnixMilli()
		}
	}

	var res struct {
		Info *RunInfo `json:"info,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "runs/update", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Info, nil
}

// Delete marks an MLflow run as deleted
//
// The run can be restored using Restore(). To permanently delete,
// use the MLflow server's vacuum functionality.
func (s *RunService) Delete(ctx context.Context, id string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
	}{
		RunID: id,
	}

	_, err := s.client.Do(ctx, "POST", "runs/delete", nil, &opts, nil)
	return err
}

// Restore restores a previously deleted MLflow run
//
// The run must have been deleted using Delete() and not yet permanently
// removed from the system.
func (s *RunService) Restore(ctx context.Context, id string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
	}{
		RunID: id,
	}

	_, err := s.client.Do(ctx, "POST", "runs/restore", nil, &opts, nil)
	return err
}

// Get retrieves detailed information about a specific MLflow run
//
// Returns the complete run including metadata, metrics, parameters,
// tags, inputs, and outputs.
func (s *RunService) Get(ctx context.Context, id string) (*Run, error) {
	var res struct {
		Run *Run `json:"run,omitempty"`
	}

	params := url.Values{}
	params.Set("run_id", id)

	_, err := s.client.Do(ctx, "GET", "runs/get", params, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Run, nil
}

// Search searches for MLflow runs matching the specified criteria
//
// Supports filtering, ordering, and pagination. Use RunSearchOptions to
// specify search criteria including experiment IDs, filter expressions,
// view type, ordering, and pagination.
func (s *RunService) Search(ctx context.Context, opts *RunSearchOptions) (*RunSearchResults, error) {
	var res RunSearchResults

	_, err := s.client.Do(ctx, "POST", "runs/search", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SetTag sets a tag on an MLflow run
//
// If the tag already exists, its value will be updated.
func (s *RunService) SetTag(ctx context.Context, id, key, value string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}{
		RunID: id,
		Key:   key,
		Value: value,
	}

	_, err := s.client.Do(ctx, "POST", "runs/set-tag", nil, &opts, nil)
	return err
}

// DeleteTag removes a tag from an MLflow run
func (s *RunService) DeleteTag(ctx context.Context, id, key string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
		Key   string `json:"key,omitempty"`
	}{
		RunID: id,
		Key:   key,
	}

	_, err := s.client.Do(ctx, "POST", "runs/delete-tag", nil, &opts, nil)
	return err
}

// LogMetric logs a metric value for an MLflow run at a specific step
//
// Metrics are numeric values that can be logged multiple times during
// a run's execution, typically to track model performance over iterations.
//
// Parameters:
//   - id: Run ID
//   - key: Metric name
//   - value: Metric value
//   - timestamp: Time in milliseconds since epoch (0 uses current time)
//   - step: Training step or iteration number
func (s *RunService) LogMetric(ctx context.Context, id, key string, value float64, timestamp int64, step int64) error {
	opts := struct {
		RunID     string  `json:"run_id,omitempty"`
		Key       string  `json:"key,omitempty"`
		Value     float64 `json:"value,omitempty"`
		Timestamp int64   `json:"timestamp,omitempty"`
		Step      int64   `json:"step,omitempty"`
	}{
		RunID:     id,
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
		Step:      step,
	}

	_, err := s.client.Do(ctx, "POST", "runs/log-metric", nil, &opts, nil)
	return err
}

// LogParam logs a parameter to an MLflow run
//
// Parameters are key-value pairs that describe the configuration used
// for a run. Each parameter can only be set once per run.
func (s *RunService) LogParam(ctx context.Context, id, key, value string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}{
		RunID: id,
		Key:   key,
		Value: value,
	}

	_, err := s.client.Do(ctx, "POST", "runs/log-parameter", nil, &opts, nil)
	return err
}

// LogBatch logs multiple metrics, parameters, and tags to a run in a single request
//
// This is more efficient than calling LogMetric, LogParam, and SetTag individually
// when logging multiple values.
func (s *RunService) LogBatch(ctx context.Context, id string, data *RunData) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
		*RunData
	}{
		RunID:   id,
		RunData: data,
	}

	_, err := s.client.Do(ctx, "POST", "runs/log-batch", nil, &opts, nil)
	return err
}

// LogInputs logs dataset inputs used by a run
//
// This tracks the datasets consumed during the run's execution,
// enabling data lineage and reproducibility.
func (s *RunService) LogInputs(ctx context.Context, id string, datasets []*DatasetInput) error {
	opts := struct {
		RunID    string          `json:"run_id,omitempty"`
		Datasets []*DatasetInput `json:"datasets,omitempty"`
	}{
		RunID:    id,
		Datasets: datasets,
	}

	_, err := s.client.Do(ctx, "POST", "runs/log-inputs", nil, &opts, nil)
	return err
}

// LogModel logs a model artifact to a run (legacy method)
//
// This method records that a model was logged to the run. For MLflow 2.0+,
// prefer using the LoggedModels API for more comprehensive model tracking.
//
// The model parameter should be a JSON string containing model metadata.
func (s *RunService) LogModel(ctx context.Context, id, model string) error {
	opts := struct {
		RunID     string `json:"run_id,omitempty"`
		ModelJson string `json:"model_json,omitempty"`
	}{
		RunID:     id,
		ModelJson: model,
	}

	_, err := s.client.Do(ctx, "POST", "runs/log-model", nil, &opts, nil)
	return err
}

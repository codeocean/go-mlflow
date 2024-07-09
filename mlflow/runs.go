package mlflow

import (
	"context"
	"time"
)

type RunService service

type RunStatus string

const (
	RunStatusRunning   RunStatus = "RUNNING"
	RunStatusScheduled RunStatus = "SCHEDULED"
	RunStatusFinished  RunStatus = "FINISHED"
	RunStatusFailed    RunStatus = "FAILED"
	RunStatusKilled    RunStatus = "KILLED"
)

type ViewType string

const (
	ViewTypeActiveOnly  ViewType = "ACTIVE_ONLY"
	ViewTypeDeletedOnly ViewType = "DELETED_ONLY"
	ViewTypeAll         ViewType = "ALL"
)

type Run struct {
	Info   *RunInfo   `json:"info,omitempty"`
	Data   *RunData   `json:"data,omitempty"`
	Inputs *RunInputs `json:"inputs,omitempty"`
}

type RunInfo struct {
	RunID          string    `json:"run_id,omitempty"`
	RunName        string    `json:"run_name,omitempty"`
	ExperimentID   string    `json:"experiment_id,omitempty"`
	Status         RunStatus `json:"status,omitempty"`
	StartTime      int64     `json:"start_time,omitempty"`
	EndTime        int64     `json:"end_time,omitempty"`
	ArtifactUri    string    `json:"artifact_uri,omitempty"`
	LifecycleStage string    `json:"lifecycle_stage,omitempty"`
}

type RunData struct {
	Metrics []*Metric `json:"metrics,omitempty"`
	Params  []*Param  `json:"params,omitempty"`
	Tags    []*RunTag `json:"tags,omitempty"`
}

type Metric struct {
	Key       string  `json:"key,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
	Step      int64   `json:"step,omitempty"`
}

type Param struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RunTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RunInputs struct {
	DatasetInputs []*DatasetInput `json:"dataset_inputs,omitempty"`
}

type DatasetInput struct {
	Tags    []*InputTag `json:"tags,omitempty"`
	Dataset *Dataset    `json:"dataset,omitempty"`
}

type InputTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Dataset struct {
	Name       string `json:"name,omitempty"`
	Digest     string `json:"digest,omitempty"`
	SourceType string `json:"source_type,omitempty"`
	Source     string `json:"source,omitempty"`
	Schema     string `json:"schema,omitempty"`
	Profile    string `json:"profile,omitempty"`
}

type RunSearchOptions struct {
	ExperimentIDs []string `json:"experiment_ids,omitempty"`
	Filter        string   `json:"filter,omitempty"`
	RunViewType   ViewType `json:"run_view_type,omitempty"`
	MaxResults    int32    `json:"max_results,omitempty"`
	OrderBy       []string `json:"order_by,omitempty"`
	PageToken     string   `json:"page_token,omitempty"`
}

type RunSearchResults struct {
	Runs          []*Run `json:"runs,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

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

func (s *RunService) Delete(ctx context.Context, id string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
	}{
		RunID: id,
	}

	_, err := s.client.Do(ctx, "POST", "runs/delete", nil, &opts, nil)
	return err
}

func (s *RunService) Restore(ctx context.Context, id string) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
	}{
		RunID: id,
	}

	_, err := s.client.Do(ctx, "POST", "runs/restore", nil, &opts, nil)
	return err
}

func (s *RunService) Get(ctx context.Context, id string) (*Run, error) {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
	}{
		RunID: id,
	}

	var res struct {
		Run *Run `json:"run,omitempty"`
	}

	_, err := s.client.Do(ctx, "GET", "runs/get", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Run, nil
}

func (s *RunService) Search(ctx context.Context, opts *RunSearchOptions) (*RunSearchResults, error) {
	var res RunSearchResults

	_, err := s.client.Do(ctx, "POST", "runs/search", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

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

func (s *RunService) LogBatch(ctx *context.Context, id string, data *RunData) error {
	opts := struct {
		RunID string `json:"run_id,omitempty"`
		*RunData
	}{
		RunID:   id,
		RunData: data,
	}

	_, err := s.client.Do(*ctx, "POST", "runs/log-batch", nil, &opts, nil)
	return err
}

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

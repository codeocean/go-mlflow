package mlflow

import (
	"context"
	"time"
)

type RunService service

type RunStatus int32

const (
	RunStatusRunning   RunStatus = 1
	RunStatusScheduled RunStatus = 2
	RunStatusFinished  RunStatus = 3
	RunStatusFailed    RunStatus = 4
	RunStatusKilled    RunStatus = 5
)

type Run struct {
	Info *RunInfo `json:"info,omitempty"`
	Data *RunData `json:"data,omitempty"`
}

type RunInfo struct {
	RunID          string    `json:"run_id,omitempty"`
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

func (s *RunService) Create(ctx context.Context, experimentID, name string, tags map[string]string) (*Run, error) {
	opts := struct {
		ExperimentID string    `json:"experiment_id,omitempty"`
		RunName      string    `json:"run_name,omitempty"`
		StartTime    int64     `json:"start_time,omitempty"`
		Tags         []*RunTag `json:"tags,omitempty"`
	}{
		ExperimentID: experimentID,
		RunName:      name,
		StartTime:    time.Now().UnixMilli(),
	}

	for key, value := range tags {
		opts.Tags = append(opts.Tags, &RunTag{Key: key, Value: value})
	}

	var res struct {
		Run *Run `json:"run,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "/runs/create", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Run, nil
}

func (s *RunService) Update(ctx context.Context, id, name string, status RunStatus) (*RunInfo, error) {
	opts := struct {
		RunID   string    `json:"run_id,omitempty"`
		RunName string    `json:"run_name,omitempty"`
		Status  RunStatus `json:"status,omitempty"`
		EndTime int64     `json:"end_time,omitempty"`
	}{
		RunID:   id,
		RunName: name,
		Status:  status,
	}

	if status == RunStatusFinished || status == RunStatusFailed || status == RunStatusKilled {
		opts.EndTime = time.Now().UnixMilli()
	}

	var res struct {
		Info *RunInfo `json:"info,omitempty"`
	}

	_, err := s.client.Do(ctx, "POST", "/runs/update", nil, &opts, &res)
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

	_, err := s.client.Do(ctx, "POST", "/runs/delete", nil, &opts, nil)
	return err
}
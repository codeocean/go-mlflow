package mlflow

import "context"

type ExperimentService service

type Experiment struct {
	ExperimentID     string           `json:"experiment_id,omitempty"`
	Name             string           `json:"name,omitempty"`
	ArtifactLocation string           `json:"artifact_location,omitempty"`
	LifecycleStage   string           `json:"lifecycle_stage,omitempty"`
	LastUpdateTime   int64            `json:"last_update_time,omitempty"`
	CreationTime     int64            `json:"creation_time,omitempty"`
	Tags             []*ExperimentTag `json:"tags,omitempty"`
}

type ExperimentTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type ExperimentsSearchOptions struct {
	Filter     string   `json:"filter,omitempty"`
	ViewType   ViewType `json:"view_type,omitempty"`
	MaxResults int64    `json:"max_results,omitempty"`
	OrderBy    []string `json:"order_by,omitempty"`
	PageToken  string   `json:"page_token,omitempty"`
}

type ExperimentsSearchResults struct {
	Experiments   []*Experiment `json:"experiments,omitempty"`
	NextPageToken string        `json:"next_page_token,omitempty"`
}

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

func (s *ExperimentService) Delete(ctx context.Context, id string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/delete", nil, &opts, nil)
	return err
}

func (s *ExperimentService) Restore(ctx context.Context, id string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	_, err := s.client.Do(ctx, "POST", "experiments/restore", nil, &opts, nil)
	return err
}

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

func (s *ExperimentService) Get(ctx context.Context, id string) (*Experiment, error) {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	var res struct {
		Experiment *Experiment `json:"experiment,omitempty"`
	}

	_, err := s.client.Do(ctx, "GET", "experiments/get", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Experiment, nil
}

func (s *ExperimentService) GetByName(ctx context.Context, name string) (*Experiment, error) {
	opts := struct {
		ExperimentName string `json:"experiment_name,omitempty"`
	}{
		ExperimentName: name,
	}

	var res struct {
		Experiment *Experiment `json:"experiment,omitempty"`
	}

	_, err := s.client.Do(ctx, "GET", "experiments/get-by-name", nil, &opts, &res)
	if err != nil {
		return nil, err
	}

	return res.Experiment, nil
}

func (s *ExperimentService) Search(ctx context.Context, opts *ExperimentsSearchOptions) (*ExperimentsSearchResults, error) {
	var res ExperimentsSearchResults

	_, err := s.client.Do(ctx, "POST", "experiments/search", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

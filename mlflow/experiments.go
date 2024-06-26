package mlflow

import "context"

type ExperimentService service

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

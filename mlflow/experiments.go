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

	_, err := s.client.Do(ctx, "POST", "/experiments/create", &opts, &res)
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

	_, err := s.client.Do(ctx, "POST", "/experiments/update", &opts, nil)
	return err
}

func (s *ExperimentService) Delete(ctx context.Context, id string) error {
	opts := struct {
		ExperimentID string `json:"experiment_id,omitempty"`
	}{
		ExperimentID: id,
	}

	_, err := s.client.Do(ctx, "POST", "/experiments/delete", &opts, nil)
	return err
}

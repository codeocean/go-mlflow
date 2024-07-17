package mlflow

type ModelVersionStatus string

const (
	ModelVersionStatusPending ModelVersionStatus = "PENDING_REGISTRATION"
	ModelVersionStatusFailed  ModelVersionStatus = "FAILED_REGISTRATION"
	ModelVersionStatusReady   ModelVersionStatus = "READY"
)

type ModelVersion struct {
	Name                 string             `json:"name,omitempty"`
	Version              string             `json:"version,omitempty"`
	CreationTimestamp    int64              `json:"creation_timestamp,omitempty"`
	LastUpdatedTimestamp int64              `json:"last_updated_timestamp,omitempty"`
	UserID               string             `json:"user_id,omitempty"`
	CurrentStage         string             `json:"current_stage,omitempty"`
	Description          string             `json:"description,omitempty"`
	Source               string             `json:"source,omitempty"`
	RunID                string             `json:"run_id,omitempty"`
	Status               ModelVersionStatus `json:"status,omitempty"`
	StatusMessage        string             `json:"status_message,omitempty"`
	Tags                 []*ModelVersionTag `json:"tags,omitempty"`
	RunLink              string             `json:"run_link,omitempty"`
	Aliases              []string           `json:"aliases,omitempty"`
}

type ModelVersionTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

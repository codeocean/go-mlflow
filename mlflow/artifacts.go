package mlflow

import "context"

type ArtifactsService service

type ListArtifactsRequest struct {
	RunID     string `json:"run_id,omitempty" binding:"required"`
	Path      string `json:"path,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

type ListArtifactsResponse struct {
	RootURI   string      `json:"root_uri,omitempty"`
	Files     []*FileInfo `json:"files,omitempty"`
	NextToken string      `json:"next_page_token,omitempty"`
}

type FileInfo struct {
	Path     string `json:"path,omitempty"`
	IsDir    bool   `json:"is_dir,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
}

func (s *ArtifactsService) List(ctx context.Context, opts *ListArtifactsRequest) (*ListArtifactsResponse, error) {
	var res ListArtifactsResponse

	_, err := s.client.Do(ctx, "GET", "artifacts/list", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

package mlflow

import "context"

// ArtifactsService handles communication with the artifacts related methods of the MLflow API.
type ArtifactsService service

// ListArtifactsRequest represents the request parameters for listing artifacts
type ListArtifactsRequest struct {
	// RunID is the unique identifier for the run (required)
	RunID string `json:"run_id,omitempty" binding:"required"`
	// Path is the relative path within the artifact root to list (optional, defaults to root)
	Path string `json:"path,omitempty"`
	// PageToken is used for pagination to fetch the next page of results
	PageToken string `json:"page_token,omitempty"`
}

// ListArtifactsResponse represents the response from listing artifacts
type ListArtifactsResponse struct {
	// RootURI is the root location for storing artifacts
	RootURI string `json:"root_uri,omitempty"`
	// Files is the list of artifacts at the specified path
	Files []*FileInfo `json:"files,omitempty"`
	// NextPageToken is used to retrieve the next page of results (empty if no more results)
	NextPageToken string `json:"next_page_token,omitempty"`
}

// FileInfo represents metadata about a file or directory artifact
type FileInfo struct {
	// Path is the relative path from the artifact root
	Path string `json:"path,omitempty"`
	// IsDir indicates whether this entry is a directory
	IsDir bool `json:"is_dir,omitempty"`
	// FileSize is the size of the file in bytes (0 for directories)
	FileSize int64 `json:"file_size,omitempty"`
}

// List retrieves the list of artifacts for a run
//
// Artifacts can be filtered by specifying a path. Use the PageToken for pagination
// when dealing with large numbers of artifacts.
func (s *ArtifactsService) List(ctx context.Context, opts *ListArtifactsRequest) (*ListArtifactsResponse, error) {
	var res ListArtifactsResponse

	_, err := s.client.Do(ctx, "GET", "artifacts/list", nil, opts, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

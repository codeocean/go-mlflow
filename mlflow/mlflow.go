package mlflow

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client manages communication with the MLflow API.
//
// Create a new Client using NewClient, then use the various services to access
// different parts of the MLflow REST API.
type Client struct {
	client  *http.Client
	baseURL *url.URL

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the MLflow API.

	// Artifacts provides access to artifact-related operations
	Artifacts *ArtifactsService
	// Experiments provides access to experiment-related operations
	Experiments *ExperimentService
	// LoggedModels provides access to logged model operations
	LoggedModels *LoggedModelService
	// Metrics provides access to metric-related operations
	Metrics *MetricsService
	// ModelVersions provides access to model version operations in the Model Registry
	ModelVersions *ModelVersionService
	// RegisteredModels provides access to registered model operations in the Model Registry
	RegisteredModels *RegisteredModelService
	// Runs provides access to run-related operations
	Runs *RunService
	// Users provides access to user management operations
	Users *UserService
}

// service is the base struct for all service types.
//
// It provides access to the client for making API calls.
type service struct {
	client *Client
}

// NewClient creates a new MLflow API client
//
// The baseURL should point to your MLflow tracking server (e.g., "http://localhost:5000").
// If httpClient is nil, a default http.Client will be used.
//
// The client automatically appends "/api/2.0/mlflow/" to the base URL for API requests.
func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}
	parsedURL.Path += "api/2.0/mlflow/"

	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient

	c := &Client{
		client:  &httpClient2,
		baseURL: parsedURL,
	}

	c.common.client = c
	c.Artifacts = (*ArtifactsService)(&c.common)
	c.Experiments = (*ExperimentService)(&c.common)
	c.LoggedModels = (*LoggedModelService)(&c.common)
	c.Metrics = (*MetricsService)(&c.common)
	c.ModelVersions = (*ModelVersionService)(&c.common)
	c.RegisteredModels = (*RegisteredModelService)(&c.common)
	c.Runs = (*RunService)(&c.common)
	c.Users = (*UserService)(&c.common)

	return c, nil
}

// Do sends an API request and returns the API response
//
// This is the core method that all service methods use to communicate with the MLflow API.
// It handles request construction, error handling, and response parsing.
//
// Parameters:
//   - ctx: Context for the request
//   - method: HTTP method (GET, POST, PATCH, DELETE, etc.)
//   - path: API endpoint path relative to the base URL
//   - params: URL query parameters (can be nil)
//   - body: Request body to be JSON-encoded (can be nil)
//   - response: Pointer to struct for decoding the response (can be nil or io.Writer)
//
// Returns the raw http.Response and any error encountered.
func (c *Client) Do(ctx context.Context, method string, path string, params url.Values, body interface{}, response interface{}) (*http.Response, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	if params != nil {
		u.RawQuery = params.Encode()
	}

	bodyReader, err := c.encodeBody(body)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	req := r.WithContext(ctx)

	req.Header.Set("content-type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		e := Error{StatusCode: res.StatusCode}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, res.Body)
			e.Message = buf.String()
		}
		return res, &e
	}

	switch v := response.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		err = json.NewDecoder(res.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
	}

	return res, err
}

// encodeBody encodes the request body as JSON
//
// Returns nil if body is nil, otherwise returns a reader containing the JSON-encoded body.
func (c *Client) encodeBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

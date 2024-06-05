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

type Client struct {
	client  *http.Client
	baseURL string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the MLflow API.
	Experiments *ExperimentService
	Users       *UserService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient

	c := &Client{
		client:  &httpClient2,
		baseURL: parsedURL.String(),
	}

	c.common.client = c
	c.Experiments = (*ExperimentService)(&c.common)
	c.Users = (*UserService)(&c.common)

	return c, nil
}

func (c *Client) Do(ctx context.Context, method string, path string, body interface{}, response interface{}) (*http.Response, error) {
	bodyReader, err := c.encodeBody(body)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req := r.WithContext(ctx)

	req.Header.Set("content-type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return res, err
	}
	res.Body.Close()

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

// QueryParams defines the query parameters to be sent when calling APIs
type QueryParams map[string]string

func (m QueryParams) String() string {
	var buffer bytes.Buffer

	for key, value := range m {
		if buffer.Len() > 0 {
			buffer.WriteByte('&')
		}
		buffer.WriteString(key)
		buffer.WriteByte('=')
		buffer.WriteString(url.QueryEscape(value))
	}

	return buffer.String()
}

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
	baseURL *url.URL

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
		baseURL: parsedURL,
	}

	c.common.client = c
	c.Experiments = (*ExperimentService)(&c.common)
	c.Users = (*UserService)(&c.common)

	return c, nil
}

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

package utils

import (
	"fmt"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func SimpleGetRequest(url string) (string, error){
	resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return "", err
    }
    defer resp.Body.Close() // Always close response body

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Read error:", err)
        return "", err
    }
	return string(body), nil
}

// APIClient is a reusable HTTP client for making API requests.
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
	ReturnRawResponse bool
}

// New creates a new APIClient with a base URL and default timeout.
func NewHttpClient(baseURL string, timeout time.Duration, headers map[string]string, returnRawResponse bool)  *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		Headers: headers,
		ReturnRawResponse: returnRawResponse,
	}
}

// doRequest is an internal helper that executes HTTP requests.
func (c *APIClient) doRequest(ctx context.Context, method, endpoint string, body any, result any) error {
	url := c.BaseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	// Add default headers
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(string(b))
	}
	if c.ReturnRawResponse {
		result = resp.Body
		return nil
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// Get sends a GET request.
func (c *APIClient) Get(ctx context.Context, endpoint string, result any) error {
	return c.doRequest(ctx, http.MethodGet, endpoint, nil, result)
}

// Post sends a POST request.
func (c *APIClient) Post(ctx context.Context, endpoint string, body any, result any) error {
	return c.doRequest(ctx, http.MethodPost, endpoint, body, result)
}

// Put sends a PUT request.
func (c *APIClient) Put(ctx context.Context, endpoint string, body any, result any) error {
	return c.doRequest(ctx, http.MethodPut, endpoint, body, result)
}

// Delete sends a DELETE request.
func (c *APIClient) Delete(ctx context.Context, endpoint string, result any) error {
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil, result)
}

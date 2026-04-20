// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultBaseURL = "https://console.jumpcloud.com/api"
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
	pageSize       = 100
)

// JumpCloudClient is an HTTP client for the JumpCloud API.
type JumpCloudClient struct {
	BaseURL    string
	APIKey     string
	OrgID      string
	HTTPClient *http.Client
	UserAgent  string
}

// APIError represents an error response from the JumpCloud API.
type APIError struct {
	StatusCode int
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("JumpCloud API error (HTTP %d): %s", e.StatusCode, e.Message)
}

// NewJumpCloudClient creates a new JumpCloud API client.
func NewJumpCloudClient(apiKey, orgID, apiURL, version string) *JumpCloudClient {
	if apiURL == "" {
		apiURL = defaultBaseURL
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	retryClient.Logger = nil // suppress default logging; we use tflog

	httpClient := retryClient.StandardClient()
	httpClient.Timeout = defaultTimeout

	return &JumpCloudClient{
		BaseURL:    apiURL,
		APIKey:     apiKey,
		OrgID:      orgID,
		HTTPClient: httpClient,
		UserAgent:  fmt.Sprintf("terraform-provider-jumpcloud/%s", version),
	}
}

// doRequest executes an HTTP request against the JumpCloud API.
func (c *JumpCloudClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	reqURL := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if c.OrgID != "" {
		req.Header.Set("x-org-id", c.OrgID)
	}

	tflog.Debug(ctx, "JumpCloud API request", map[string]interface{}{
		"method": method,
		"url":    reqURL,
	})

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	tflog.Debug(ctx, "JumpCloud API response", map[string]interface{}{
		"status": resp.StatusCode,
	})

	if err := checkResponse(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

// doRequestWithQuery executes an HTTP GET request with query parameters.
func (c *JumpCloudClient) doRequestWithQuery(ctx context.Context, path string, params url.Values) ([]byte, error) {
	reqURL := c.BaseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if c.OrgID != "" {
		req.Header.Set("x-org-id", c.OrgID)
	}

	tflog.Debug(ctx, "JumpCloud API request", map[string]interface{}{
		"method": "GET",
		"url":    reqURL,
	})

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	tflog.Debug(ctx, "JumpCloud API response", map[string]interface{}{
		"status": resp.StatusCode,
	})

	if err := checkResponse(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

// doListRequest handles paginated GET requests and collects all results.
func (c *JumpCloudClient) doListRequest(ctx context.Context, path string, params url.Values) ([]json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}

	var allResults []json.RawMessage
	skip := 0

	for {
		params.Set("limit", strconv.Itoa(pageSize))
		params.Set("skip", strconv.Itoa(skip))

		body, err := c.doRequestWithQuery(ctx, path, params)
		if err != nil {
			return nil, err
		}

		// V1 API returns {"results": [...], "totalCount": N}
		// V2 API returns a JSON array directly
		var results []json.RawMessage

		// Try V2 format first (plain array)
		if err := json.Unmarshal(body, &results); err != nil {
			// Try V1 format
			var v1Response struct {
				Results    []json.RawMessage `json:"results"`
				TotalCount int               `json:"totalCount"`
			}
			if err := json.Unmarshal(body, &v1Response); err != nil {
				return nil, fmt.Errorf("parsing list response: %w", err)
			}
			results = v1Response.Results
		}

		allResults = append(allResults, results...)

		if len(results) < pageSize {
			break
		}

		skip += pageSize
	}

	return allResults, nil
}

// checkResponse checks the HTTP response status code and returns an error for non-2xx.
func checkResponse(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    string(body),
	}

	// Try to parse structured error message
	var errResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil {
		if errResp.Message != "" {
			apiErr.Message = errResp.Message
		} else if errResp.Error != "" {
			apiErr.Message = errResp.Error
		}
	}

	return apiErr
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

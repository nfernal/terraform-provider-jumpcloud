// Copyright (c) Nfernal
// SPDX-License-Identifier: MIT

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(server *httptest.Server) *JumpCloudClient {
	return &JumpCloudClient{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		OrgID:      "test-org-id",
		HTTPClient: server.Client(),
		UserAgent:  "terraform-provider-jumpcloud/test",
	}
}

func TestNewJumpCloudClient(t *testing.T) {
	c := NewJumpCloudClient("key", "org", "", "1.0.0")
	if c.BaseURL != defaultBaseURL {
		t.Errorf("expected default base URL %q, got %q", defaultBaseURL, c.BaseURL)
	}
	if c.APIKey != "key" {
		t.Errorf("expected API key %q, got %q", "key", c.APIKey)
	}
	if c.OrgID != "org" {
		t.Errorf("expected OrgID %q, got %q", "org", c.OrgID)
	}
	if c.UserAgent != "terraform-provider-jumpcloud/1.0.0" {
		t.Errorf("expected UserAgent %q, got %q", "terraform-provider-jumpcloud/1.0.0", c.UserAgent)
	}
}

func TestNewJumpCloudClient_customURL(t *testing.T) {
	c := NewJumpCloudClient("key", "org", "https://custom.api.com", "2.0.0")
	if c.BaseURL != "https://custom.api.com" {
		t.Errorf("expected custom base URL, got %q", c.BaseURL)
	}
}

func TestDoRequest_setsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("expected x-api-key %q, got %q", "test-api-key", got)
		}
		if got := r.Header.Get("x-org-id"); got != "test-org-id" {
			t.Errorf("expected x-org-id %q, got %q", "test-org-id", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type %q, got %q", "application/json", got)
		}
		if got := r.Header.Get("User-Agent"); got != "terraform-provider-jumpcloud/test" {
			t.Errorf("expected User-Agent %q, got %q", "terraform-provider-jumpcloud/test", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_noOrgID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-org-id"); got != "" {
			t.Errorf("expected no x-org-id header, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	c.OrgID = ""
	_, err := c.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_sendsBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["name"] != "test" {
			t.Errorf("expected body name %q, got %q", "test", body["name"])
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"id":"123"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	resp, err := c.doRequest(context.Background(), http.MethodPost, "/test", map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"id":"123"}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestDoRequest_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, `{"message":"bad request"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("expected message %q, got %q", "bad request", apiErr.Message)
	}
}

func TestDoRequest_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"error":"internal server error"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
	if apiErr.Message != "internal server error" {
		t.Errorf("expected message %q, got %q", "internal server error", apiErr.Message)
	}
}

func TestDoRequestWithQuery_addsParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("filter"); got != "name:eq:test" {
			t.Errorf("expected filter param %q, got %q", "name:eq:test", got)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("expected limit param %q, got %q", "10", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	params := map[string][]string{
		"filter": {"name:eq:test"},
		"limit":  {"10"},
	}
	_, err := c.doRequestWithQuery(context.Background(), "/test", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequestWithQuery_noParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %q", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doRequestWithQuery(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoListRequest_v2Format(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `[{"id":"1"},{"id":"2"}]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	results, err := c.doListRequest(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestDoListRequest_v1Format(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"results":[{"id":"1"},{"id":"2"},{"id":"3"}],"totalCount":3}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	results, err := c.doListRequest(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestDoListRequest_pagination(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		// First call: return exactly pageSize items to trigger pagination
		if callCount == 1 {
			items := make([]map[string]string, pageSize)
			for i := range items {
				items[i] = map[string]string{"id": fmt.Sprintf("%d", i)}
			}
			_ = json.NewEncoder(w).Encode(items)
		} else {
			// Second call: return fewer than pageSize to stop
			_, _ = fmt.Fprint(w, `[{"id":"extra"}]`)
		}
	}))
	defer server.Close()

	c := newTestClient(server)
	results, err := c.doListRequest(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != pageSize+1 {
		t.Errorf("expected %d results, got %d", pageSize+1, len(results))
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls for pagination, got %d", callCount)
	}
}

func TestDoListRequest_emptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	results, err := c.doListRequest(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDoListRequest_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, `{"message":"forbidden"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doListRequest(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDoListRequest_invalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `not json`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.doListRequest(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestCheckResponse_success(t *testing.T) {
	for _, code := range []int{200, 201, 204, 299} {
		if err := checkResponse(code, []byte(`{}`)); err != nil {
			t.Errorf("expected no error for status %d, got %v", code, err)
		}
	}
}

func TestCheckResponse_errorWithMessage(t *testing.T) {
	err := checkResponse(400, []byte(`{"message":"validation failed"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.Message != "validation failed" {
		t.Errorf("expected message %q, got %q", "validation failed", apiErr.Message)
	}
}

func TestCheckResponse_errorWithErrorField(t *testing.T) {
	err := checkResponse(500, []byte(`{"error":"something broke"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.Message != "something broke" {
		t.Errorf("expected message %q, got %q", "something broke", apiErr.Message)
	}
}

func TestCheckResponse_errorWithUnstructuredBody(t *testing.T) {
	err := checkResponse(502, []byte(`Bad Gateway`))
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.Message != "Bad Gateway" {
		t.Errorf("expected message %q, got %q", "Bad Gateway", apiErr.Message)
	}
}

func TestCheckResponse_messagePreferredOverError(t *testing.T) {
	err := checkResponse(400, []byte(`{"message":"specific","error":"generic"}`))
	apiErr := err.(*APIError)
	if apiErr.Message != "specific" {
		t.Errorf("expected message field to take precedence, got %q", apiErr.Message)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	expected := "JumpCloud API error (HTTP 404): not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"404 error", &APIError{StatusCode: 404, Message: "not found"}, true},
		{"400 error", &APIError{StatusCode: 400, Message: "bad request"}, false},
		{"500 error", &APIError{StatusCode: 500, Message: "server error"}, false},
		{"non-API error", fmt.Errorf("some error"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotFound(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

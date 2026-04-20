// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SystemGroup represents a JumpCloud V2 system group.
type SystemGroup struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// CreateSystemGroup creates a new system group.
func (c *JumpCloudClient) CreateSystemGroup(ctx context.Context, group SystemGroup) (*SystemGroup, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/systemgroups", group)
	if err != nil {
		return nil, fmt.Errorf("creating system group: %w", err)
	}

	var result SystemGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing system group response: %w", err)
	}

	return &result, nil
}

// GetSystemGroup retrieves a system group by ID.
func (c *JumpCloudClient) GetSystemGroup(ctx context.Context, id string) (*SystemGroup, error) {
	body, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/systemgroups/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("getting system group: %w", err)
	}

	var result SystemGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing system group response: %w", err)
	}

	return &result, nil
}

// UpdateSystemGroup updates an existing system group.
func (c *JumpCloudClient) UpdateSystemGroup(ctx context.Context, id string, group SystemGroup) (*SystemGroup, error) {
	body, err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/v2/systemgroups/%s", id), group)
	if err != nil {
		return nil, fmt.Errorf("updating system group: %w", err)
	}

	var result SystemGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing system group response: %w", err)
	}

	return &result, nil
}

// DeleteSystemGroup deletes a system group by ID.
func (c *JumpCloudClient) DeleteSystemGroup(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/v2/systemgroups/%s", id), nil)
	if err != nil {
		return fmt.Errorf("deleting system group: %w", err)
	}

	return nil
}

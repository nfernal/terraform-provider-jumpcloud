// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// UserGroup represents a JumpCloud V2 user group.
type UserGroup struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// GroupMemberResponse represents the API response for a group member.
type GroupMemberResponse struct {
	To GroupMember `json:"to"`
}

// GroupMember represents a member of a JumpCloud group.
type GroupMember struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// GroupMembershipOp represents an add/remove membership operation.
type GroupMembershipOp struct {
	Op   string `json:"op"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

// CreateUserGroup creates a new user group.
func (c *JumpCloudClient) CreateUserGroup(ctx context.Context, group UserGroup) (*UserGroup, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/v2/usergroups", group)
	if err != nil {
		return nil, fmt.Errorf("creating user group: %w", err)
	}

	var result UserGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user group response: %w", err)
	}

	return &result, nil
}

// GetUserGroup retrieves a user group by ID.
func (c *JumpCloudClient) GetUserGroup(ctx context.Context, id string) (*UserGroup, error) {
	body, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/usergroups/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("getting user group: %w", err)
	}

	var result UserGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user group response: %w", err)
	}

	return &result, nil
}

// UpdateUserGroup updates an existing user group.
func (c *JumpCloudClient) UpdateUserGroup(ctx context.Context, id string, group UserGroup) (*UserGroup, error) {
	body, err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("/v2/usergroups/%s", id), group)
	if err != nil {
		return nil, fmt.Errorf("updating user group: %w", err)
	}

	var result UserGroup
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user group response: %w", err)
	}

	return &result, nil
}

// DeleteUserGroup deletes a user group by ID.
func (c *JumpCloudClient) DeleteUserGroup(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/v2/usergroups/%s", id), nil)
	if err != nil {
		return fmt.Errorf("deleting user group: %w", err)
	}

	return nil
}

// GetUserGroupByName looks up a user group by name.
func (c *JumpCloudClient) GetUserGroupByName(ctx context.Context, name string) (*UserGroup, error) {
	params := url.Values{}
	params.Set("filter", fmt.Sprintf("name:eq:%s", name))
	params.Set("limit", "1")

	body, err := c.doRequestWithQuery(ctx, "/v2/usergroups", params)
	if err != nil {
		return nil, fmt.Errorf("searching user group by name: %w", err)
	}

	var results []UserGroup
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("parsing user group search response: %w", err)
	}

	if len(results) == 0 {
		return nil, &APIError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("no user group found with name: %s", name),
		}
	}

	for _, group := range results {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, &APIError{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("no user group found with exact name: %s", name),
	}
}

// ListUserGroupMembers lists members of a user group.
func (c *JumpCloudClient) ListUserGroupMembers(ctx context.Context, groupID string) ([]GroupMember, error) {
	results, err := c.doListRequest(ctx, fmt.Sprintf("/v2/usergroups/%s/members", groupID), nil)
	if err != nil {
		return nil, fmt.Errorf("listing user group members: %w", err)
	}

	var members []GroupMember
	for _, raw := range results {
		var resp GroupMemberResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, fmt.Errorf("parsing group member: %w", err)
		}
		members = append(members, resp.To)
	}

	return members, nil
}

// ModifyUserGroupMembership adds or removes a member from a user group.
func (c *JumpCloudClient) ModifyUserGroupMembership(ctx context.Context, groupID string, op GroupMembershipOp) error {
	_, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/v2/usergroups/%s/members", groupID), op)
	if err != nil {
		return fmt.Errorf("modifying user group membership: %w", err)
	}

	return nil
}

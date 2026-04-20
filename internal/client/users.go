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

// SystemUser represents a JumpCloud V1 system user.
type SystemUser struct {
	ID                          string `json:"_id,omitempty"`
	Username                    string `json:"username"`
	Email                       string `json:"email"`
	Firstname                   string `json:"firstname,omitempty"`
	Lastname                    string `json:"lastname,omitempty"`
	Password                    string `json:"password,omitempty"`
	AccountLocked               *bool  `json:"account_locked,omitempty"`
	Activated                   bool   `json:"activated"`
	AllowPublicKey              *bool  `json:"allow_public_key,omitempty"`
	Sudo                        *bool  `json:"sudo,omitempty"`
	EnableManagedUID            *bool  `json:"enable_managed_uid,omitempty"`
	UnixUID                     *int   `json:"unix_uid,omitempty"`
	UnixGUID                    *int   `json:"unix_guid,omitempty"`
	PasswordlessSudo            *bool  `json:"passwordless_sudo,omitempty"`
	ExternallyManaged           *bool  `json:"externally_managed,omitempty"`
	LdapBindingUser             *bool  `json:"ldap_binding_user,omitempty"`
	EnableUserPortalMultifactor *bool  `json:"enable_user_portal_multifactor,omitempty"`
	Suspended                   *bool  `json:"suspended,omitempty"`
	State                       string `json:"state,omitempty"`
	Description                 string `json:"description,omitempty"`
	Department                  string `json:"department,omitempty"`
	CostCenter                  string `json:"costCenter,omitempty"`
	Company                     string `json:"company,omitempty"`
	EmployeeType                string `json:"employeeType,omitempty"`
	EmployeeIdentifier          string `json:"employeeIdentifier,omitempty"`
	JobTitle                    string `json:"jobTitle,omitempty"`
	Location                    string `json:"location,omitempty"`
	Manager                     string `json:"manager,omitempty"`
	MiddleName                  string `json:"middlename,omitempty"`
	DisplayName                 string `json:"displayname,omitempty"`
	ExternalDN                  string `json:"external_dn,omitempty"`
	ExternalSourceType          string `json:"external_source_type,omitempty"`
	MFA                         *MFA   `json:"mfa,omitempty"`
}

// MFA represents the MFA configuration for a system user.
type MFA struct {
	Configured *bool `json:"configured,omitempty"`
	Exclusion  *bool `json:"exclusion,omitempty"`
}

// CreateUser creates a new system user in JumpCloud.
func (c *JumpCloudClient) CreateUser(ctx context.Context, user SystemUser) (*SystemUser, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/systemusers", user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	var result SystemUser
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user response: %w", err)
	}

	return &result, nil
}

// GetUser retrieves a system user by ID.
func (c *JumpCloudClient) GetUser(ctx context.Context, id string) (*SystemUser, error) {
	body, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/systemusers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	var result SystemUser
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user response: %w", err)
	}

	return &result, nil
}

// UpdateUser updates an existing system user.
func (c *JumpCloudClient) UpdateUser(ctx context.Context, id string, user SystemUser) (*SystemUser, error) {
	body, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/systemusers/%s", id), user)
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	var result SystemUser
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing user response: %w", err)
	}

	return &result, nil
}

// DeleteUser deletes a system user by ID.
func (c *JumpCloudClient) DeleteUser(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/systemusers/%s", id), nil)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return nil
}

// GetUserByEmail looks up a system user by email address.
func (c *JumpCloudClient) GetUserByEmail(ctx context.Context, email string) (*SystemUser, error) {
	params := url.Values{}
	params.Set("filter", fmt.Sprintf("email:$eq:%s", email))
	params.Set("limit", "1")

	body, err := c.doRequestWithQuery(ctx, "/systemusers", params)
	if err != nil {
		return nil, fmt.Errorf("searching user by email: %w", err)
	}

	var response struct {
		Results    []SystemUser `json:"results"`
		TotalCount int          `json:"totalCount"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, &APIError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("no user found with email: %s", email),
		}
	}

	// Verify exact match
	for _, user := range response.Results {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, &APIError{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("no user found with exact email: %s", email),
	}
}

// GetUserByUsername looks up a system user by username.
func (c *JumpCloudClient) GetUserByUsername(ctx context.Context, username string) (*SystemUser, error) {
	params := url.Values{}
	params.Set("filter", fmt.Sprintf("username:$eq:%s", username))
	params.Set("limit", "1")

	body, err := c.doRequestWithQuery(ctx, "/systemusers", params)
	if err != nil {
		return nil, fmt.Errorf("searching user by username: %w", err)
	}

	var response struct {
		Results    []SystemUser `json:"results"`
		TotalCount int          `json:"totalCount"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, &APIError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("no user found with username: %s", username),
		}
	}

	// Verify exact match
	for _, user := range response.Results {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, &APIError{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("no user found with exact username: %s", username),
	}
}

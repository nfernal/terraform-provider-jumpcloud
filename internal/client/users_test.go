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

func TestCreateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/systemusers" {
			t.Errorf("expected path /systemusers, got %s", r.URL.Path)
		}

		var user SystemUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("expected username %q, got %q", "testuser", user.Username)
		}
		if user.Email != "test@example.com" {
			t.Errorf("expected email %q, got %q", "test@example.com", user.Email)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SystemUser{
			ID:       "user123",
			Username: user.Username,
			Email:    user.Email,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.CreateUser(context.Background(), SystemUser{
		Username: "testuser",
		Email:    "test@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user123" {
		t.Errorf("expected ID %q, got %q", "user123", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username %q, got %q", "testuser", user.Username)
	}
}

func TestCreateUser_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"message":"user already exists"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.CreateUser(context.Background(), SystemUser{
		Username: "testuser",
		Email:    "test@example.com",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateUser_withAllFields(t *testing.T) {
	boolTrue := true
	intVal := 5000

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user SystemUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if user.Sudo == nil || !*user.Sudo {
			t.Error("expected sudo to be true")
		}
		if user.UnixUID == nil || *user.UnixUID != 5000 {
			t.Error("expected unix_uid to be 5000")
		}
		if user.MFA == nil || user.MFA.Configured == nil || !*user.MFA.Configured {
			t.Error("expected MFA configured to be true")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.CreateUser(context.Background(), SystemUser{
		Username:  "fulluser",
		Email:     "full@example.com",
		Firstname: "Full",
		Lastname:  "User",
		Sudo:      &boolTrue,
		UnixUID:   &intVal,
		MFA:       &MFA{Configured: &boolTrue},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/systemusers/user123" {
			t.Errorf("expected path /systemusers/user123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SystemUser{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.GetUser(context.Background(), "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user123" {
		t.Errorf("expected ID %q, got %q", "user123", user.ID)
	}
}

func TestGetUser_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"user not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUser(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// GetUser wraps the error, so check the error message
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestUpdateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/systemusers/user123" {
			t.Errorf("expected path /systemusers/user123, got %s", r.URL.Path)
		}

		var user SystemUser
		json.NewDecoder(r.Body).Decode(&user)

		w.WriteHeader(http.StatusOK)
		user.ID = "user123"
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.UpdateUser(context.Background(), "user123", SystemUser{
		Username:  "testuser",
		Email:     "test@example.com",
		Firstname: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Firstname != "Updated" {
		t.Errorf("expected firstname %q, got %q", "Updated", user.Firstname)
	}
}

func TestUpdateUser_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"user not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.UpdateUser(context.Background(), "nonexistent", SystemUser{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/systemusers/user123" {
			t.Errorf("expected path /systemusers/user123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteUser(context.Background(), "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUser_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"user not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteUser(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserByEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		if filter != "email:$eq:test@example.com" {
			t.Errorf("expected filter %q, got %q", "email:$eq:test@example.com", filter)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []SystemUser{
				{ID: "user123", Username: "testuser", Email: "test@example.com"},
			},
			"totalCount": 1,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.GetUserByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user123" {
		t.Errorf("expected ID %q, got %q", "user123", user.ID)
	}
}

func TestGetUserByEmail_noResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results":    []SystemUser{},
			"totalCount": 0,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserByEmail(context.Background(), "notfound@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetUserByEmail_noExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []SystemUser{
				{ID: "user123", Username: "testuser", Email: "other@example.com"},
			},
			"totalCount": 1,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserByEmail(context.Background(), "test@example.com")
	if err == nil {
		t.Fatal("expected error for non-exact match, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		if filter != "username:$eq:testuser" {
			t.Errorf("expected filter %q, got %q", "username:$eq:testuser", filter)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []SystemUser{
				{ID: "user123", Username: "testuser", Email: "test@example.com"},
			},
			"totalCount": 1,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.GetUserByUsername(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user123" {
		t.Errorf("expected ID %q, got %q", "user123", user.ID)
	}
}

func TestGetUserByUsername_noResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results":    []SystemUser{},
			"totalCount": 0,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserByUsername(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetUserByUsername_noExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []SystemUser{
				{ID: "user123", Username: "testuser_other", Email: "test@example.com"},
			},
			"totalCount": 1,
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserByUsername(context.Background(), "testuser")
	if err == nil {
		t.Fatal("expected error for non-exact match, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetUser_withNilPointerFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Return user with no optional pointer fields
		fmt.Fprint(w, `{"_id":"user123","username":"testuser","email":"test@example.com","activated":true}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.GetUser(context.Background(), "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Sudo != nil {
		t.Error("expected Sudo to be nil")
	}
	if user.AccountLocked != nil {
		t.Error("expected AccountLocked to be nil")
	}
	if user.UnixUID != nil {
		t.Error("expected UnixUID to be nil")
	}
	if user.MFA != nil {
		t.Error("expected MFA to be nil")
	}
}

func TestGetUser_withMFAFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"_id":"user123","username":"testuser","email":"test@example.com","mfa":{"configured":true,"exclusion":false}}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	user, err := c.GetUser(context.Background(), "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.MFA == nil {
		t.Fatal("expected MFA to be non-nil")
	}
	if user.MFA.Configured == nil || !*user.MFA.Configured {
		t.Error("expected MFA.Configured to be true")
	}
	if user.MFA.Exclusion == nil || *user.MFA.Exclusion {
		t.Error("expected MFA.Exclusion to be false")
	}
}

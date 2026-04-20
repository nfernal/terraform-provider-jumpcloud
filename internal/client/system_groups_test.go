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

func TestCreateSystemGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/systemgroups" {
			t.Errorf("expected path /v2/systemgroups, got %s", r.URL.Path)
		}

		var group SystemGroup
		json.NewDecoder(r.Body).Decode(&group)

		w.WriteHeader(http.StatusCreated)
		group.ID = "sysgrp123"
		json.NewEncoder(w).Encode(group)
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.CreateSystemGroup(context.Background(), SystemGroup{
		Name:        "test-system-group",
		Description: "test desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "sysgrp123" {
		t.Errorf("expected ID %q, got %q", "sysgrp123", group.ID)
	}
	if group.Name != "test-system-group" {
		t.Errorf("expected name %q, got %q", "test-system-group", group.Name)
	}
}

func TestCreateSystemGroup_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"invalid system group"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.CreateSystemGroup(context.Background(), SystemGroup{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSystemGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/systemgroups/sysgrp123" {
			t.Errorf("expected path /v2/systemgroups/sysgrp123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SystemGroup{
			ID:          "sysgrp123",
			Name:        "test-system-group",
			Description: "desc",
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.GetSystemGroup(context.Background(), "sysgrp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "test-system-group" {
		t.Errorf("expected name %q, got %q", "test-system-group", group.Name)
	}
	if group.Description != "desc" {
		t.Errorf("expected description %q, got %q", "desc", group.Description)
	}
}

func TestGetSystemGroup_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"system group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetSystemGroup(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateSystemGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v2/systemgroups/sysgrp123" {
			t.Errorf("expected path /v2/systemgroups/sysgrp123, got %s", r.URL.Path)
		}

		var group SystemGroup
		json.NewDecoder(r.Body).Decode(&group)

		w.WriteHeader(http.StatusOK)
		group.ID = "sysgrp123"
		json.NewEncoder(w).Encode(group)
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.UpdateSystemGroup(context.Background(), "sysgrp123", SystemGroup{
		Name:        "updated-group",
		Description: "new desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Description != "new desc" {
		t.Errorf("expected description %q, got %q", "new desc", group.Description)
	}
}

func TestUpdateSystemGroup_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"system group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.UpdateSystemGroup(context.Background(), "nonexistent", SystemGroup{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteSystemGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v2/systemgroups/sysgrp123" {
			t.Errorf("expected path /v2/systemgroups/sysgrp123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteSystemGroup(context.Background(), "sysgrp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteSystemGroup_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"system group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteSystemGroup(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateSystemGroup_withDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var group SystemGroup
		json.NewDecoder(r.Body).Decode(&group)

		if group.Description != "detailed description" {
			t.Errorf("expected description %q, got %q", "detailed description", group.Description)
		}

		w.WriteHeader(http.StatusCreated)
		group.ID = "sysgrp456"
		json.NewEncoder(w).Encode(group)
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.CreateSystemGroup(context.Background(), SystemGroup{
		Name:        "desc-group",
		Description: "detailed description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Description != "detailed description" {
		t.Errorf("expected description %q, got %q", "detailed description", group.Description)
	}
}

// Copyright (c) Nfernal
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/usergroups" {
			t.Errorf("expected path /v2/usergroups, got %s", r.URL.Path)
		}

		var group UserGroup
		json.NewDecoder(r.Body).Decode(&group)

		w.WriteHeader(http.StatusCreated)
		group.ID = "grp123"
		json.NewEncoder(w).Encode(group)
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.CreateUserGroup(context.Background(), UserGroup{
		Name:        "test-group",
		Description: "test description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "grp123" {
		t.Errorf("expected ID %q, got %q", "grp123", group.ID)
	}
	if group.Name != "test-group" {
		t.Errorf("expected name %q, got %q", "test-group", group.Name)
	}
}

func TestCreateUserGroup_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"invalid group"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.CreateUserGroup(context.Background(), UserGroup{Name: "bad"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/usergroups/grp123" {
			t.Errorf("expected path /v2/usergroups/grp123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserGroup{
			ID:          "grp123",
			Name:        "test-group",
			Description: "desc",
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.GetUserGroup(context.Background(), "grp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "test-group" {
		t.Errorf("expected name %q, got %q", "test-group", group.Name)
	}
}

func TestGetUserGroup_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserGroup(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateUserGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v2/usergroups/grp123" {
			t.Errorf("expected path /v2/usergroups/grp123, got %s", r.URL.Path)
		}

		var group UserGroup
		json.NewDecoder(r.Body).Decode(&group)

		w.WriteHeader(http.StatusOK)
		group.ID = "grp123"
		json.NewEncoder(w).Encode(group)
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.UpdateUserGroup(context.Background(), "grp123", UserGroup{
		Name:        "updated-group",
		Description: "updated desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Description != "updated desc" {
		t.Errorf("expected description %q, got %q", "updated desc", group.Description)
	}
}

func TestDeleteUserGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v2/usergroups/grp123" {
			t.Errorf("expected path /v2/usergroups/grp123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteUserGroup(context.Background(), "grp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUserGroup_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.DeleteUserGroup(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserGroupByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		if filter != "name:eq:test-group" {
			t.Errorf("expected filter %q, got %q", "name:eq:test-group", filter)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]UserGroup{
			{ID: "grp123", Name: "test-group", Description: "found"},
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	group, err := c.GetUserGroupByName(context.Background(), "test-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "grp123" {
		t.Errorf("expected ID %q, got %q", "grp123", group.ID)
	}
}

func TestGetUserGroupByName_noResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserGroupByName(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetUserGroupByName_noExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]UserGroup{
			{ID: "grp123", Name: "test-group-other"},
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.GetUserGroupByName(context.Background(), "test-group")
	if err == nil {
		t.Fatal("expected error for non-exact match, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestListUserGroupMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/usergroups/grp123/members" {
			t.Errorf("expected path /v2/usergroups/grp123/members, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupMemberResponse{
			{To: GroupMember{ID: "user1", Type: "user"}},
			{To: GroupMember{ID: "user2", Type: "user"}},
		})
	}))
	defer server.Close()

	c := newTestClient(server)
	members, err := c.ListUserGroupMembers(context.Background(), "grp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
	if members[0].ID != "user1" {
		t.Errorf("expected first member ID %q, got %q", "user1", members[0].ID)
	}
	if members[1].ID != "user2" {
		t.Errorf("expected second member ID %q, got %q", "user2", members[1].ID)
	}
}

func TestListUserGroupMembers_empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	c := newTestClient(server)
	members, err := c.ListUserGroupMembers(context.Background(), "grp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 0 {
		t.Errorf("expected 0 members, got %d", len(members))
	}
}

func TestListUserGroupMembers_notFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"group not found"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.ListUserGroupMembers(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModifyUserGroupMembership_add(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/usergroups/grp123/members" {
			t.Errorf("expected path /v2/usergroups/grp123/members, got %s", r.URL.Path)
		}

		var op GroupMembershipOp
		json.NewDecoder(r.Body).Decode(&op)
		if op.Op != "add" {
			t.Errorf("expected op %q, got %q", "add", op.Op)
		}
		if op.Type != "user" {
			t.Errorf("expected type %q, got %q", "user", op.Type)
		}
		if op.ID != "user456" {
			t.Errorf("expected ID %q, got %q", "user456", op.ID)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.ModifyUserGroupMembership(context.Background(), "grp123", GroupMembershipOp{
		Op:   "add",
		Type: "user",
		ID:   "user456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestModifyUserGroupMembership_remove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var op GroupMembershipOp
		json.NewDecoder(r.Body).Decode(&op)
		if op.Op != "remove" {
			t.Errorf("expected op %q, got %q", "remove", op.Op)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.ModifyUserGroupMembership(context.Background(), "grp123", GroupMembershipOp{
		Op:   "remove",
		Type: "user",
		ID:   "user456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestModifyUserGroupMembership_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"invalid operation"}`)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.ModifyUserGroupMembership(context.Background(), "grp123", GroupMembershipOp{
		Op:   "add",
		Type: "user",
		ID:   "user456",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListUserGroupMembers_pagination(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)

		if callCount == 1 {
			// Return pageSize members to trigger pagination
			members := make([]GroupMemberResponse, pageSize)
			for i := range members {
				members[i] = GroupMemberResponse{
					To: GroupMember{ID: fmt.Sprintf("user%d", i), Type: "user"},
				}
			}
			json.NewEncoder(w).Encode(members)
		} else {
			// Return fewer to stop pagination
			json.NewEncoder(w).Encode([]GroupMemberResponse{
				{To: GroupMember{ID: "extra", Type: "user"}},
			})
		}
	}))
	defer server.Close()

	c := newTestClient(server)
	members, err := c.ListUserGroupMembers(context.Background(), "grp123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != pageSize+1 {
		t.Errorf("expected %d members, got %d", pageSize+1, len(members))
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTag_CRUD(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "tag_crud@test.com", "password123")

	// CREATE
	w := performRequest("POST", "/v1/tag/",
		map[string]string{"name": "Test Tag"}, accessToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// LIST
	w = performRequest("GET", "/v1/tag/", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	items := resp["response"].([]interface{})
	if len(items) == 0 {
		t.Fatal("list: expected at least one tag")
	}

	lastItem := items[len(items)-1].(map[string]interface{})
	tagID := lastItem["_id"].(string)

	// FIND
	w = performRequest("GET", fmt.Sprintf("/v1/tag/%s", tagID), nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("find: expected 200, got %d", w.Code)
	}

	// UPDATE
	w = performRequest("PATCH", fmt.Sprintf("/v1/tag/%s", tagID),
		map[string]string{"name": "Updated Tag"}, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify update
	w = performRequest("GET", fmt.Sprintf("/v1/tag/%s", tagID), nil, "")
	findResp := parseJSON(w)
	found := findResp["response"].(map[string]interface{})
	if found["name"] != "Updated Tag" {
		t.Errorf("expected 'Updated Tag', got '%v'", found["name"])
	}

	// DELETE
	w = performRequest("DELETE", fmt.Sprintf("/v1/tag/%s", tagID), nil, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", w.Code)
	}

	// Verify delete
	w = performRequest("GET", fmt.Sprintf("/v1/tag/%s", tagID), nil, "")
	if w.Code != http.StatusNotFound {
		t.Errorf("find after delete: expected 404, got %d", w.Code)
	}
}

func TestTag_FindInvalidID(t *testing.T) {
	w := performRequest("GET", "/v1/tag/invalid-id", nil, "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTag_DeleteNonExistent(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "tag_dne@test.com", "password123")

	w := performRequest("DELETE", "/v1/tag/507f1f77bcf86cd799439011", nil, accessToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTag_UpdateNonExistent(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "tag_upd_ne@test.com", "password123")

	w := performRequest("PATCH", "/v1/tag/507f1f77bcf86cd799439011",
		map[string]string{"name": "test"}, accessToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

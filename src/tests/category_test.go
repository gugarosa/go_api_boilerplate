package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCategory_ListEmpty(t *testing.T) {
	w := performRequest("GET", "/v1/category/", nil, "")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	items, ok := resp["response"].([]interface{})
	if !ok {
		t.Fatalf("expected array response, got %T: %v", resp["response"], resp["response"])
	}
	_ = items // may or may not be empty depending on test ordering
}

func TestCategory_CreateWithoutAuth(t *testing.T) {
	w := performRequest("POST", "/v1/category/",
		map[string]string{"name": "Unauthorized"}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCategory_CRUD(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "cat_crud@test.com", "password123")

	// CREATE
	w := performRequest("POST", "/v1/category/",
		map[string]string{"name": "Test Category"}, accessToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// LIST
	w = performRequest("GET", "/v1/category/", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	items := resp["response"].([]interface{})
	if len(items) == 0 {
		t.Fatal("list: expected at least one category")
	}

	// Get the last created category's ID
	lastItem := items[len(items)-1].(map[string]interface{})
	catID := lastItem["_id"].(string)

	// FIND
	w = performRequest("GET", fmt.Sprintf("/v1/category/%s", catID), nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("find: expected 200, got %d", w.Code)
	}

	findResp := parseJSON(w)
	found := findResp["response"].(map[string]interface{})
	if found["name"] != "Test Category" {
		t.Errorf("find: expected 'Test Category', got '%v'", found["name"])
	}

	// UPDATE
	w = performRequest("PATCH", fmt.Sprintf("/v1/category/%s", catID),
		map[string]string{"name": "Updated Category"}, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify update
	w = performRequest("GET", fmt.Sprintf("/v1/category/%s", catID), nil, "")
	findResp = parseJSON(w)
	found = findResp["response"].(map[string]interface{})
	if found["name"] != "Updated Category" {
		t.Errorf("expected 'Updated Category', got '%v'", found["name"])
	}

	// DELETE
	w = performRequest("DELETE", fmt.Sprintf("/v1/category/%s", catID), nil, accessToken)
	if w.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", w.Code)
	}

	// Verify delete
	w = performRequest("GET", fmt.Sprintf("/v1/category/%s", catID), nil, "")
	if w.Code != http.StatusNotFound {
		t.Errorf("find after delete: expected 404, got %d", w.Code)
	}
}

func TestCategory_CreateMissingName(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "cat_noname@test.com", "password123")

	w := performRequest("POST", "/v1/category/", map[string]string{}, accessToken)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCategory_FindInvalidID(t *testing.T) {
	w := performRequest("GET", "/v1/category/not-a-valid-id", nil, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCategory_FindNonExistent(t *testing.T) {
	w := performRequest("GET", "/v1/category/507f1f77bcf86cd799439011", nil, "")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCategory_DeleteWithoutAuth(t *testing.T) {
	w := performRequest("DELETE", "/v1/category/507f1f77bcf86cd799439011", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCategory_UpdateWithoutAuth(t *testing.T) {
	w := performRequest("PATCH", "/v1/category/507f1f77bcf86cd799439011",
		map[string]string{"name": "test"}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCategory_DeleteNonExistent(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "cat_del_ne@test.com", "password123")

	w := performRequest("DELETE", "/v1/category/507f1f77bcf86cd799439011", nil, accessToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCategory_UpdateNonExistent(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "cat_upd_ne@test.com", "password123")

	w := performRequest("PATCH", "/v1/category/507f1f77bcf86cd799439011",
		map[string]string{"name": "test"}, accessToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

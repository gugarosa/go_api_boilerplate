package tests

import (
	"net/http"
	"testing"
)

func TestNonExistentRoute(t *testing.T) {
	w := performRequest("GET", "/v1/nonexistent", nil, "")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	resp := parseJSON(w)
	if resp["response"] != "This route is not available." {
		t.Errorf("unexpected message: %v", resp["response"])
	}
}

func TestInvalidJSONBody(t *testing.T) {
	w := performRawRequest("POST", "/v1/register", "{invalid json}", "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestEmptyBody_Register(t *testing.T) {
	w := performRequest("POST", "/v1/register", nil, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty body, got %d", w.Code)
	}
}

func TestProtectedEndpoint_MalformedJWT(t *testing.T) {
	w := performRequest("POST", "/v1/category/",
		map[string]string{"name": "test"},
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.expired.invalid")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestWrongContentType(t *testing.T) {
	w := performRequestWithContentType("POST", "/v1/register",
		"email=test@test.com&password=password123", "", "application/x-www-form-urlencoded")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for wrong content type, got %d", w.Code)
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "del_nonex_edge@test.com", "password123")

	endpoints := []string{
		"/v1/category/507f1f77bcf86cd799439011",
		"/v1/tag/507f1f77bcf86cd799439011",
		"/v1/product/507f1f77bcf86cd799439011",
	}

	for _, ep := range endpoints {
		w := performRequest("DELETE", ep, nil, accessToken)
		if w.Code != http.StatusNotFound {
			t.Errorf("%s: expected 404, got %d", ep, w.Code)
		}
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "upd_nonex_edge@test.com", "password123")

	endpoints := []struct {
		path string
		body map[string]string
	}{
		{"/v1/category/507f1f77bcf86cd799439011", map[string]string{"name": "x"}},
		{"/v1/tag/507f1f77bcf86cd799439011", map[string]string{"name": "x"}},
	}

	for _, ep := range endpoints {
		w := performRequest("PATCH", ep.path, ep.body, accessToken)
		if w.Code != http.StatusNotFound {
			t.Errorf("%s: expected 404, got %d", ep.path, w.Code)
		}
	}
}

func TestInvalidObjectID_AllResources(t *testing.T) {
	paths := []string{
		"/v1/category/not-an-id",
		"/v1/tag/not-an-id",
		"/v1/product/not-an-id",
	}

	for _, path := range paths {
		w := performRequest("GET", path, nil, "")
		if w.Code != http.StatusBadRequest {
			t.Errorf("GET %s: expected 400, got %d", path, w.Code)
		}
	}
}

func TestMethodNotAllowed(t *testing.T) {
	// PUT is not a registered method on these routes
	w := performRequest("PUT", "/v1/category/", nil, "")

	// Gin returns 404 for unregistered method+path combos (no MethodNotAllowed handler)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

package tests

import (
	"net/http"
	"testing"
)

// --- Register Tests ---

func TestRegister_Success(t *testing.T) {
	w := performRequest("POST", "/v1/register",
		map[string]string{"email": "reg_success@test.com", "password": "password123"}, "")

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["response"] != "Document successfully inserted." {
		t.Errorf("unexpected response: %v", resp["response"])
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	performRequest("POST", "/v1/register",
		map[string]string{"email": "dup@test.com", "password": "password123"}, "")

	w := performRequest("POST", "/v1/register",
		map[string]string{"email": "dup@test.com", "password": "password456"}, "")

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	w := performRequest("POST", "/v1/register",
		map[string]string{"email": "not-an-email", "password": "password123"}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	w := performRequest("POST", "/v1/register",
		map[string]string{"email": "short@test.com", "password": "short"}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	w := performRequest("POST", "/v1/register", map[string]string{}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Login Tests ---

func TestLogin_Success(t *testing.T) {
	performRequest("POST", "/v1/register",
		map[string]string{"email": "login_ok@test.com", "password": "password123"}, "")

	w := performRequest("POST", "/v1/login",
		map[string]string{"email": "login_ok@test.com", "password": "password123"}, "")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["access_token"] == nil || resp["access_token"] == "" {
		t.Error("expected access_token in response")
	}
	if resp["refresh_token"] == nil || resp["refresh_token"] == "" {
		t.Error("expected refresh_token in response")
	}
}

func TestLogin_WrongEmail(t *testing.T) {
	w := performRequest("POST", "/v1/login",
		map[string]string{"email": "nonexistent@test.com", "password": "password123"}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	performRequest("POST", "/v1/register",
		map[string]string{"email": "wrongpw@test.com", "password": "password123"}, "")

	w := performRequest("POST", "/v1/login",
		map[string]string{"email": "wrongpw@test.com", "password": "wrongpassword"}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	w := performRequest("POST", "/v1/login", map[string]string{}, "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Refresh Tests ---

func TestRefresh_Success(t *testing.T) {
	_, refreshToken := registerAndLogin(t, "refresh_ok@test.com", "password123")

	w := performRequest("POST", "/v1/refresh",
		map[string]string{"refresh_token": refreshToken}, "")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	resp := parseJSON(w)
	if resp["access_token"] == nil || resp["access_token"] == "" {
		t.Error("expected new access_token")
	}
	if resp["refresh_token"] == nil || resp["refresh_token"] == "" {
		t.Error("expected new refresh_token")
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	w := performRequest("POST", "/v1/refresh",
		map[string]string{"refresh_token": "invalid-refresh-token"}, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRefresh_UsedTokenIsInvalidated(t *testing.T) {
	_, refreshToken := registerAndLogin(t, "refresh_used@test.com", "password123")

	// Use refresh token
	w := performRequest("POST", "/v1/refresh",
		map[string]string{"refresh_token": refreshToken}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("first refresh: expected 200, got %d", w.Code)
	}

	// Try to reuse the same refresh token — the old UUID was deleted
	w = performRequest("POST", "/v1/refresh",
		map[string]string{"refresh_token": refreshToken}, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("reuse refresh: expected 401, got %d", w.Code)
	}
}

// --- Logout Tests ---

func TestLogout_Success(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "logout_ok@test.com", "password123")

	w := performRequest("POST", "/v1/logout", nil, accessToken)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLogout_NoToken(t *testing.T) {
	w := performRequest("POST", "/v1/logout", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogout_InvalidToken(t *testing.T) {
	w := performRequest("POST", "/v1/logout", nil, "invalid-token")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogout_ThenAccessProtected(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "logout_then@test.com", "password123")

	// Logout
	performRequest("POST", "/v1/logout", nil, accessToken)

	// Try accessing a protected endpoint — JWT is valid but Redis session gone
	w := performRequest("POST", "/v1/category/",
		map[string]string{"name": "test"}, accessToken)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 after logout, got %d", w.Code)
	}
}

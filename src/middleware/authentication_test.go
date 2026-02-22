package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("ACCESS_SECRET", "test_access_secret")
	os.Setenv("REFRESH_SECRET", "test_refresh_secret")
}

func TestCreateToken_Success(t *testing.T) {
	id := primitive.NewObjectID()
	token, err := CreateToken(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if token.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if token.AccessUUID == "" {
		t.Error("expected non-empty access UUID")
	}
	if token.RefreshUUID == "" {
		t.Error("expected non-empty refresh UUID")
	}
	if token.AccessExpires == 0 {
		t.Error("expected non-zero access expiry")
	}
	if token.RefreshExpires == 0 {
		t.Error("expected non-zero refresh expiry")
	}
}

func TestCreateToken_UniqueUUIDs(t *testing.T) {
	id := primitive.NewObjectID()
	t1, _ := CreateToken(id)
	t2, _ := CreateToken(id)

	if t1.AccessUUID == t2.AccessUUID {
		t.Error("expected unique access UUIDs across calls")
	}
	if t1.RefreshUUID == t2.RefreshUUID {
		t.Error("expected unique refresh UUIDs across calls")
	}
}

func TestGetToken_ValidBearer(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer my-token-123")

	result := GetToken(req)
	if result != "my-token-123" {
		t.Errorf("expected 'my-token-123', got '%s'", result)
	}
}

func TestGetToken_NoHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	if result := GetToken(req); result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestGetToken_MalformedHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidNoSpace")
	if result := GetToken(req); result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestGetToken_ExtraSpaces(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer token extra")
	// Split produces 3 parts, len != 2, returns ""
	if result := GetToken(req); result != "" {
		t.Errorf("expected empty string for triple-part header, got '%s'", result)
	}
}

func TestVerifyToken_Valid(t *testing.T) {
	id := primitive.NewObjectID()
	tokenModel, _ := CreateToken(id)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenModel.AccessToken)

	token, err := VerifyToken(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !token.Valid {
		t.Error("expected token to be valid")
	}
}

func TestVerifyToken_InvalidString(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")

	_, err := VerifyToken(req)
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	_, err := VerifyToken(req)
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestGetTokenData_Valid(t *testing.T) {
	id := primitive.NewObjectID()
	tokenModel, _ := CreateToken(id)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenModel.AccessToken)

	data, err := GetTokenData(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.AccessUUID != tokenModel.AccessUUID {
		t.Errorf("UUID mismatch: expected '%s', got '%s'", tokenModel.AccessUUID, data.AccessUUID)
	}
	if data.UserID != id.Hex() {
		t.Errorf("user ID mismatch: expected '%s', got '%s'", id.Hex(), data.UserID)
	}
}

func TestGetTokenData_InvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	_, err := GetTokenData(req)
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestAuthGuard_ValidToken(t *testing.T) {
	id := primitive.NewObjectID()
	tokenModel, _ := CreateToken(id)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	handlerCalled := false
	r.GET("/test", AuthGuard(), func(c *gin.Context) {
		handlerCalled = true
		// Verify token_data was stored in context
		val, exists := c.Get("token_data")
		if !exists {
			t.Error("expected token_data in context")
		}
		if val == nil {
			t.Error("expected non-nil token_data")
		}
		c.Status(http.StatusOK)
	})

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenModel.AccessToken)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !handlerCalled {
		t.Error("expected handler to be called")
	}
}

func TestAuthGuard_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/test", AuthGuard(), func(c *gin.Context) {
		t.Error("handler should not be called without token")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthGuard_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/test", AuthGuard(), func(c *gin.Context) {
		t.Error("handler should not be called with invalid token")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-jwt-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

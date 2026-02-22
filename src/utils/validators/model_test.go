package validators

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testModel struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestBindModel_ValidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"test","email":"a@b.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testModel
	if err := BindModel(c, &model); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if model.Name != "test" {
		t.Errorf("expected 'test', got '%s'", model.Name)
	}
	if model.Email != "a@b.com" {
		t.Errorf("expected 'a@b.com', got '%s'", model.Email)
	}
}

func TestBindModel_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{invalid json}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testModel
	if err := BindModel(c, &model); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestBindModel_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(""))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testModel
	if err := BindModel(c, &model); err == nil {
		t.Error("expected error for empty body, got nil")
	}
}

func TestValidateModel_Valid(t *testing.T) {
	model := testModel{Name: "test", Email: "a@b.com"}
	if err := ValidateModel(model); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateModel_MissingRequired(t *testing.T) {
	model := testModel{}
	if err := ValidateModel(model); err == nil {
		t.Error("expected validation error, got nil")
	}
}

func TestValidateModel_InvalidEmail(t *testing.T) {
	model := testModel{Name: "test", Email: "not-an-email"}
	if err := ValidateModel(model); err == nil {
		t.Error("expected validation error for invalid email, got nil")
	}
}

func TestValidateModel_PartiallyValid(t *testing.T) {
	model := testModel{Name: "test", Email: ""}
	if err := ValidateModel(model); err == nil {
		t.Error("expected validation error for missing email, got nil")
	}
}

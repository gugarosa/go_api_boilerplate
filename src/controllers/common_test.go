package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type sampleStruct struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func TestDecodeStruct(t *testing.T) {
	s := sampleStruct{ID: primitive.NewObjectID(), Name: "test"}
	result, err := DecodeStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected 'test', got '%v'", result["name"])
	}
}

func TestEncodeStruct(t *testing.T) {
	id := primitive.NewObjectID()
	m := bson.M{"_id": id, "name": "test"}

	var s sampleStruct
	if err := EncodeStruct(m, &s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.Name != "test" {
		t.Errorf("expected 'test', got '%s'", s.Name)
	}
	if s.ID != id {
		t.Errorf("expected ID '%s', got '%s'", id.Hex(), s.ID.Hex())
	}
}

func TestDecodeEncodeRoundTrip(t *testing.T) {
	original := sampleStruct{ID: primitive.NewObjectID(), Name: "roundtrip"}

	decoded, err := DecodeStruct(original)
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	var result sampleStruct
	if err := EncodeStruct(decoded, &result); err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}

	if result.ID != original.ID || result.Name != original.Name {
		t.Errorf("round-trip mismatch: %+v vs %+v", original, result)
	}
}

func TestDecodeStruct_WithTimestamps(t *testing.T) {
	type ts struct {
		Name      string    `bson:"name"`
		CreatedAt time.Time `bson:"created_at,omitempty"`
	}
	s := ts{Name: "test", CreatedAt: time.Now()}
	result, err := DecodeStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected 'test', got '%v'", result["name"])
	}
}

func TestBindAndValidateRequest_Valid(t *testing.T) {
	type testReq struct {
		Name string `json:"name" validate:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"valid"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testReq
	if err := BindAndValidateRequest(c, &model); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if model.Name != "valid" {
		t.Errorf("expected 'valid', got '%s'", model.Name)
	}
}

func TestBindAndValidateRequest_InvalidJSON(t *testing.T) {
	type testReq struct {
		Name string `json:"name" validate:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{bad}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testReq
	if err := BindAndValidateRequest(c, &model); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestBindAndValidateRequest_MissingRequired(t *testing.T) {
	type testReq struct {
		Name string `json:"name" validate:"required"`
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":""}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var model testReq
	if err := BindAndValidateRequest(c, &model); err == nil {
		t.Error("expected validation error for empty required field, got nil")
	}
}

package utils

import (
	"errors"
	"testing"
)

func TestLogError_NoArgs(t *testing.T) {
	if err := LogError(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestLogError_ReturnsFirstError(t *testing.T) {
	err1 := errors.New("first")
	err2 := errors.New("second")

	result := LogError(err1, err2)
	if result != err1 {
		t.Errorf("expected first error, got %v", result)
	}
}

func TestLogError_SkipsNils(t *testing.T) {
	expected := errors.New("real error")
	result := LogError(nil, nil, expected)
	if result != expected {
		t.Errorf("expected 'real error', got %v", result)
	}
}

func TestLogError_AllNils(t *testing.T) {
	if err := LogError(nil, nil, nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestLogError_SingleNil(t *testing.T) {
	if err := LogError(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestLogError_SingleError(t *testing.T) {
	expected := errors.New("single")
	if err := LogError(expected); err != expected {
		t.Errorf("expected 'single', got %v", err)
	}
}

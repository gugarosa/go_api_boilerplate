package tests

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkListTags benchmarks the public list endpoint (no auth).
func BenchmarkListTags(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest("GET", "/v1/tag/", nil, "")
		if w.Code != http.StatusOK {
			b.Fatalf("expected 200, got %d", w.Code)
		}
	}
}

// BenchmarkListCategories benchmarks the public list endpoint for categories.
func BenchmarkListCategories(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest("GET", "/v1/category/", nil, "")
		if w.Code != http.StatusOK {
			b.Fatalf("expected 200, got %d", w.Code)
		}
	}
}

// BenchmarkLogin benchmarks the login flow.
func BenchmarkLogin(b *testing.B) {
	email := "bench_login@test.com"
	password := "password123"
	performRequest("POST", "/v1/register",
		map[string]string{"email": email, "password": password}, "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest("POST", "/v1/login",
			map[string]string{"email": email, "password": password}, "")
		if w.Code != http.StatusOK {
			b.Fatalf("expected 200, got %d", w.Code)
		}
	}
}

// BenchmarkCreateTag benchmarks authenticated tag creation.
func BenchmarkCreateTag(b *testing.B) {
	email := "bench_tag@test.com"
	password := "password123"
	performRequest("POST", "/v1/register",
		map[string]string{"email": email, "password": password}, "")

	w := performRequest("POST", "/v1/login",
		map[string]string{"email": email, "password": password}, "")
	resp := parseJSON(w)
	token := resp["access_token"].(string)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest("POST", "/v1/tag/",
			map[string]string{"name": fmt.Sprintf("Bench Tag %d", i)}, token)
		if w.Code != http.StatusCreated {
			b.Fatalf("expected 201, got %d", w.Code)
		}
	}
}

// TestConcurrentListRequests fires many goroutines making read requests simultaneously.
func TestConcurrentListRequests(t *testing.T) {
	const numRequests = 100
	var wg sync.WaitGroup
	var failures int64

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			w := performRequest("GET", "/v1/tag/", nil, "")
			if w.Code != http.StatusOK {
				atomic.AddInt64(&failures, 1)
			}
		}()
	}
	wg.Wait()

	if failures > 0 {
		t.Errorf("%d/%d concurrent list requests failed", failures, numRequests)
	}
}

// TestConcurrentAuthenticatedWrites fires many goroutines creating tags concurrently.
func TestConcurrentAuthenticatedWrites(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "concurrent_write@test.com", "password123")

	const numRequests = 50
	var wg sync.WaitGroup
	var failures int64

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			defer wg.Done()
			w := performRequest("POST", "/v1/tag/",
				map[string]string{"name": fmt.Sprintf("Concurrent Tag %d", idx)}, accessToken)
			if w.Code != http.StatusCreated {
				atomic.AddInt64(&failures, 1)
			}
		}(i)
	}
	wg.Wait()

	if failures > 0 {
		t.Errorf("%d/%d concurrent create requests failed", failures, numRequests)
	}
}

// TestConcurrentRegistrations tests many users registering simultaneously.
func TestConcurrentRegistrations(t *testing.T) {
	const numRequests = 20
	var wg sync.WaitGroup
	var successes int64

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			defer wg.Done()
			w := performRequest("POST", "/v1/register",
				map[string]string{
					"email":    fmt.Sprintf("concurrent_%d@test.com", idx),
					"password": "password123",
				}, "")
			if w.Code == http.StatusCreated {
				atomic.AddInt64(&successes, 1)
			}
		}(i)
	}
	wg.Wait()

	if successes != numRequests {
		t.Errorf("expected %d successful registrations, got %d", numRequests, successes)
	}
}

// TestConcurrentMixedReadWrite tests simultaneous reads and writes.
func TestConcurrentMixedReadWrite(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "mixed_rw@test.com", "password123")

	const numOps = 50
	var wg sync.WaitGroup
	var failures int64

	wg.Add(numOps * 2) // reads + writes

	// Concurrent writes
	for i := 0; i < numOps; i++ {
		go func(idx int) {
			defer wg.Done()
			w := performRequest("POST", "/v1/category/",
				map[string]string{"name": fmt.Sprintf("Mixed Cat %d", idx)}, accessToken)
			if w.Code != http.StatusCreated {
				atomic.AddInt64(&failures, 1)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numOps; i++ {
		go func() {
			defer wg.Done()
			w := performRequest("GET", "/v1/category/", nil, "")
			if w.Code != http.StatusOK {
				atomic.AddInt64(&failures, 1)
			}
		}()
	}

	wg.Wait()

	if failures > 0 {
		t.Errorf("%d/%d mixed read/write operations failed", failures, numOps*2)
	}
}

// TestHighVolumeCreation creates many items and verifies they all appear in list.
func TestHighVolumeCreation(t *testing.T) {
	accessToken, _ := registerAndLogin(t, "highvol@test.com", "password123")

	const numItems = 100
	for i := 0; i < numItems; i++ {
		w := performRequest("POST", "/v1/tag/",
			map[string]string{"name": fmt.Sprintf("Volume Tag %d", i)}, accessToken)
		if w.Code != http.StatusCreated {
			t.Fatalf("create #%d: expected 201, got %d", i, w.Code)
		}
	}

	// Verify all are listed
	w := performRequest("GET", "/v1/tag/", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}

	resp := parseJSON(w)
	items := resp["response"].([]interface{})
	if len(items) < numItems {
		t.Errorf("expected at least %d tags, got %d", numItems, len(items))
	}
}

// TestRapidLoginLogout cycles through login/logout rapidly.
func TestRapidLoginLogout(t *testing.T) {
	email := "rapid_cycle@test.com"
	password := "password123"
	performRequest("POST", "/v1/register",
		map[string]string{"email": email, "password": password}, "")

	const cycles = 20
	for i := 0; i < cycles; i++ {
		w := performRequest("POST", "/v1/login",
			map[string]string{"email": email, "password": password}, "")
		if w.Code != http.StatusOK {
			t.Fatalf("cycle %d: login failed with %d", i, w.Code)
		}

		resp := parseJSON(w)
		token := resp["access_token"].(string)

		w = performRequest("POST", "/v1/logout", nil, token)
		if w.Code != http.StatusOK {
			t.Fatalf("cycle %d: logout failed with %d", i, w.Code)
		}
	}
}

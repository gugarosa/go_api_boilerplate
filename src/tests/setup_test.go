package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_api_boilerplate/database"
	"go_api_boilerplate/server"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// testRouter is the shared Gin engine for all integration tests
var testRouter *gin.Engine

// TestMain initializes the database, cache, and router for integration tests.
// Requires MongoDB and Redis to be running (provided by docker-compose.test.yml).
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASS")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		url.QueryEscape(dbUser), url.QueryEscape(dbPass), dbHost, dbPort)
	database.InitMongo(mongoURI, dbName)
	database.InitRedis(redisHost, redisPort, redisPass)

	testRouter = gin.New()
	server.InitRouter(testRouter)

	// Clean state before tests
	cleanupCollections()

	code := m.Run()

	// Clean state after tests
	cleanupCollections()

	os.Exit(code)
}

func cleanupCollections() {
	ctx := context.Background()
	database.UserCollection.Drop(ctx)
	database.CategoryCollection.Drop(ctx)
	database.ProductCollection.Drop(ctx)
	database.QuestionCollection.Drop(ctx)
	database.SurveyCollection.Drop(ctx)
	database.TagCollection.Drop(ctx)
}

// performRequest sends an HTTP request to the test router and returns the recorder.
func performRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

// performRawRequest sends raw string body with JSON content-type.
func performRawRequest(method, path, rawBody, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

// performRequestWithContentType sends a request with a custom content-type.
func performRequestWithContentType(method, path, rawBody, token, contentType string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", contentType)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

// parseJSON unmarshals the response body into a map.
func parseJSON(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

// registerAndLogin creates a user and returns access + refresh tokens.
func registerAndLogin(t *testing.T, email, password string) (accessToken, refreshToken string) {
	t.Helper()

	reg := performRequest("POST", "/v1/register",
		map[string]string{"email": email, "password": password}, "")
	if reg.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", reg.Code, reg.Body.String())
	}

	login := performRequest("POST", "/v1/login",
		map[string]string{"email": email, "password": password}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d %s", login.Code, login.Body.String())
	}

	resp := parseJSON(login)
	return resp["access_token"].(string), resp["refresh_token"].(string)
}

// loginOnly logs in an existing user and returns tokens.
func loginOnly(t *testing.T, email, password string) (accessToken, refreshToken string) {
	t.Helper()

	login := performRequest("POST", "/v1/login",
		map[string]string{"email": email, "password": password}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d %s", login.Code, login.Body.String())
	}

	resp := parseJSON(login)
	return resp["access_token"].(string), resp["refresh_token"].(string)
}

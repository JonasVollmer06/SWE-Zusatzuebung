package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"swe-zusatzuebung/internal/config"
	"swe-zusatzuebung/internal/database"
	"swe-zusatzuebung/internal/fussballer"
	"swe-zusatzuebung/internal/server"

	"gorm.io/gorm"
)

func newReadAPIRouter(t *testing.T) http.Handler {
	t.Helper()

	router, _ := newAPITestRouter(t)
	return router
}

func newAPITestRouter(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Skipf("PostgreSQL is not available for integration tests: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected SQL database handle, got %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	repository := fussballer.NewRepository(db)
	readService := fussballer.NewReadService(repository)
	writeService := fussballer.NewWriteService(repository)

	return server.NewRouter(fussballer.NewRouter(readService, writeService)), db
}

func request(router http.Handler, method string, target string) *httptest.ResponseRecorder {
	return requestWithHeaders(router, method, target, nil)
}

func requestWithHeaders(
	router http.Handler,
	method string,
	target string,
	headers map[string]string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	request.Header.Set("Accept", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func requestJSON(router http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func resetFussballerByUsername(t *testing.T, db *gorm.DB, username string) {
	t.Helper()

	deleteByUsername := func() {
		if err := db.Exec("DELETE FROM fussballer.fussballer WHERE username = ?", username).Error; err != nil {
			t.Fatalf("expected cleanup for username %q to succeed, got %v", username, err)
		}
	}

	deleteByUsername()
	t.Cleanup(deleteByUsername)
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if response.Code != expected {
		t.Fatalf("expected status %d, got %d with body %q", expected, response.Code, response.Body.String())
	}
}

func assertJSONContentType(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected JSON content type, got %q", contentType)
	}
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
}

func assertAll(
	t *testing.T,
	players []fussballer.Fussballer,
	matches func(fussballer.Fussballer) bool,
) {
	t.Helper()

	if len(players) == 0 {
		t.Fatal("expected at least one player")
	}

	for _, player := range players {
		if !matches(player) {
			t.Fatalf("unexpected player in result: %+v", player)
		}
	}
}

func itoa(value int) string {
	return strconv.FormatInt(int64(value), 10)
}

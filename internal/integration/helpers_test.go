package integration

import (
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
)

func newReadAPIRouter(t *testing.T) http.Handler {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.Connect(ctx, config.Load().DatabaseURL)
	if err != nil {
		t.Skipf("PostgreSQL is not available for integration tests: %v", err)
	}
	t.Cleanup(pool.Close)

	repository := fussballer.NewRepository(pool)
	readService := fussballer.NewReadService(repository)
	return server.NewRouter(fussballer.NewRouter(readService))
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

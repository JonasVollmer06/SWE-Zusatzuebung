package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"swe-zusatzuebung/internal/config"
	"swe-zusatzuebung/internal/database"
	"swe-zusatzuebung/internal/fussballer"
	"swe-zusatzuebung/internal/server"
)

func TestReadAPIFindByID(t *testing.T) {
	router := newReadAPIRouter(t)

	for _, id := range []int{20, 40} {
		t.Run("existing id", func(t *testing.T) {
			response := request(router, http.MethodGet, "/fussballer/"+itoa(id))

			assertStatus(t, response, http.StatusOK)
			assertJSONContentType(t, response)

			var player fussballer.Fussballer
			decodeJSON(t, response, &player)

			if player.ID != id {
				t.Fatalf("expected id %d, got %d", id, player.ID)
			}
			if player.Adresse == nil {
				t.Fatal("expected address to be included")
			}
		})
	}
}

func TestReadAPIFindByIDNotFound(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer/9999")

	assertStatus(t, response, http.StatusNotFound)
}

func TestReadAPIFindByIDInvalidID(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer/xyz")

	assertStatus(t, response, http.StatusNotFound)
}

func TestReadAPIFindByIDNotModified(t *testing.T) {
	router := newReadAPIRouter(t)

	response := requestWithHeaders(
		router,
		http.MethodGet,
		"/fussballer/40",
		map[string]string{"If-None-Match": `"0"`},
	)

	assertStatus(t, response, http.StatusNotModified)
	if body := response.Body.String(); body != "" {
		t.Fatalf("expected empty body, got %q", body)
	}
}

func TestReadAPIFindAll(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer")

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)

	var page fussballer.Page
	decodeJSON(t, response, &page)

	if len(page.Content) != fussballer.DefaultPageSize {
		t.Fatalf("expected default page size %d, got %d", fussballer.DefaultPageSize, len(page.Content))
	}
	if page.Page.TotalElements != 7 {
		t.Fatalf("expected totalElements 7, got %d", page.Page.TotalElements)
	}
}

func TestReadAPIFindByNachname(t *testing.T) {
	router := newReadAPIRouter(t)

	for _, nachname := range []string{"Vollmer", "Hery"} {
		t.Run(nachname, func(t *testing.T) {
			target := "/fussballer?nachname=" + url.QueryEscape(nachname)
			response := request(router, http.MethodGet, target)

			assertStatus(t, response, http.StatusOK)
			assertJSONContentType(t, response)

			var page fussballer.Page
			decodeJSON(t, response, &page)

			assertAll(t, page.Content, func(player fussballer.Fussballer) bool {
				return player.Nachname == nachname
			})
		})
	}
}

func TestReadAPIFindByNationalitaet(t *testing.T) {
	router := newReadAPIRouter(t)

	for _, nationalitaet := range []string{"Türkei", "Senegal"} {
		t.Run(nationalitaet, func(t *testing.T) {
			target := "/fussballer?nationalitaet=" + url.QueryEscape(nationalitaet)
			response := request(router, http.MethodGet, target)

			assertStatus(t, response, http.StatusOK)
			assertJSONContentType(t, response)

			var page fussballer.Page
			decodeJSON(t, response, &page)

			assertAll(t, page.Content, func(player fussballer.Fussballer) bool {
				return player.Nationalitaet == nationalitaet
			})
		})
	}
}

func TestReadAPIFindByPosition(t *testing.T) {
	router := newReadAPIRouter(t)

	for _, position := range []fussballer.Position{
		fussballer.PositionMittelfeldspieler,
		fussballer.PositionVerteidiger,
	} {
		t.Run(string(position), func(t *testing.T) {
			target := "/fussballer?position=" + url.QueryEscape(string(position))
			response := request(router, http.MethodGet, target)

			assertStatus(t, response, http.StatusOK)
			assertJSONContentType(t, response)

			var page fussballer.Page
			decodeJSON(t, response, &page)

			assertAll(t, page.Content, func(player fussballer.Fussballer) bool {
				return player.Position != nil && *player.Position == position
			})
		})
	}
}

func TestReadAPIFindNotFound(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer?nachname=Mustermann")

	assertStatus(t, response, http.StatusNotFound)
}

func TestReadAPIRejectsInvalidSearchParameter(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer?foo=bar")

	assertStatus(t, response, http.StatusBadRequest)
}

func TestReadAPICountOnly(t *testing.T) {
	router := newReadAPIRouter(t)

	response := request(router, http.MethodGet, "/fussballer?count-only=true")

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)

	var result map[string]int
	decodeJSON(t, response, &result)

	if result["count"] != 7 {
		t.Fatalf("expected count 7, got %d", result["count"])
	}
}

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

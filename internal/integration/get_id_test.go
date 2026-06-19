package integration

import (
	"net/http"
	"testing"

	"swe-zusatzuebung/internal/fussballer"
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

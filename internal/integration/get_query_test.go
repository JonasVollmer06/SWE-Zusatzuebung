package integration

import (
	"net/http"
	"net/url"
	"testing"

	"swe-zusatzuebung/internal/fussballer"
)

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

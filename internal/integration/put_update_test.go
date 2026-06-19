package integration

import (
	"net/http"
	"testing"

	"swe-zusatzuebung/internal/fussballer"
)

func TestWriteAPIUpdate(t *testing.T) {
	router, _ := newAPITestRouter(t)
	t.Cleanup(func() {
		response := request(router, http.MethodPost, "/fussballer/reset")
		assertStatus(t, response, http.StatusOK)
	})

	body := `{
		"nachname": "Aktualisiert",
		"nationalitaet": "Deutschland",
		"position": "STUERMER",
		"geburtsdatum": "2002-02-02T00:00:00Z",
		"username": "updated-20",
		"adresse": {
			"plz": "76131",
			"ort": "Karlsruhe",
			"bundesland": "Baden-Wuerttemberg"
		}
	}`
	response := requestJSON(router, http.MethodPut, "/fussballer/20", body)

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)
	if response.Header().Get("ETag") != `"1"` {
		t.Fatalf("expected ETag %q, got %q", `"1"`, response.Header().Get("ETag"))
	}

	var updated fussballer.Fussballer
	decodeJSON(t, response, &updated)
	if updated.ID != 20 || updated.Nachname != "Aktualisiert" || updated.Username != "updated-20" {
		t.Fatalf("unexpected updated player: %+v", updated)
	}
	if updated.Version != 1 {
		t.Fatalf("expected version 1, got %d", updated.Version)
	}
	if updated.Position == nil || *updated.Position != fussballer.PositionStuermer {
		t.Fatalf("expected position %q, got %+v", fussballer.PositionStuermer, updated.Position)
	}
	if updated.Adresse == nil || updated.Adresse.Ort != "Karlsruhe" {
		t.Fatalf("expected updated address, got %+v", updated.Adresse)
	}

	getResponse := request(router, http.MethodGet, "/fussballer/20")
	assertStatus(t, getResponse, http.StatusOK)

	var persisted fussballer.Fussballer
	decodeJSON(t, getResponse, &persisted)
	if persisted.Nachname != "Aktualisiert" || persisted.Username != "updated-20" {
		t.Fatalf("expected persisted update, got %+v", persisted)
	}
}

func TestWriteAPIUpdateNotFound(t *testing.T) {
	router, _ := newAPITestRouter(t)

	body := `{
		"nachname": "Aktualisiert",
		"nationalitaet": "Deutschland",
		"position": "STUERMER",
		"geburtsdatum": "2002-02-02T00:00:00Z",
		"username": "updated-not-found"
	}`
	response := requestJSON(router, http.MethodPut, "/fussballer/9999", body)

	assertStatus(t, response, http.StatusNotFound)
}

func TestWriteAPIUpdateRejectsInvalidPosition(t *testing.T) {
	router, _ := newAPITestRouter(t)

	body := `{
		"nachname": "Aktualisiert",
		"nationalitaet": "Deutschland",
		"position": "TRAINER",
		"geburtsdatum": "2002-02-02T00:00:00Z",
		"username": "updated-invalid"
	}`
	response := requestJSON(router, http.MethodPut, "/fussballer/20", body)

	assertStatus(t, response, http.StatusBadRequest)
}

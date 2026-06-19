package integration

import (
	"net/http"
	"testing"

	"swe-zusatzuebung/internal/fussballer"
)

func TestWriteAPIDelete(t *testing.T) {
	router, db := newAPITestRouter(t)
	username := "integration-delete"
	resetFussballerByUsername(t, db, username)

	createResponse := requestJSON(router, http.MethodPost, "/fussballer", `{
		"nachname": "Delete",
		"nationalitaet": "Deutschland",
		"position": "VERTEIDIGER",
		"geburtsdatum": "2001-01-01T00:00:00Z",
		"username": "integration-delete"
	}`)
	assertStatus(t, createResponse, http.StatusCreated)

	var created fussballer.Fussballer
	decodeJSON(t, createResponse, &created)

	deleteResponse := request(router, http.MethodDelete, "/fussballer/"+itoa(created.ID))
	assertStatus(t, deleteResponse, http.StatusNoContent)

	getResponse := request(router, http.MethodGet, "/fussballer/"+itoa(created.ID))
	assertStatus(t, getResponse, http.StatusNotFound)
}

func TestWriteAPIDeleteNotFound(t *testing.T) {
	router, _ := newAPITestRouter(t)

	response := request(router, http.MethodDelete, "/fussballer/9999")

	assertStatus(t, response, http.StatusNotFound)
}

func TestWriteAPIReset(t *testing.T) {
	router, _ := newAPITestRouter(t)

	createResponse := requestJSON(router, http.MethodPost, "/fussballer", `{
		"nachname": "Reset",
		"nationalitaet": "Deutschland",
		"position": "STUERMER",
		"geburtsdatum": "2001-01-01T00:00:00Z",
		"username": "integration-reset"
	}`)
	assertStatus(t, createResponse, http.StatusCreated)

	resetResponse := request(router, http.MethodPost, "/fussballer/reset")
	assertStatus(t, resetResponse, http.StatusOK)
	assertJSONContentType(t, resetResponse)

	countResponse := request(router, http.MethodGet, "/fussballer?count-only=true")
	assertStatus(t, countResponse, http.StatusOK)

	var result map[string]int
	decodeJSON(t, countResponse, &result)
	if result["count"] != 7 {
		t.Fatalf("expected CSV count 7 after reset, got %d", result["count"])
	}

	searchResponse := request(router, http.MethodGet, "/fussballer?nachname=Reset")
	assertStatus(t, searchResponse, http.StatusNotFound)
}

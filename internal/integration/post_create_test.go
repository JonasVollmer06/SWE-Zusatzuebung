package integration

import (
	"net/http"
	"strings"
	"testing"

	"swe-zusatzuebung/internal/fussballer"
)

func TestWriteAPICreate(t *testing.T) {
	router, db := newAPITestRouter(t)
	username := "integration-create"
	resetFussballerByUsername(t, db, username)

	body := `{
		"nachname": "Musiala",
		"nationalitaet": "Deutschland",
		"position": "MITTELFELDSPIELER",
		"geburtsdatum": "2003-02-26T00:00:00Z",
		"username": "integration-create",
		"adresse": {
			"plz": "80331",
			"ort": "Muenchen",
			"bundesland": "Bayern"
		}
	}`
	response := requestJSON(router, http.MethodPost, "/fussballer", body)

	assertStatus(t, response, http.StatusCreated)
	assertJSONContentType(t, response)

	location := response.Header().Get("Location")
	if !strings.HasPrefix(location, "/fussballer/") {
		t.Fatalf("expected Location header for created player, got %q", location)
	}
	if response.Header().Get("ETag") != `"0"` {
		t.Fatalf("expected ETag %q, got %q", `"0"`, response.Header().Get("ETag"))
	}

	var created fussballer.Fussballer
	decodeJSON(t, response, &created)

	if created.ID == 0 {
		t.Fatal("expected generated id")
	}
	if created.Nachname != "Musiala" || created.Username != username {
		t.Fatalf("unexpected created player: %+v", created)
	}
	if created.Position == nil || *created.Position != fussballer.PositionMittelfeldspieler {
		t.Fatalf("expected position %q, got %+v", fussballer.PositionMittelfeldspieler, created.Position)
	}
	if created.Adresse == nil || created.Adresse.Ort != "Muenchen" {
		t.Fatalf("expected created address, got %+v", created.Adresse)
	}

	getResponse := request(router, http.MethodGet, location)
	assertStatus(t, getResponse, http.StatusOK)

	var persisted fussballer.Fussballer
	decodeJSON(t, getResponse, &persisted)
	if persisted.ID != created.ID || persisted.Username != username {
		t.Fatalf("expected persisted player %+v, got %+v", created, persisted)
	}
}

func TestWriteAPIRejectsInvalidJSON(t *testing.T) {
	router, _ := newAPITestRouter(t)

	response := requestJSON(router, http.MethodPost, "/fussballer", `{"nachname":`)

	assertStatus(t, response, http.StatusBadRequest)
}

func TestWriteAPIRejectsMissingRequiredField(t *testing.T) {
	router, _ := newAPITestRouter(t)

	body := `{
		"nationalitaet": "Deutschland",
		"position": "MITTELFELDSPIELER",
		"geburtsdatum": "2003-02-26T00:00:00Z",
		"username": "integration-missing-name"
	}`
	response := requestJSON(router, http.MethodPost, "/fussballer", body)

	assertStatus(t, response, http.StatusBadRequest)
}

func TestWriteAPIRejectsInvalidPosition(t *testing.T) {
	router, _ := newAPITestRouter(t)

	body := `{
		"nachname": "Musiala",
		"nationalitaet": "Deutschland",
		"position": "TRAINER",
		"geburtsdatum": "2003-02-26T00:00:00Z",
		"username": "integration-invalid-position"
	}`
	response := requestJSON(router, http.MethodPost, "/fussballer", body)

	assertStatus(t, response, http.StatusBadRequest)
}

func TestWriteAPIRejectsUnsupportedMediaType(t *testing.T) {
	router, _ := newAPITestRouter(t)

	request := requestWithHeaders(
		router,
		http.MethodPost,
		"/fussballer",
		map[string]string{"Content-Type": "text/plain"},
	)

	assertStatus(t, request, http.StatusUnsupportedMediaType)
}

func TestWriteAPIRejectsUnknownField(t *testing.T) {
	router, _ := newAPITestRouter(t)

	body := `{
		"nachname": "Musiala",
		"nationalitaet": "Deutschland",
		"position": "MITTELFELDSPIELER",
		"geburtsdatum": "2003-02-26T00:00:00Z",
		"username": "integration-unknown-field",
		"verein": "FC Bayern"
	}`
	response := requestJSON(router, http.MethodPost, "/fussballer", body)

	assertStatus(t, response, http.StatusBadRequest)
}

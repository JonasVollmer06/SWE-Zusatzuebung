package fussballer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type fakeReader struct {
	findByIDFunc func(ctx context.Context, id int) (*Fussballer, error)
	findFunc     func(ctx context.Context, criteria SearchCriteria, pageable Pageable) (*Slice, error)
	countFunc    func(ctx context.Context, criteria SearchCriteria) (int, error)
}

func (r fakeReader) FindByID(ctx context.Context, id int) (*Fussballer, error) {
	if r.findByIDFunc == nil {
		return nil, errors.New("unexpected FindByID call")
	}

	return r.findByIDFunc(ctx, id)
}

func (r fakeReader) Find(ctx context.Context, criteria SearchCriteria, pageable Pageable) (*Slice, error) {
	if r.findFunc == nil {
		return nil, errors.New("unexpected Find call")
	}

	return r.findFunc(ctx, criteria, pageable)
}

func (r fakeReader) Count(ctx context.Context, criteria SearchCriteria) (int, error) {
	if r.countFunc == nil {
		return 0, errors.New("unexpected Count call")
	}

	return r.countFunc(ctx, criteria)
}

func TestRouterFindByID(t *testing.T) {
	router := NewRouter(fakeReader{
		findByIDFunc: func(_ context.Context, id int) (*Fussballer, error) {
			if id != 1000 {
				t.Fatalf("expected id 1000, got %d", id)
			}

			return &Fussballer{ID: 1000, Version: 3, Nachname: "Neuer"}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/1000", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if response.Header().Get("ETag") != `"3"` {
		t.Fatalf("expected ETag %q, got %q", `"3"`, response.Header().Get("ETag"))
	}

	var player Fussballer
	if err := json.NewDecoder(response.Body).Decode(&player); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if player.ID != 1000 || player.Nachname != "Neuer" {
		t.Fatalf("unexpected player response: %+v", player)
	}
}

func TestRouterFindByIDReturnsNotModified(t *testing.T) {
	router := NewRouter(fakeReader{
		findByIDFunc: func(_ context.Context, _ int) (*Fussballer, error) {
			return &Fussballer{ID: 1000, Version: 3}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/1000", nil)
	request.Header.Set("If-None-Match", `"3"`)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotModified {
		t.Fatalf("expected status %d, got %d", http.StatusNotModified, response.Code)
	}
}

func TestRouterFind(t *testing.T) {
	position := PositionTorwart
	expectedCriteria := SearchCriteria{
		Nachname: "Neuer",
		Position: &position,
	}
	expectedPageable := Pageable{Number: 1, Size: 10}

	router := NewRouter(fakeReader{
		findFunc: func(_ context.Context, criteria SearchCriteria, pageable Pageable) (*Slice, error) {
			if !reflect.DeepEqual(criteria, expectedCriteria) {
				t.Fatalf("expected criteria %+v, got %+v", expectedCriteria, criteria)
			}
			if !reflect.DeepEqual(pageable, expectedPageable) {
				t.Fatalf("expected pageable %+v, got %+v", expectedPageable, pageable)
			}

			return &Slice{
				Content:       []Fussballer{{ID: 1000, Nachname: "Neuer"}},
				TotalElements: 42,
			}, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/?nachname=Neuer&position=TORWART&page=2&size=10", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var page Page
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if len(page.Content) != 1 || page.Content[0].ID != 1000 {
		t.Fatalf("unexpected page content: %+v", page.Content)
	}
	if page.Page.Number != 1 || page.Page.Size != 10 || page.Page.TotalElements != 42 || page.Page.TotalPages != 5 {
		t.Fatalf("unexpected page metadata: %+v", page.Page)
	}
}

func TestRouterCountOnly(t *testing.T) {
	router := NewRouter(fakeReader{
		countFunc: func(_ context.Context, criteria SearchCriteria) (int, error) {
			if criteria.Nationalitaet != "Deutschland" {
				t.Fatalf("expected nationalitaet Deutschland, got %q", criteria.Nationalitaet)
			}

			return 7, nil
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/?count-only=true&nationalitaet=Deutschland", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var result map[string]int
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if result["count"] != 7 {
		t.Fatalf("expected count 7, got %d", result["count"])
	}
}

func TestRouterRejectsInvalidAcceptHeader(t *testing.T) {
	router := NewRouter(fakeReader{})

	request := httptest.NewRequest(http.MethodGet, "/1000", nil)
	request.Header.Set("Accept", "application/xml")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotAcceptable {
		t.Fatalf("expected status %d, got %d", http.StatusNotAcceptable, response.Code)
	}
}

func TestRouterRejectsInvalidPosition(t *testing.T) {
	router := NewRouter(fakeReader{})

	request := httptest.NewRequest(http.MethodGet, "/?position=TRAINER", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

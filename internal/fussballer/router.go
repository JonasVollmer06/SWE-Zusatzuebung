package fussballer

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Reader interface {
	FindByID(ctx context.Context, id int) (*Fussballer, error)
	Find(ctx context.Context, criteria SearchCriteria, pageable Pageable) (*Slice, error)
	Count(ctx context.Context, criteria SearchCriteria) (int, error)
}

type Page struct {
	Content []Fussballer `json:"content"`
	Page    PageMetadata `json:"page"`
}

type PageMetadata struct {
	Size          int `json:"size"`
	Number        int `json:"number"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

func NewRouter(service Reader) http.Handler {
	router := chi.NewRouter()
	handler := routerHandler{service: service}

	router.Get("/", handler.find)
	router.Get("/{id}", handler.findByID)

	return router
}

type routerHandler struct {
	service Reader
}

func (h routerHandler) findByID(w http.ResponseWriter, r *http.Request) {
	if !acceptsJSON(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	player, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	etag := `"` + strconv.Itoa(player.Version) + `"`
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)
	writeJSON(w, http.StatusOK, player)
}

func (h routerHandler) find(w http.ResponseWriter, r *http.Request) {
	if !acceptsJSON(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	criteria, err := parseSearchCriteria(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if _, ok := r.URL.Query()["count-only"]; ok {
		count, err := h.service.Count(r.Context(), criteria)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]int{"count": count})
		return
	}

	pageable := parsePageable(r)
	slice, err := h.service.Find(r.Context(), criteria, pageable)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, createPage(slice, pageable))
}

func parseSearchCriteria(r *http.Request) (SearchCriteria, error) {
	query := r.URL.Query()
	criteria := SearchCriteria{
		Nachname:      query.Get("nachname"),
		Nationalitaet: query.Get("nationalitaet"),
	}

	for key := range query {
		switch key {
		case "nachname", "nationalitaet", "position", "page", "size", "count-only":
		default:
			return SearchCriteria{}, ErrInvalidSearchParameter
		}
	}

	if positionValue := query.Get("position"); positionValue != "" {
		position := Position(positionValue)
		if !position.IsValid() {
			return SearchCriteria{}, ErrInvalidSearchParameter
		}
		criteria.Position = &position
	}

	return criteria, nil
}

func parsePageable(r *http.Request) Pageable {
	number := DefaultPageNumber
	if value := r.URL.Query().Get("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			number = parsed - 1
		}
	}

	size := DefaultPageSize
	if value := r.URL.Query().Get("size"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			size = parsed
		}
	}

	return NewPageable(number, size)
}

func createPage(slice *Slice, pageable Pageable) Page {
	return Page{
		Content: slice.Content,
		Page: PageMetadata{
			Size:          pageable.Size,
			Number:        pageable.Number,
			TotalElements: slice.TotalElements,
			TotalPages:    int(math.Ceil(float64(slice.TotalElements) / float64(pageable.Size))),
		},
	}
}

func acceptsJSON(r *http.Request) bool {
	accept := strings.ToLower(r.Header.Get("Accept"))
	return accept == "" || accept == "*/*" || strings.Contains(accept, "json") || strings.Contains(accept, "html")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrInvalidID):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, ErrInvalidSearchParameter):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

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

type Writer interface {
	Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error)
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

func NewRouter(reader Reader, writers ...Writer) http.Handler {
	router := chi.NewRouter()
	handler := routerHandler{reader: reader}

	if len(writers) > 0 {
		handler.writer = writers[0]
	}

	router.Get("/", handler.find)
	router.Get("/{id}", handler.findByID)
	if handler.writer != nil {
		router.Post("/", handler.create)
	}

	return router
}

type routerHandler struct {
	reader Reader
	writer Writer
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

	player, err := h.reader.FindByID(r.Context(), id)
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
		count, err := h.reader.Count(r.Context(), criteria)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]int{"count": count})
		return
	}

	pageable := parsePageable(r)
	slice, err := h.reader.Find(r.Context(), criteria, pageable)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, createPage(slice, pageable))
}

func (h routerHandler) create(w http.ResponseWriter, r *http.Request) {
	if !acceptsJSON(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if !isJSONContent(r) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var request CreateFussballerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player, err := h.writer.Create(r.Context(), request)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Location", "/fussballer/"+strconv.Itoa(player.ID))
	w.Header().Set("ETag", `"`+strconv.Itoa(player.Version)+`"`)
	writeJSON(w, http.StatusCreated, player)
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

func isJSONContent(r *http.Request) bool {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	return strings.Contains(contentType, "application/json")
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
	case errors.Is(err, ErrInvalidSearchParameter), errors.Is(err, ErrValidation):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

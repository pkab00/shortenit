package link

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/pkab00/shortenit/internal/apperr"
	"github.com/pkab00/shortenit/internal/redirect"
	"github.com/pkab00/shortenit/internal/statistics"
)

type Handler struct {
	linkService       *Service
	redirectService   *redirect.Service
	statisticsService *statistics.Service
}

type CreateLinkRequest struct {
	URL string `json:"url"`
}

func NewHandler(
	linkService *Service,
	redirectService *redirect.Service,
	statisticsService *statistics.Service,
) *Handler {
	return &Handler{
		linkService:       linkService,
		redirectService:   redirectService,
		statisticsService: statisticsService,
	}
}

func (h *Handler) handleServerError(w http.ResponseWriter, err error) {
	log.Println("error getting record by code: ", err)
	switch {
	case errors.Is(err, apperr.ErrorLinkNotFound):
		http.Error(w, "Link Not Found", http.StatusNotFound)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create godoc
// @Summary Creates a new short code based on a provided link
// @Accept json
// @Produce json
// @Param request body CreateLinkRequest true "URL"
// @Success 200 {object} LinkResponse
// @Success 201 {object} LinkResponse
// @Failure 400
// @Failure 500
// @Router /shorten [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error decoding create request: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.Create(ctx, req.URL)
	if err != nil {
		log.Println("error creating new link: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if res.Created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(res.Link)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.PathValue("code")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.Delete(ctx, code)
	if err != nil {
		log.Println("error deleting link: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.All(ctx)
	if err != nil {
		log.Println("error getting links: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.PathValue("code")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.ByCode(ctx, code)
	if err != nil {
		h.handleServerError(w, err)
		return
	}

	_, err = h.redirectService.Increment(ctx, code)
	if err != nil {
		log.Println("error incrementing redirect counter: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	link := res.URL
	http.Redirect(w, r, link, http.StatusFound)
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.PathValue("code")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.statisticsService.Get(ctx, code)
	if err != nil {
		h.handleServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

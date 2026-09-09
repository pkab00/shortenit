package link

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	s *Service
}

type LinkRequest struct {
	URL string `json:"url"`
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req LinkRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error decoding create request: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.s.Create(ctx, req.URL)
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

	res, err := h.s.Delete(ctx, code)
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

	res, err := h.s.All(ctx)
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

	res, err := h.s.ByCode(ctx, code)
	if err != nil {
		log.Println("error getting record by code: ", err)
		switch {
		case errors.Is(err, ErrorLinkNotFound):
			http.Error(w, "Link Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	link := res.Body
	log.Println("redirecting to", link)
	http.Redirect(w, r, link, http.StatusFound)
}

package link

import (
	"context"
	"encoding/json"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.s.Create(ctx, req.URL)
	if err != nil {
		log.Println("error creating new link: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var req LinkRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error decoding delete request: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.s.Delete(ctx, req.URL)
	if err != nil {
		log.Println("error deleting link: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.s.All(ctx)
	if err != nil {
		log.Println("error getting links: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

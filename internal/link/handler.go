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

type LinkResponse struct {
	Code      string `json:"code"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type LinkErrorResponse struct {
	Code         string `json:"code"`
	URL          string `json:"url"`
	ErrorMessage string `json:"error_message"`
}

type CreateManyResponse struct {
	Success []LinkResponse      `json:"success"`
	Failure []LinkErrorResponse `json:"failure"`
}

type CreateLinkRequest struct {
	URL  string  `json:"url"`
	Code *string `json:"code"`
}

type Handler struct {
	linkService       *Service
	redirectService   *redirect.Service
	statisticsService *statistics.Service
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
	switch {
	case errors.Is(err, apperr.ErrorLinkNotFound):
		http.Error(w, "Link Not Found", http.StatusNotFound)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Create godoc
//
//	@Summary	Creates a new short code based on a provided link. If such a code already exists, returns an existing record.
//	@Accept		json
//	@Produce	json
//	@Param		request	body		CreateLinkRequest	true	"Create request"
//	@Success	200		{object}	LinkResponse
//	@Success	201		{object}	LinkResponse
//	@Failure	400
//	@Failure	500
//	@Router		/shorten [post]
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

	res, err := h.linkService.Create(ctx, &req)
	if err != nil {
		log.Println("error creating new link: ", err)
		switch {
		case errors.Is(err, apperr.ErrorInvalidCreateRequest):
			http.Error(w, "Bad Create Request", http.StatusBadRequest)
		case errors.Is(err, apperr.ErrorInvalidCustomCode):
			http.Error(w, "Invalid Custom Code", http.StatusBadRequest)
		case errors.Is(err, apperr.ErrorCustomCodeInUse):
			http.Error(w, "Custom Code In Use", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if res.Created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(res.Response)
}

// Create godoc
//
//	@Summary	Creates multiple records. Returns all the requests grouped by the result as "succeess" and "failure".
//	@Accept		json
//	@Produce	json
//	@Param		request	body		[]CreateLinkRequest	true	"Create request"
//	@Success	200		{object}	CreateManyResponse
//	@Failure	500		{object}	CreateManyResponse
//	@Router		/shorten [post]
func (h *Handler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var reqs []CreateLinkRequest
	var success []LinkResponse
	var failure []LinkErrorResponse
	var response CreateManyResponse
	var err error

	err = json.NewDecoder(r.Body).Decode(&reqs)
	if err != nil {
		log.Println("error decoding create request: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	results := h.linkService.CreateMany(ctx, reqs)
	for index, result := range results {
		if result.Error != nil {
			code := reqs[index].Code
			url := reqs[index].URL
			err = result.Error
			failure = append(
				failure,
				LinkErrorResponse{
					Code:         *code,
					URL:          url,
					ErrorMessage: err.Error(),
				},
			)
		} else {
			success = append(success, *result.Response)
		}
	}

	response = CreateManyResponse{
		Success: success,
		Failure: failure,
	}
	w.Header().Set("Content-Type", "application/json")
	if len(success) == 0 && len(failure) > 0 {
		log.Println("all create operations failed")
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: handler should not always return code 500
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// Create godoc
//
//	@Summary	Deletes a URL by its short code.
//	@Produce	json
//	@Success	200	{object}	LinkResponse
//	@Failure	400
//	@Failure	500
//	@Router		/shorten/{code} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.PathValue("code")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.Delete(ctx, code)
	if err != nil {
		log.Println("error deleting link: ", err)
		h.handleServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Create godoc
//
//	@Summary	Returns the list of all accessable URLs.
//	@Produce	json
//	@Success	200	{array}	LinkResponse
//	@Failure	500
//	@Router		/shorten [get]
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

// Create godoc
//
//	@Summary	Processes redirect to a URL by its short code.
//	@Success	320
//	@Failure	400
//	@Failure	500
//	@Router		/shorten/{code} [get]
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	var err error

	code := r.PathValue("code")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := h.linkService.ByCode(ctx, code)
	if err != nil {
		log.Println("error getting record by code: ", err)
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

// Create godoc
//
//	@Summary	Returns statistics related to a URL by its short code.
//	@Produce	json
//	@Success	200	{object}	statistics.StatisticsResponse
//	@Failure	400
//	@Failure	500
//	@Router		/shorten/{code}/statistics [get]
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

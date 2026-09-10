package link

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/shortenit/internal/redirect"
	"github.com/pkab00/shortenit/internal/statistics"
)

func RegisterRoutes(db *sql.DB, mux *http.ServeMux) {
	var (
		linkRepo       = NewRepository(db)
		redirectRepo   = redirect.NewRepository(db)
		statisticsRepo = statistics.NewRepository(db)

		linkService       = NewService(linkRepo)
		redirectService   = redirect.NewService(redirectRepo)
		statisticsService = statistics.NewService(statisticsRepo)

		hand = NewHandler(linkService, redirectService, statisticsService)
	)

	mux.HandleFunc("POST /shorten", hand.Create)
	mux.HandleFunc("DELETE /shorten/{code}", hand.Delete)
	mux.HandleFunc("GET /shorten/{code}", hand.Redirect)
	mux.HandleFunc("GET /shorten/{code}/statistics", hand.GetStatistics)
	mux.HandleFunc("GET /shorten", hand.All)
}

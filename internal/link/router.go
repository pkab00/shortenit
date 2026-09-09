package link

import (
	"database/sql"
	"net/http"

	"github.com/pkab00/shortenit/internal/redirect"
)

func RegisterRoutes(db *sql.DB, mux *http.ServeMux) {
	linkRepo, redirectRepo := NewRepository(db), redirect.NewRepository(db)
	linkServ, redirectServ := NewService(linkRepo), redirect.NewService(redirectRepo)
	hand := NewHandler(linkServ, redirectServ)

	mux.HandleFunc("POST /shorten", hand.Create)
	mux.HandleFunc("DELETE /shorten/{code}", hand.Delete)
	mux.HandleFunc("GET /shorten/{code}", hand.Redirect)
	mux.HandleFunc("GET /shorten", hand.All)
}

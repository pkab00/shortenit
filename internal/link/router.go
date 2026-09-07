package link

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB, mux *http.ServeMux) {
	repo := NewRepository(db)
	serv := NewService(repo)
	hand := NewHandler(serv)

	mux.HandleFunc("POST /shorten", hand.Create)
	mux.HandleFunc("DELETE /shorten", hand.Delete)
	mux.HandleFunc("GET /links", hand.All)
}

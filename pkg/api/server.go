package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Handler interface {
	RegisterRoutes(*mux.Router)
}

type DBClient interface {
	Query()
}

type Server struct {
	addr    string
	handler Handler
	db      *pgx.Conn
}

func NewServer(addr string, handler Handler, db *pgx.Conn) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
		db:      db,
	}
}

func (a *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	a.handler.RegisterRoutes(subrouter)

	log.Printf("Server listening on %s...", a.addr)

	return http.ListenAndServe(a.addr, router)
}

package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	RegisterRoutes(*mux.Router)
}

// type DBClient interface {
// 	Query()
// 	Exec()
// }

type Server struct {
	addr    string
	handler Handler
}

func NewServer(addr string, handler Handler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (a *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	a.handler.RegisterRoutes(subrouter)

	log.Printf("Server listening on %s...", a.addr)

	return http.ListenAndServe(a.addr, router)
}

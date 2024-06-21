package payments

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct{}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/deposit", h.HandleDeposit).Methods("POST")
	router.HandleFunc("/withdraw", h.HandleWithdraw).Methods("POST")
	router.HandleFunc("/pay", h.HandlePay).Methods("POST")
}

func (h *handler) HandleDeposit(rw http.ResponseWriter, req *http.Request) {

}

func (h *handler) HandleWithdraw(rw http.ResponseWriter, req *http.Request) {

}

func (h *handler) HandlePay(rw http.ResponseWriter, req *http.Request) {

}

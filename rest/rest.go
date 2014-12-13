package rest

import (
	"github.com/gorilla/mux"
	"github.com/jawr/dns/rest/domain"
	"net/http"
)

func Setup() http.Handler {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	domain.Setup(sr)
	return r
}

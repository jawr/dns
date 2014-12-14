package rest

import (
	"github.com/gorilla/mux"
	"github.com/jawr/dns/rest/domain"
	"github.com/jawr/dns/rest/record"
	"github.com/jawr/dns/rest/tld"
	"net/http"
)

func Setup() http.Handler {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	domain.Setup(sr)
	tld.Setup(sr)
	record.Setup(sr)
	return r
}

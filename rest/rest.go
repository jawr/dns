package rest

import (
	"github.com/gorilla/mux"
	"github.com/jawr/dns/rest/domains"
	"github.com/jawr/dns/rest/records"
	"github.com/jawr/dns/rest/whois"
	"net/http"
)

func Setup() http.Handler {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	domains.Setup(sr)
	whois.Setup(sr)
	records.Setup(sr)
	return r
}

package rest

import (
	"github.com/gorilla/mux"
	"github.com/jawr/dns/rest/domain"
	"github.com/jawr/dns/rest/record"
	"github.com/jawr/dns/rest/record_type"
	"github.com/jawr/dns/rest/tld"
	"github.com/jawr/dns/rest/whois"
	"net/http"
)

func Setup() http.Handler {
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	domain.Setup(sr)
	tld.Setup(sr)
	record.Setup(sr)
	record_type.Setup(sr)
	whois.Setup(sr)
	return r
}

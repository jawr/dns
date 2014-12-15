package domain

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/record"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/rest/whois"
	"net/http"
)

type Domain struct {
}

func Setup(r *mux.Router) {
	d := &Domain{}
	sr := r.PathPrefix("/domain").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(d.Search))
	sr.HandleFunc("/{uuid}", d.GetUUID)

	w := whois.Result{}
	wsr := sr.PathPrefix("/{duuid}/whois").Subrouter()
	wsr.HandleFunc("/", injectDomainUUID(paginator.Paginate(w.Search)))

	rec := record.Record{}
	rsr := sr.PathPrefix("/{duuid}/record").Subrouter()
	rsr.HandleFunc("/", injectDomainUUID(paginator.Paginate(rec.Search)))
}

func (d Domain) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (d Domain) GetUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	i, err := db.Get(db.GetByUUID(), uuid)
	util.ToJSON(i, err, w)
}

func injectDomainUUID(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		vars := mux.Vars(r)
		uuid := vars["duuid"]
		context.Set(r, "domain", uuid)
		fn(w, r)
	}
}

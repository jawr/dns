package domain

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/record"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/rest/whois"
	"net/http"
)

type Domain struct {
}

const UUID_REGEX string = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func Setup(r *mux.Router) {
	d := &Domain{}
	sr := r.PathPrefix("/domain").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(d.Search))
	sr.HandleFunc("/{uuid:"+UUID_REGEX+"}", d.GetUUID)
	sr.HandleFunc("/{name}", d.GetName)

	w := whois.Result{}
	sr.HandleFunc("/{duuid:"+UUID_REGEX+"}/whois",
		injectDomainUUID(paginator.Paginate(w.Search)),
	)
	sr.HandleFunc("/{name}/whois",
		injectDomainName(paginator.Paginate(w.Search)),
	)

	rec := record.Record{}
	sr.HandleFunc("/{duuid:"+UUID_REGEX+"}/records",
		injectDomainUUID(paginator.Paginate(rec.Search)),
	)
	sr.HandleFunc("/{name}/records",
		injectDomainName(paginator.Paginate(rec.Search)),
	)
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

func (d Domain) GetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, t, err := tld.DetectDomainAndTLD(vars["name"])
	if err != nil {
		util.ToJSON(db.Domain{}, err, w)
		return
	}
	i, err := db.Get(db.GetByNameAndTLD(), s, t.ID)
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

func injectDomainName(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		vars := mux.Vars(r)
		s, t, err := tld.DetectDomainAndTLD(vars["name"])
		if err != nil {
			util.ToJSON(db.Domain{}, err, w)
			return
		}
		i, err := db.Get(db.GetByNameAndTLD(), s, t.ID)
		if err != nil {
			util.ToJSON(i, err, w)
			return
		}
		context.Set(r, "domain", i.UUID.String())
		fn(w, r)
	}
}

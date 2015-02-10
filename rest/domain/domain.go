package domain

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/rest/whois"
	"net/http"
	"net/url"
)

type Domain struct {
}

const UUID_REGEX string = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"

func Setup(r *mux.Router) {
	d := &Domain{}
	sr := r.PathPrefix("/domain").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(d.Search))

	sr.HandleFunc("/query/email/{query}", d.Query)
	sr.HandleFunc("/{uuid:"+UUID_REGEX+"}", d.GetUUID)
	sr.HandleFunc("/{name}", d.GetName)

	sr.HandleFunc("/{duuid:"+UUID_REGEX+"}/whois",
		injectDomainUUID(paginator.Paginate(whois.Search)),
	)
	sr.HandleFunc("/{name}/whois",
		injectDomainName(paginator.Paginate(whois.Search)),
	)

}

func (d Domain) Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]
	fmt.Println(query)
	list, err := db.GetByJoinWhoisEmails(query).GetAll()
	util.ToJSON(list, err, w)
}

func (d Domain) Search(w http.ResponseWriter, r *http.Request, query url.Values, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (d Domain) GetUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	i, err := db.GetByUUID(uuid).Get()
	util.ToJSON(i, err, w)
}

func (d Domain) GetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, t, err := tld.DetectDomainAndTLD(vars["name"])
	if err != nil {
		util.ToJSON(db.Domain{}, err, w)
		return
	}
	i, err := db.GetByNameAndTLD(s, t.ID).Get()
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
		i, err := db.GetByNameAndTLD(s, t.ID).Get()
		if err != nil {
			util.ToJSON(i, err, w)
			return
		}
		context.Set(r, "domain", i.UUID.String())
		fn(w, r)
	}
}

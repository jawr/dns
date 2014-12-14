package domain

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
)

type Domain struct {
}

func Setup(r *mux.Router) {
	d := &Domain{}
	sr := r.PathPrefix("/domain").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(d.Search))
	sr.HandleFunc("/{uuid}", d.GetUUID)
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

package domain

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domain"
	//"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
)

type Domain struct {
}

func Setup(r *mux.Router) {
	d := &Domain{}
	sr := r.PathPrefix("/domain").Subrouter()
	sr.HandleFunc("/hello", d.Hello)
	sr.HandleFunc("/", paginator.Paginate(d.List))
}

func (d Domain) Hello(w http.ResponseWriter, r *http.Request) {
	util.ToJSON([]string{"a", "b", "c", "d"}, nil, w)
}

func (d Domain) List(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

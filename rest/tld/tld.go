package tld

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
	"strconv"
)

type TLD struct {
}

func Setup(r *mux.Router) {
	t := &TLD{}
	sr := r.PathPrefix("/tld").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(t.Search))
	sr.HandleFunc("/{id}", t.GetID)
}

func (t TLD) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (t TLD) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		util.ToJSON(id, err, w)
	}
	i, err := db.Get(db.GetByID(), id)
	util.ToJSON(i, err, w)
}

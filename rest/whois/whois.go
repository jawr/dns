package whois

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
	"strconv"
)

type Record struct {
}

func Setup(r *mux.Router) {
	rec := &Record{}
	sr := r.PathPrefix("/whois").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(rec.Search))
	sr.HandleFunc("/{id}", rec.GetID)
}

func (rec Record) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (rec Record) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		util.ToJSON(id, err, w)
	}
	i, err := db.Get(db.GetByID(), id)
	util.ToJSON(i, err, w)
}

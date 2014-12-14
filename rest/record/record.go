package record

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
)

type Record struct {
}

func Setup(r *mux.Router) {
	rec := &Record{}
	sr := r.PathPrefix("/record").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(rec.Search))
	sr.HandleFunc("/{uuid}", rec.GetUUID)
}

func (rec Record) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (rec Record) GetUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	i, err := db.Get(db.GetByUUID(), uuid)
	util.ToJSON(i, err, w)
}

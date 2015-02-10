package record_type

import (
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/record_type"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
	"net/url"
	"strconv"
)

type RecordType struct {
}

func Setup(r *mux.Router) {
	rec := &RecordType{}
	sr := r.PathPrefix("/record_type").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(rec.Search))
	sr.HandleFunc("/{id:[0-9]+}", rec.GetID)
	sr.HandleFunc("/{name}", rec.GetName)
}

func (rec RecordType) Search(w http.ResponseWriter, r *http.Request, query url.Values, idx, limit int) {
	list, err := db.Search(query, idx, limit)
	util.ToJSON(list, err, w)
}

func (rec RecordType) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		util.ToJSON(id, err, w)
	}
	i, err := db.Get(db.GetByID(), id)
	util.ToJSON(i, err, w)
}

func (rec RecordType) GetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	i, err := db.Get(db.GetByName(), name)
	util.ToJSON(i, err, w)
}

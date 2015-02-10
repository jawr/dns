package watcher

import (
	//"encoding/json"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/watcher"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
	"net/url"
	"strconv"
)

type Result struct {
}

func Setup(r *mux.Router) {
	res := &Result{}
	sr := r.PathPrefix("/watcher").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(res.Root))
	sr.HandleFunc("/{id}", res.GetID)
}

func (res Result) Root(w http.ResponseWriter, r *http.Request, query url.Values, idx, limit int) {
	switch r.Method {
	case "GET":
		list, err := db.Search(query, limit, idx)
		log.Info("watcher: %+v", list)
		util.ToJSON(list, err, w)
	case "POST":
		res.Post(w, r)
	}
}

func (res Result) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		util.ToJSON(id, err, w)
		return
	}
	i, err := db.Get(db.GetByID(), id)
	util.ToJSON(i, err, w)
}

func (res Result) Post(w http.ResponseWriter, r *http.Request) {
	log.Info("Do post...")
	util.ToJSON(r, nil, w)
}

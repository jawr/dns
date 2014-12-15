package whois

import (
	"errors"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domain"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/whois/parser"
	"net/http"
	"strconv"
)

type Result struct {
}

func Setup(r *mux.Router) {
	res := &Result{}
	sr := r.PathPrefix("/whois").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(res.Search))
	sr.HandleFunc("/{id}", res.GetID)
}

func (res Result) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	switch r.Method {
	case "GET":
		list, err := db.Search(query, idx, limit)
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
	uuid, ok := context.GetOk(r, "domain")
	if !ok {
		util.ToJSON(db.Result{}, errors.New("No object found."), w)
		return
	}
	d, err := domain.Get(domain.GetByUUID(), uuid)
	if err != nil {
		util.ToJSON(db.Result{}, err, w)
		return
	}

	p := parser.New()
	wret, err := p.Parse(d)
	util.ToJSON(wret, err, w)
}

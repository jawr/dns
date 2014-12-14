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

type Whois struct {
}

func Setup(r *mux.Router) {
	wh := &Whois{}
	sr := r.PathPrefix("/whois").Subrouter()
	sr.HandleFunc("/", paginator.Paginate(wh.Search))
	sr.HandleFunc("/{id}", wh.GetID)
}

func (wh Whois) Search(w http.ResponseWriter, r *http.Request, query map[string][]string, idx, limit int) {
	switch r.Method {
	case "GET":
		list, err := db.Search(query, idx, limit)
		util.ToJSON(list, err, w)
	case "POST":
		wh.Post(w, r)
	}
}

func (wh Whois) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		util.ToJSON(id, err, w)
		return
	}
	i, err := db.Get(db.GetByID(), id)
	util.ToJSON(i, err, w)
}

func (wh Whois) Post(w http.ResponseWriter, r *http.Request) {
	uuid, ok := context.GetOk(r, "domain")
	if !ok {
		util.ToJSON(db.Whois{}, errors.New("No object found."), w)
		return
	}
	d, err := domain.Get(domain.GetByUUID(), uuid)
	if err != nil {
		util.ToJSON(db.Whois{}, err, w)
		return
	}

	p := parser.New()
	wret, err := p.Parse(d)
	util.ToJSON(wret, err, w)
}

package whois

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domain"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/whois/dispatcher"
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

type Post struct {
	UUID  string `json:"uuid,omitempty"`
	Query string `json:"query,omitempty"`
}

func (res Result) Post(w http.ResponseWriter, r *http.Request) {
	log.Info("%+v", r.Body)
	decoder := json.NewDecoder(r.Body)
	var post Post
	err := decoder.Decode(&post)
	if err != nil {
		util.ToJSON(db.Result{}, err, w)
		return
	}
	log.Info("%+v", post)
	if len(post.UUID) > 0 {
		d, err := domain.Get(domain.GetByUUID(), post.UUID)
		if err != nil {
			util.ToJSON(db.Result{}, err, w)
			return
		}
		c := dispatcher.AddDomain(d)
		result := <-c
		util.ToJSON(result, err, w)
		return

	} else if len(post.Query) > 0 {
		c := dispatcher.AddQuery(post.Query)
		result := <-c
		util.ToJSON(result, err, w)
		return
	}
	// do error

}

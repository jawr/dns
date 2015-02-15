package watchers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domains"
	db "github.com/jawr/dns/database/models/watchers"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"net/http"
	"net/url"
	"strings"
)

var routes = util.Routes{
	util.Route{
		"Index",
		"GET",
		"/",
		paginator.Paginate(Search),
	},
	util.Route{
		"Create",
		"POST",
		"/",
		Create,
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/watchers").Subrouter()
	util.SetupRouter(subRouter, "Watchers", routes)
}

/*
	Search is used for retrieving watchers. It accepts the following
	GET params:

	+ duuid 	- uuid of domain
	+ domain	- string of domain
	+ name		- string of name

*/
func Search(w http.ResponseWriter, r *http.Request, params url.Values, limit, offset int) {
	query := db.SELECT
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		case "name":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "duuid", "domain":
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, params.Get(k))
			i++
		}
	}
	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + " "
	}
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	log.Info("Query: " + query)
	log.Info("Args: %+v", args)
	watcherList, err := db.GetList(query, args...)
	util.ToJSON(watcherList, err, w)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var post struct {
		Domain   string `json:"domain"`
		Interval string `json:"interval"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		util.Error(err, w)
		return
	}
	domain, err := domains.GetByUUID(post.Domain).One()
	if err != nil {
		util.Error(err, w)
		return
	}
	watcher, err := db.New(domain, post.Interval)
	util.ToJSON(watcher, err, w)
}

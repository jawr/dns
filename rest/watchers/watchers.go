package watchers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/tlds"
	db "github.com/jawr/dns/database/models/watchers"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/auth"
	"github.com/jawr/dns/rest/paginator"
	"github.com/jawr/dns/rest/util"
	"github.com/jawr/dns/whois/dispatcher"
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
		case "user":
			user, err := auth.GetUser(r)
			if err != nil {
				util.Error(err, w)
				return
			}
			where = append(where, fmt.Sprintf("users@> '[%d]'", user.ID))
			i++
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
		Name     string `json:"name"`
		Interval string `json:"interval"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		util.Error(err, w)
		return
	}
	var domain domains.Domain
	if len(post.Name) > 0 {
		name, tld, err := tlds.DetectDomainAndTLD(post.Name)
		if err != nil {
			util.Error(err, w)
			return
		}
		domain, err = domains.GetByNameAndTLD(name, tld).One()
		if err != nil {
			domain = domains.New(name, tld)
			err = domain.Insert()
			if err != nil {
				util.Error(err, w)
				return
			}
		}
	} else {
		domain, err = domains.GetByUUID(post.Domain).One()
		if err != nil {
			util.Error(err, w)
			return
		}
	}
	watcher, err := db.GetByDomain(domain).One()
	if err != nil {
		watcher, err = db.New(domain, post.Interval)
		if err != nil {
			util.Error(err, w)
			return
		}
	}
	watcher.SetLowerInterval(post.Interval)
	user, err := auth.GetUser(r)
	if err != nil {
		util.Error(err, w)
		return
	}
	watcher.AddUser(*user)
	err = watcher.Save()
	log.Info("%+v", watcher)
	dispatcher.AddDomain(domain)
	util.ToJSON(watcher, err, w)
}

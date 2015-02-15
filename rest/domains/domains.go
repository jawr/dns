package domains

import (
	"fmt"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/tlds"
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
		"ByUUID",
		"GET",
		"/{uuid:" + util.UUID_REGEX + "}",
		ByUUID(serveDomain),
	},
	util.Route{
		"ByName",
		"GET",
		"/{name}",
		ByName(serveDomain),
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/domain").Subrouter()
	util.SetupRouter(subRouter, "Domain", routes)
}

/*
	Search is used for retrieving domain records. It accepts the following GET params:

	+ uuid		- uuid of domain
	+ name		- string representation of domain (can include wildcards)
	+ email		- domains attached to an email
*/
func Search(w http.ResponseWriter, r *http.Request, params url.Values, limit, offset int) {
	query := db.SELECT
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		case "uuid":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
		case "email":
			// special case which overrides results
			domainList, err := db.GetByJoinWhoisEmails(params.Get(k)).List()
			util.ToJSON(domainList, err, w)
			return
		case "name":
			name := params.Get(k)
			if strings.ContainsAny(name, ".") {
				// attempt to detect tld
				tld, err := tlds.Detect(name)
				if err == nil {
					where = append(where, fmt.Sprintf("tld = $%d", i))
					args = append(args, tld.ID)
					// strip tld from name
					name = strings.TrimSuffix(name, "."+tld.Name)
				}
			}
			if strings.ContainsAny(name, "* || % ||  ") {
				// contains wildcards
				name = strings.Replace(name, "*", "%", -1)
				name = strings.Replace(name, " ", "%", -1)
			}
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, name)
		}
	}
	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + " "
	}
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	log.Info("Query: " + query)
	log.Info("Args: %+v", args)
	domainList, err := db.GetList(query, args...)
	util.ToJSON(domainList, err, w)
}

func ByUUID(fn func(http.ResponseWriter, *http.Request, db.Domain)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		domain, err := db.GetByUUID(uuid).One()
		if err != nil {
			util.Error(err, w)
			return
		}
		fn(w, r, domain)
	}
}

func ByName(fn func(http.ResponseWriter, *http.Request, db.Domain)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, tld, err := tlds.DetectDomainAndTLD(vars["name"])
		if err != nil {
			util.Error(err, w)
			return
		}
		domain, err := db.GetByNameAndTLD(name, tld.ID).One()
		fn(w, r, domain)
	}
}

func serveDomain(w http.ResponseWriter, r *http.Request, domain db.Domain) {
	util.ToJSON(domain, nil, w)
}

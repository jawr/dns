package whois

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/tlds"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	domainsAPI "github.com/jawr/dns/rest/domains"
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
}

var overloadRoutes = util.Routes{
	util.Route{
		"ByDomain",
		"GET",
		"/domain/{uuid:" + util.UUID_REGEX + "}/whois",
		domainsAPI.ByUUID(ByDomainUUID),
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/whois").Subrouter()
	util.SetupRouter(subRouter, "Whois", routes)
	util.SetupRouter(router, "WhoisOverload", overloadRoutes)
}

/*
	Search is used for retrieving whois records for a domain. It accepts the following
	GET params:

	+ uuid		- uuid of record
	+ duuid 	- uuid of domain
	+ domain	- string of domain

*/
func Search(w http.ResponseWriter, r *http.Request, params url.Values, limit, offset int) {
	query := db.SELECT
	var domain domains.Domain
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		case "uuid":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "duuid", "domain":
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, params.Get(k))
			v := params.Get(k)
			params.Del("duuid")
			params.Del("domain")
			params.Set("domain", v)
			i++
		case "name":
			name, tld, err := tlds.DetectDomainAndTLD(params.Get(k))
			if err != nil {
				util.Error(err, w)
				return
			}
			domain, err = domains.GetByNameAndTLD(name, tld).One()
			if err != nil {
				util.Error(err, w)
				return
			}
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, domain.UUID.String())
			i++
		case "email":
			where = append(where, fmt.Sprintf("emails ? $%d", i))
			args = append(args, params.Get(k))
			i++
		case "raw":
			where = append(where, fmt.Sprintf("raw_whois ->>0 ILIKE $%d", i))
			args = append(args, "%"+params.Get(k)+"%")
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
	recordList, err := db.GetList(query, args...)
	// check for non sql errors
	// if we have no results dispatch a worker to get one
	if len(recordList) == 0 {
		if len(domain.Name) == 0 {
			// no domain lets grab one using what we assume is a duuid
			if duuid := params.Get("domain"); duuid != "" {
				log.Info("duuid: " + duuid)
				domain, err = domains.GetByUUID(duuid).One()
				if err != nil {
					util.Error(err, w)
					return
				}
			}
		}
		if len(domain.Name) == 0 {
			log.Error("Unable to detect domain for whois lookup, params: %+v", params)
			util.Error(errors.New("Unable to find domain for whois lookup."), w)
			return
		}
		result := <-dispatcher.AddDomain(domain)
		recordList = append(recordList, result)
	}
	util.ToJSON(recordList, err, w)
}

func ByDomainUUID(w http.ResponseWriter, r *http.Request, domain domains.Domain) {
	record, err := db.GetByDomain(domain).One()
	// could push to dispatcher based on query params
	util.ToJSON(record, err, w)
}

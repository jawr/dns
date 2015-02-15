package records

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domains"
	db "github.com/jawr/dns/database/models/records"
	"github.com/jawr/dns/dig/dispatcher"
	"github.com/jawr/dns/log"
	domainsAPI "github.com/jawr/dns/rest/domains"
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
}

var overloadRoutes = util.Routes{
	util.Route{
		"ByDomain",
		"GET",
		"/domain/{uuid:" + util.UUID_REGEX + "}/records",
		domainsAPI.ByUUID(ByDomainUUID),
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/records").Subrouter()
	util.SetupRouter(subRouter, "Records", routes)
	util.SetupRouter(router, "RecordsOverload", overloadRoutes)
}

type Domain map[string]interface{}

/*
	Search is used for retrieving records for a domain. It accepts the following
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
			v := params.Get(k)
			params.Del("duuid")
			params.Del("domain")
			params.Set("domain", v)
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
	if err != nil {
		util.Error(err, w)
		return
	}
	// check for none sql errors
	// if we have no results dispatch a worker to get one
	if len(recordList) == 0 {
		// no domain lets grab one using what we assume is a duuid
		log.Info("no records")
		if duuid := params.Get("domain"); duuid != "" {
			log.Info("duuid: " + duuid)
			domain, err := domains.GetByUUID(duuid).One()
			if err != nil {
				log.Info("herererere")
				util.Error(err, w)
				return
			}
			results := <-dispatcher.AddDomain(domain)
			log.Info("%v", results)
			for _, result := range results {
				recordList = append(recordList, result)
			}
		}
	}
	util.ToJSON(cleanParserFromRecords(recordList), err, w)
}

func ByDomainUUID(w http.ResponseWriter, r *http.Request, domain domains.Domain) {
	recordList, err := db.GetByDomain(domain).List()
	// could push to dispatcher based on query params
	util.ToJSON(cleanParserFromRecords(recordList), err, w)
}

func cleanParserFromRecords(recordList []db.Record) []Domain {
	var parsed []Domain
	structs.DefaultTagName = "json"
	for _, r := range recordList {
		m := structs.Map(r)
		m["parse_date"] = r.Date
		m["added"] = r.Added
		//m["domain_uuid"] = domain.UUID.String()
		//m["domain"] = domain.String()
		delete(m, "parser")
		parsed = append(parsed, m)
	}
	return parsed
}

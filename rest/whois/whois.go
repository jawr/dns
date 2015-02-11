package whois

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	domains "github.com/jawr/dns/database/models/domain"
	tlds "github.com/jawr/dns/database/models/tld"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
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

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/whois").Subrouter()
	util.SetupRouter(subRouter, "Whois", routes)
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
			i++
		case "name":
			name, tld, err := tlds.DetectDomainAndTLD(params.Get(k))
			if err != nil {
				util.Error(err, w)
				return
			}
			domain, err = domains.GetByNameAndTLD(name, tld.ID).Get()
			if err != nil {
				util.Error(err, w)
				return
			}
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, domain.UUID.String())
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
	records, err := db.GetList(query, args...)
	if err != nil {
		log.Error("Here")
		util.Error(err, w)
		return
	}
	// if we have no results dispatch a worker to get one
	if len(records) == 0 {
		if len(domain.Name) == 0 {
			// no domain lets grab one using what we assume is a duuid
			if duuid := params.Get("domain"); duuid != "" {
				log.Info("duuid: " + duuid)
				domain, err = domains.GetByUUID(duuid).Get()
				if err != nil {
					util.Error(err, w)
					return
				}
			}
		}
		if len(domain.Name) == 0 {
			log.Error("Unable to detect domain for whois lookup, params: %+v", params)
			util.Error(errors.New("Unable to detect domain for whois lookup."), w)
			return
		}
		result := <-dispatcher.AddDomain(domain)
		records = append(records, result)
	}
	util.ToJSON(records, err, w)
}

/*
func Post(w http.ResponseWriter, r *http.Request) {
	log.Info("%+v", r.Body)
	decoder := json.NewDecoder(r.Body)
	var post struct {
		UUID  string `json:"uuid,omitempty"`
		Query string `json:"query,omitempty"`
	}
	err := decoder.Decode(&post)
	if err != nil {
		util.ToJSON(db.Result{}, err, w)
		return
	}
	log.Info("%+v", post)
	if len(post.UUID) > 0 {
		d, err := domain.GetByUUID(post.UUID).Get()
		if err != nil {
			util.ToJSON(db.Result{}, err, w)
			return
		}
		c := dispatcher.AddDomain(d)
		result := <-c
		util.ToJSON([]db.Result{result}, err, w)
		return

	} else if len(post.Query) > 0 {
		c := dispatcher.AddQuery(post.Query)
		result := <-c
		util.ToJSON([]db.Result{result}, err, w)
		return
	}
	// do error
}
*/

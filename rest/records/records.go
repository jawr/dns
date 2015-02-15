package records

import (
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/database/models/domains"
	db "github.com/jawr/dns/database/models/records"
	domainsAPI "github.com/jawr/dns/rest/domains"
	"github.com/jawr/dns/rest/util"
	"net/http"
)

var overloadRoutes = util.Routes{
	util.Route{
		"ByDomain",
		"GET",
		"/domain/{uuid:" + util.UUID_REGEX + "}/records",
		domainsAPI.ByUUID(ByDomainUUID),
	},
}

func Setup(router *mux.Router) {
	//subRouter := router.PathPrefix("/records").Subrouter()
	util.SetupRouter(router, "RecordsOverload", overloadRoutes)
}

type Domain map[string]interface{}

func ByDomainUUID(w http.ResponseWriter, r *http.Request, domain domains.Domain) {
	recordList, err := db.GetByDomain(domain).List()
	// could push to dispatcher based on query params
	var parsed []Domain
	structs.DefaultTagName = "json"
	for _, r := range recordList {
		m := structs.Map(r)
		m["type"] = r.Type.Name
		//m["domain_uuid"] = domain.UUID.String()
		//m["domain"] = domain.String()
		delete(m, "parser")
		parsed = append(parsed, m)
	}
	util.ToJSON(parsed, err, w)
}

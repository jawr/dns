package notifications

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	db "github.com/jawr/dns/database/models/notifications"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/auth"
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
		"Archive",
		"POST",
		"/archive/",
		Archive,
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/notifications").Subrouter()
	util.SetupRouter(subRouter, "Notifications", routes)
}

func Archive(w http.ResponseWriter, r *http.Request) {
	var post []db.Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		util.Error(err, w)
		return
	}
	user, err := auth.GetUser(r)
	if err != nil {
		util.Error(err, w)
		return
	}
	note, err := db.GetByUser(*user).One()
	if err != nil {
		util.Error(err, w)
		return
	}
	note.ArchiveMessages(post)
	util.ToJSON(post, nil, w)
}

/*
	Search is used for retrieving notifications. It accepts the following
	GET params:


*/
func Search(w http.ResponseWriter, r *http.Request, params url.Values, limit, offset int) {
	query := db.SELECT
	var where []string
	var args []interface{}
	var orderBy = ""
	i := 1
	// inject user
	user, err := auth.GetUser(r)
	if err != nil {
		util.Error(err, w)
		return
	}
	where = append(where, fmt.Sprintf("user_id = $%d", i))
	args = append(args, user.ID)
	i++
	for k, _ := range params {
		switch k {
		case "duuid":
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "orderBy":
			orderBy = fmt.Sprintf("ORDER BY $%d ", i)
			args = append(args, params.Get(k))
			if v, ok := params["order"]; ok {
				order := strings.ToLower(v[0])
				if order == "desc" {
					orderBy += "DESC "
				} else if order == "asc" {
					orderBy += "ASC "
				}
			}
		}
	}
	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + " "
	}
	if len(orderBy) > 0 {
		query += orderBy + " "
	}
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	log.Info("Query: " + query)
	log.Info("Args: %+v", args)
	notificationList, err := db.GetList(query, args...)
	util.ToJSON(notificationList, err, w)
}

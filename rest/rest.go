package rest

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/jawr/dns/rest/auth"
	"github.com/jawr/dns/rest/domains"
	"github.com/jawr/dns/rest/notifications"
	"github.com/jawr/dns/rest/records"
	"github.com/jawr/dns/rest/watchers"
	"github.com/jawr/dns/rest/whois"
	"github.com/yosssi/boltstore/reaper"
	"net/http"
)

func Setup() {
	// setup sessions
	sessionDB, err := bolt.Open("./sessions.db", 0666, nil)
	if err != nil {
		panic(err)
	}
	defer sessionDB.Close()
	defer reaper.Quit(reaper.Run(sessionDB, reaper.Options{}))

	// setup routes
	r := mux.NewRouter()
	sr := r.PathPrefix("/api/v1").Subrouter()
	domains.Setup(sr)
	whois.Setup(sr)
	records.Setup(sr)
	watchers.Setup(sr)
	auth.Setup(sr)
	notifications.Setup(sr)

	// setup authoriser
	authoriser := auth.New(sessionDB)
	h := authoriser(r)
	http.ListenAndServe(":8080", h)
}

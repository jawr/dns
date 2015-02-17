package auth

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jawr/dns/database/models/users"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest/util"
	"github.com/stathat/jconfig"
	"github.com/yosssi/boltstore/store"
	"net/http"
	"strings"
	"time"
)

type httpHandler func(http.Handler) http.Handler

var sessionName string
var str *store.Store

func init() {
	gob.Register(time.Now())
}

var routes = util.Routes{
	util.Route{
		"Index",
		"GET",
		"/",
		CheckAuth,
	},
}

func Setup(router *mux.Router) {
	subRouter := router.PathPrefix("/auth").Subrouter()
	util.SetupRouter(subRouter, "Auth", routes)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	var params = r.URL.Query()
	if check := params.Get("check"); len(check) > 0 {
		w.WriteHeader(http.StatusOK)
		return
	}
	var hash = params.Get("hash")
	if len(hash) == 0 {
		util.Error(errors.New("No hash supplied for auth."), w)
		return
	}
	data, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		util.Error(err, w)
		return
	}
	args := strings.Split(string(data), ":")
	fmt.Println(args)
	if len(args) == 2 && users.CheckPassword(args[0], args[1]) {
		session, err := GetSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["expires"] = time.Now().Add(time.Minute * 10)
		if err := sessions.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	util.Error(errors.New("Unable to authenticate."), w)

}

func New(db *bolt.DB) httpHandler {
	config := jconfig.LoadConfig("config.json")
	key := config.GetString("session_secret")
	var err error
	str, err = store.New(db, store.Config{}, []byte(key))
	if err != nil {
		panic(err)
	}
	sessionName = config.GetString("session_name")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				session, err := GetSession(r)
				if err != nil {
					util.Error(err, w)
					return
				}
				if t, ok := session.Values["expires"]; ok {
					log.Info("Got existing cookie...")
					if err == nil && time.Now().Before(t.(time.Time)) {
						// authed
						w.Header().Set("X-Expires", t.(time.Time).String())
						h.ServeHTTP(w, r)
						return
					}
				}
				user, pass, ok := r.BasicAuth()
				if ok {
					log.Info("Got basic auth...")
					if users.CheckPassword(user, pass) {
						session.Values["expires"] = time.Now().Add(time.Minute * 10)
						// save our session
						if err := sessions.Save(r, w); err != nil {
							util.Error(err, w)
							return
						}
						h.ServeHTTP(w, r)
						return
					}
				}
				if r.URL.Path == "/api/v1/auth/" {
					h.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
			},
		)
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return str.Get(r, sessionName)
}

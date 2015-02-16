package auth

import (
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"github.com/jawr/dns/database/models/users"
	"github.com/stathat/jconfig"
	"github.com/yosssi/boltstore/store"
	"net/http"
	"time"
)

type httpHandler func(http.Handler) http.Handler

var sessionName string
var str *store.Store

func init() {
	gob.Register(time.Now())
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
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if t, ok := session.Values["expires"]; ok {
					if err == nil && time.Now().Before(t.(time.Time)) {
						if err := sessions.Save(r, w); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						h.ServeHTTP(w, r)
						return
					}
				}
				user, pass, ok := r.BasicAuth()
				if ok {
					if users.CheckPassword(user, pass) {
						session.Values["expires"] = time.Now().Add(time.Minute * 10)
						if err := sessions.Save(r, w); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						h.ServeHTTP(w, r)
						return
					}
				}
				w.WriteHeader(http.StatusUnauthorized)
			},
		)
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return str.Get(r, sessionName)
}

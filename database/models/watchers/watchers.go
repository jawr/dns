package watchers

import (
	"encoding/json"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/users"
	"github.com/jawr/dns/database/models/watchers/intervals"
	"time"
)

type Log struct {
	Time time.Time `json:"time"`
}

type Watcher struct {
	ID       int32              `json:"id"`
	Domain   domains.Domain     `json:"domain"`
	Added    time.Time          `json:"added"`
	Updated  time.Time          `json:"updated"`
	Interval intervals.Interval `json:"interval"`
	Logs     []Log              `json:"logs"`
	Users    []users.User       `json:"users"`
}

func New(d domains.Domain, interval string) (Watcher, error) {
	w := Watcher{}
	w.Domain = d
	i, err := intervals.New(interval)
	if err != nil {
		// TODO: better wrapper for our db errors so can see in log which package
		return w, err
	}
	w.Interval = i
	w.Updated = time.Now()
	w.Added = time.Now()
	conn, err := connection.Get()
	if err != nil {
		return w, err
	}
	err = conn.QueryRow("INSERT INTO watcher (domain, interval) VALUES ($1, $2) RETURNING id",
		d.UUID.String(),
		i.ID,
	).Scan(&w.ID)
	if err != nil {
		return w, err
	}
	return GetByDomain(d).One()
}

func (w *Watcher) Save() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	w.Logs = append(w.Logs, Log{w.Updated})
	w.Updated = time.Now()
	b, err := json.Marshal(w.Logs)
	if err != nil {
		return err
	}
	var usersArr []int32
	for _, u := range w.Users {
		usersArr = append(usersArr, u.ID)
	}
	b2, err := json.Marshal(usersArr)
	if err != nil {
		return err
	}
	_, err = conn.Exec("UPDATE watcher SET updated = $1, logs = $2, users = $3 WHERE id = $4", w.Updated, b, b2, w.ID)
	return err
}

func (w *Watcher) AddUser(user users.User) {
	for _, u := range w.Users {
		if user.ID == u.ID {
			return
		}
	}
	w.Users = append(w.Users, user)
}

func (w *Watcher) AddUserByString(userStr string) {
	user, err := users.GetByEmail(userStr).One()
	if err != nil {
		return
	}
	for _, u := range w.Users {
		if user.ID == u.ID {
			return
		}
	}
	w.Users = append(w.Users, user)
}

func (w *Watcher) SetLowerInterval(intervalStr string) {
	interval, err := intervals.GetByValue(intervalStr).One()
	if err != nil {
		interval, err = intervals.New(intervalStr)
		if err != nil {
			return
		}
	}
	// BUG(setLowerInterval) implement this, make it sort intervals perhaps using the cache
	// and set the watcher to use the smaller interval
	fmt.Println("%s", interval.Value)
}

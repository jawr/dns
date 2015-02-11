package watchers

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
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
}

func New(d domains.Domain, i intervals.Interval) (Watcher, error) {
	w := Watcher{}
	w.Domain = d
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
		return GetByDomain(d.UUID.String()).One()
	}
	return w, err
}

func (w *Watcher) Save() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	w.Logs = append(w.Logs, Log{w.Updated})
	b, err := json.Marshal(w.Logs)
	if err != nil {
		return err
	}
	_, err = conn.Exec("UPDATE watcher SET logs = $1 WHERE id = $2", b, w.ID)
	return err
}

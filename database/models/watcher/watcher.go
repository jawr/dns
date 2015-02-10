package watcher

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/watcher/interval"
	"time"
)

type Log struct {
	Time time.Time `json:"time"`
}

type Watcher struct {
	ID       int32             `json:"id"`
	Domain   domain.Domain     `json:"domain"`
	Added    time.Time         `json:"added"`
	Updated  time.Time         `json:"updated"`
	Interval interval.Interval `json:"interval"`
	Logs     []Log             `json:"logs"`
}

func New(d domain.Domain, i interval.Interval) (Watcher, error) {
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
		return Get(GetByDomain(), d.UUID.String())
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

const (
	SELECT string = "SELECT * FROM watcher "
)

func GetAll() string {
	return SELECT
}

func GetByID() string {
	return SELECT + "WHERE id = $1"
}

func GetByInterval() string {
	return SELECT + "WHERE interval = $1"
}

func GetByDomain() string {
	return SELECT + "WHERE domain = $1"
}

func parseRow(row connection.Row) (Watcher, error) {
	w := Watcher{}
	var intervalID int32
	var uuid string
	var b []byte
	err := row.Scan(&w.ID, &uuid, &w.Added, &w.Updated, &intervalID, &b)
	err = json.Unmarshal(b, &w.Logs)
	if err != nil {
		return w, err
	}
	w.Domain, err = domain.GetByUUID(uuid).Get()
	if err != nil {
		return w, err
	}
	w.Interval, err = interval.Get(interval.GetByID(), intervalID)
	return w, err
}

func Get(query string, args ...interface{}) (Watcher, error) {
	var result Watcher
	conn, err := connection.Get()
	if err != nil {
		return Watcher{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err = parseRow(row)
	return result, err
}

func GetList(query string, args ...interface{}) ([]Watcher, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Watcher{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Watcher{}, err
	}
	var list []Watcher
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

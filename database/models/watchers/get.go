package watchers

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/users"
	"github.com/jawr/dns/database/models/watchers/intervals"
)

const (
	SELECT string = "SELECT * FROM watcher "
)

type Result struct {
	One  func() (Watcher, error)
	List func() ([]Watcher, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Watcher, error) {
			return Get(query, args...)
		},
		func() ([]Watcher, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByUser(user users.User) Result {
	return newResult(SELECT+"WHERE users@> '[$1]'", user.ID)
}

func GetByID(id int32) Result {
	return newResult(SELECT+"WHERE id = $1", id)
}

func GetByInterval(interval intervals.Interval) Result {
	return newResult(SELECT+"WHERE interval = $1", interval.ID)
}

func GetByDomainUUID(uuid string) Result {
	return newResult(SELECT+"WHERE domain = $1", uuid)
}

func GetByDomain(domain domains.Domain) Result {
	return GetByDomainUUID(domain.UUID.String())
}

func parseRow(row connection.Row) (Watcher, error) {
	w := Watcher{}
	var intervalID int32
	var uuid string
	var usersBuf, logs []byte
	err := row.Scan(&w.ID, &uuid, &w.Added, &w.Updated, &intervalID, &logs, &usersBuf)
	err = json.Unmarshal(logs, &w.Logs)
	if err != nil {
		return w, err
	}
	var usersArr []int
	err = json.Unmarshal(usersBuf, &usersArr)
	if err != nil {
		return w, err
	}
	for _, u := range usersArr {
		user, err := users.GetByID(u).One()
		if err != nil {
			continue
		}
		w.Users = append(w.Users, user)
	}
	w.Domain, err = domains.GetByUUID(uuid).One()
	if err != nil {
		return w, err
	}
	w.Interval, err = intervals.GetByID(intervalID).One()
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

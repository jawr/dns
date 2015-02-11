package whois

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
)

const (
	SELECT string = "SELECT * FROM whois "
)

type Result struct {
	One  func() (Record, error)
	List func() ([]Record, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Record, error) {
			return Get(query, args...)
		},
		func() ([]Record, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByID(id int) Result {
	return newResult(SELECT+"WHERE id = $1", id)
}

func GetByHasEmail(email string) Result {
	return newResult(SELECT+"WHERE emails ? $1", email)
}

func parseRow(row connection.Row) (Record, error) {
	w := Record{}
	var duuid, wuuid string
	var b []byte
	err := row.Scan(&w.ID, &duuid, &w.Data, &w.Raw, &w.Contacts, &b, &w.Added, &wuuid)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(b, &w.Emails)
	if err != nil {
		return w, err
	}
	w.Domain, err = domains.GetByUUID(duuid).One()
	if err != nil {
		return w, err
	}
	w.UUID = uuid.Parse(wuuid)
	return w, nil
}

func Get(query string, args ...interface{}) (Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return Record{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Record{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Record{}, err
	}
	var list []Record
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

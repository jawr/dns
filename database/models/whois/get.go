package whois

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"strings"
)

const (
	SELECT string = "SELECT * FROM whois "
)

type Result struct {
	Count func() (int, error)
	One   func() (Record, error)
	List  func() ([]Record, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (int, error) {
			var count int
			conn, err := connection.Get()
			if err != nil {
				return count, err
			}
			query = strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
			err = conn.QueryRow(query, args...).Scan(&count)
			return count, err
		},
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

func GetByDomain(domain domains.Domain) Result {
	return newResult(SELECT+"WHERE domain = $1", domain.UUID.String())
}

func parseRow(row connection.Row) (Record, error) {
	w := Record{}
	var duuid, wuuid string
	var emails, organizations, phones, postcodes, names []byte
	err := row.Scan(
		&w.ID,
		&duuid,
		&w.Data,
		&w.Raw,
		&w.Contacts,
		&emails,
		&w.Added,
		&wuuid,
		&organizations,
		&phones,
		&postcodes,
		&names,
	)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(emails, &w.Emails)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(organizations, &w.Organizations)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(phones, &w.Phones)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(postcodes, &w.Postcodes)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(names, &w.Emails)
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

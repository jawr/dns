package domain

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tld"
)

const (
	SELECT string = "SELECT * FROM domain "
)

func GetByNameAndTLD() string {
	return SELECT + "WHERE name = $1 AND tld = $2"
}

func GetByUUID() string {
	return SELECT + "WHERE uuid = $1"
}

func GetAll() string {
	return SELECT + "LIMIT 10"
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var dUUID string
	var tldId int32
	err := row.Scan(&dUUID, &d.Name, &tldId)
	if err != nil {
		return d, err
	}
	d.UUID = uuid.Parse(dUUID)
	t, err := tld.Get(tld.GetByID(), tldId)
	d.TLD = t
	return d, nil
}

func Get(query string, args ...interface{}) (Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return Domain{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Domain{}, err
	}
	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return []Domain{}, err
	}
	var list []Domain
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

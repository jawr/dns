package domain

import (
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/tld"
)

func GetByNameAndTLD() string {
	return "SELECT * FROM domain WHERE name = $1 AND tld = $2"
}

func GetByID() string {
	return "SELECT * FROM domain WHERE id = $1"
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var tldId int32
	err := row.Scan(&d.UUID, &d.Name, &tldId)
	if err != nil {
		return d, err
	}
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

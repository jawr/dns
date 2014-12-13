package whois

import (
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
)

func parseRow(row connection.Row) (Whois, error) {
	w := Whois{}
	var uuid string
	err := row.Scan(&w.ID, &uuid, &w.Data, &w.Added)
	if err != nil {
		return w, err
	}
	d, err := domain.Get(domain.GetByUUID(), uuid)
	if err != nil {
		return w, err
	}
	w.Domain = d
	return w, nil
}

func Get(query string, args ...interface{}) (Whois, error) {
	conn, err := connection.Get()
	if err != nil {
		return Whois{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Whois, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Whois{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Whois{}, err
	}
	var list []Whois
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

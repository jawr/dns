package domain

import (
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/tld"
	"strings"
)

type Domain struct {
	ID   int32
	Name string
	TLD  tld.TLD
}

func New(name string, t tld.TLD) (Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return Domain{}, err
	}
	name = strings.TrimSuffix(name, ".")
	name = strings.TrimSuffix(name, "."+t.Name)
	var id int32
	err = conn.QueryRow("INSERT INTO domain (name, tld) VALUES ($1, $2) RETURNING id",
		name,
		t.ID,
	).Scan(&id)
	return Domain{
		ID:   id,
		Name: name,
		TLD:  t,
	}, err
}

func GetByNameAndTLD() string {
	return "SELECT * FROM domain WHERE name = $1 AND tld = $2"
}

func GetByID(id int32) (Domain, error) {
	return Get("SELECT * FROM domain WHERE id = $1", id)
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var tldId int32
	err := row.Scan(&d.ID, &d.Name, &tldId)
	if err != nil {
		return d, err
	}
	t, err := tld.GetByID(tldId)
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

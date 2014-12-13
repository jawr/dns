package tld

import (
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
)

type TLD struct {
	ID   int32
	Name string
}

var c = cache.New()

func New(name string) (TLD, error) {
	if t, ok := c.Check(name); ok {
		return t.(TLD), nil
	}
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_tld($1)", name).Scan(&id)
	t := TLD{
		ID:   id,
		Name: name,
	}
	c.Add(t)
	return t, err
}

func (t TLD) UID() string { return t.Name }

func GetByID() string {
	return "SELECT * FROM tld WHERE id = $1"
}

func parseRow(row connection.Row) (TLD, error) {
	rt := TLD{}
	err := row.Scan(&rt.ID, &rt.Name)
	return rt, err

}

func Get(query string, args ...interface{}) (TLD, error) {
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]TLD, error) {
	conn, err := connection.Get()
	if err != nil {
		return []TLD{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []TLD{}, err
	}
	var list []TLD
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

package record_type

import (
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tld"
)

type RecordType struct {
	ID   int32
	Name string
}

var c = cache.New()

func New(name string, t tld.TLD) (RecordType, error) {
	if rt, ok := c.Check(name); ok {
		return rt.(RecordType), nil
	}
	conn, err := connection.Get()
	if err != nil {
		return RecordType{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT ensure_record_table($1, $2)", name, t.ID).Scan(&id)
	rt := RecordType{
		ID:   id,
		Name: name,
	}
	c.Add(rt)
	return rt, err
}

func (rt RecordType) UID() string { return rt.Name }

func GetByID() string {
	return "SELECT * FROM record_type WHERE id = $1"
}

func parseRow(row connection.Row) (RecordType, error) {
	rt := RecordType{}
	err := row.Scan(&rt.ID, &rt.Name)
	return rt, err
}

func Get(query string, args ...interface{}) (RecordType, error) {
	conn, err := connection.Get()
	if err != nil {
		return RecordType{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]RecordType, error) {
	conn, err := connection.Get()
	if err != nil {
		return []RecordType{}, err
	}
	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return []RecordType{}, err
	}
	var list []RecordType
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

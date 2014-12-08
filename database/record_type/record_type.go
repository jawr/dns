package record_type

import (
	"github.com/jawr/dns/database/connection"
)

type RecordType struct {
	ID   int32
	Name string
}

func New(name string) (RecordType, error) {
	conn, err := connection.Get()
	if err != nil {
		return RecordType{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_record_type($1)", name).Scan(&id)
	return RecordType{
		ID:   id,
		Name: name,
	}, err
}

func GetByID(id int32) (RecordType, error) {
	return Get("SELECT * FROM record_type WHERE id = $1", id)
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

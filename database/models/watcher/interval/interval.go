package interval

import (
	"github.com/jawr/dns/database/connection"
)

type Interval struct {
	ID    int32  `json:"id"`
	Value string `json:"value"`
}

func New(value string) (Interval, error) {
	i := Interval{}
	i.Value = value
	conn, err := connection.Get()
	if err != nil {
		return i, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_interval($1)", value).Scan(&id)
	i.ID = id
	return i, err
}

func (i Interval) UID() string { return i.Value }

const (
	SELECT string = "SELECT * FROM interval "
)

func GetAll() string {
	return SELECT
}

func GetByID() string {
	return SELECT + "WHERE id = $1"
}

func GetByValue() string {
	return SELECT + "WHERE value = $1"
}

func parseRow(row connection.Row) (Interval, error) {
	i := Interval{}
	err := row.Scan(&i.ID, &i.Value)
	return i, err
}

func Get(query string, args ...interface{}) (Interval, error) {
	var result Interval
	conn, err := connection.Get()
	if err != nil {
		return Interval{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err = parseRow(row)
	return result, err
}

func GetList(query string, args ...interface{}) ([]Interval, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Interval{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Interval{}, err
	}
	var list []Interval
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

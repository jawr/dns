package intervals

import (
	"github.com/jawr/dns/database/connection"
)

const (
	SELECT string = "SELECT * FROM interval "
)

type Result struct {
	One  func() (Interval, error)
	List func() ([]Interval, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Interval, error) {
			return Get(query, args...)
		},
		func() ([]Interval, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByID(id int32) Result {
	if interval, ok := byID[id]; ok {
		return Result{
			func() (Interval, error) {
				return interval, nil
			},
			func() ([]Interval, error) {
				return []Interval{interval}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE id = $1", id)
}

func GetByValue(value string) Result {
	if interval, ok := byValue[value]; ok {
		return Result{
			func() (Interval, error) {
				return interval, nil
			},
			func() ([]Interval, error) {
				return []Interval{interval}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE value = $1", value)
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

package tlds

import (
	"github.com/jawr/dns/database/connection"
)

const (
	SELECT string = "SELECT * FROM tld "
)

type Result struct {
	One  func() (TLD, error)
	List func() ([]TLD, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (TLD, error) {
			return Get(query, args...)
		},
		func() ([]TLD, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByID(id int32) Result {
	if tld, ok := byID[id]; ok {
		return Result{
			func() (TLD, error) {
				return tld, nil
			},
			func() ([]TLD, error) {
				return []TLD{tld}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE id = $1", id)
}

func GetByName(name string) Result {
	if tld, ok := byName[name]; ok {
		return Result{
			func() (TLD, error) {
				return tld, nil
			},
			func() ([]TLD, error) {
				return []TLD{tld}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE name = $1", name)
}

func parseRow(row connection.Row) (TLD, error) {
	t := TLD{}
	err := row.Scan(&t.ID, &t.Name)
	return t, err

}

func Get(query string, args ...interface{}) (TLD, error) {
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err := parseRow(row)
	return result, err
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

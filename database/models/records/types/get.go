package types

import (
	"github.com/jawr/dns/database/connection"
)

const (
	SELECT string = "SELECT * FROM record_type "
)

type Result struct {
	One  func() (Type, error)
	List func() ([]Type, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Type, error) {
			return Get(query, args...)
		},
		func() ([]Type, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByID(id int32) Result {
	if t, ok := byID[id]; ok {
		return Result{
			func() (Type, error) {
				return t, nil
			},
			func() ([]Type, error) {
				return []Type{t}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE id = $1", id)
}

func GetByName(name string) Result {
	if t, ok := byName[name]; ok {
		return Result{
			func() (Type, error) {
				return t, nil
			},
			func() ([]Type, error) {
				return []Type{t}, nil
			},
		}
	}
	return newResult(SELECT+"WHERE name = $1", name)
}

func parseRow(row connection.Row) (Type, error) {
	rt := Type{}
	err := row.Scan(&rt.ID, &rt.Name)
	return rt, err
}

func Get(query string, args ...interface{}) (Type, error) {
	conn, err := connection.Get()
	if err != nil {
		return Type{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Type, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Type{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Type{}, err
	}
	var list []Type
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

/*
func Search(params url.Values, idx, limit int) ([]Type, error) {
	query := GetAll()
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		case "name":
		case "id":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
			i++
		}
	}
	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + " "
	}
	query += fmt.Sprintf("LIMIT %d OFFSET %d", limit, idx)
	return GetList(query, args...)
}
*/

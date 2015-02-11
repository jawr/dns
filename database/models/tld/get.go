package tlds

import (
	"fmt"
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"net/url"
	"reflect"
	"strings"
)

const (
	SELECT string = "SELECT * FROM tld "
)

func GetAll() string {
	return SELECT
}

func GetByID() string {
	return SELECT + "WHERE id = $1"
}

func GetByName() string {
	return SELECT + "WHERE name = $1"
}

func parseRow(row connection.Row) (TLD, error) {
	t := TLD{}
	err := row.Scan(&t.ID, &t.Name)
	return t, err

}

func Get(query string, args ...interface{}) (TLD, error) {
	var result TLD
	if len(args) == 1 {
		switch reflect.TypeOf(args[0]).Kind() {
		case reflect.Int32:
			if i, ok := cacheGetID.Check(query, args[0].(int32)); ok {
				return i.(TLD), nil
			}
			defer func() {
				cacheGetID.Add(result, query, args[0].(int32))
			}()
		}
	}
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err = parseRow(row)
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

package domain

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tld"
	"net/url"
	"reflect"
	"strings"
)

var cacheGetByString = cache.NewCacheString()

const (
	SELECT string = "SELECT * FROM domain "
)

func GetByJoinWhoisEmails() string {
	return `
		SELECT DISTINCT d.* FROM domain AS d JOIN whois w ON d.uuid = w.domain
		WHERE w.emails ? $1`
}

func GetByNameAndTLD() string {
	return SELECT + "WHERE name = $1 AND tld = $2"
}

func GetByUUID() string {
	return SELECT + "WHERE uuid = $1"
}

func GetAll() string {
	return SELECT
}

func GetAllLimitOffset() string {
	return SELECT + "LIMIT $1 OFFSET $2"
}

func Search(params url.Values, offset, limit int) ([]Domain, error) {
	query := GetAll()
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		case "name", "uuid":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "tld":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			t, err := tld.Get(tld.GetByName(), params.Get(k))
			if err != nil {
				return []Domain{}, err
			}
			args = append(args, t.ID)
			i++
		}
	}
	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + " "
	}
	query += fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	return GetList(query, args...)
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var dUUID string
	var tldId int32
	err := row.Scan(&dUUID, &d.Name, &tldId)
	if err != nil {
		return d, err
	}
	fmt.Println(tldId)
	d.UUID = uuid.Parse(dUUID)
	t, err := tld.Get(tld.GetByID(), tldId)
	d.TLD = t
	fmt.Println(t)
	return d, err
}

func Get(query string, args ...interface{}) (Domain, error) {
	var result Domain
	if len(args) == 1 {
		switch reflect.TypeOf(args[0]).Kind() {
		case reflect.String:
			if i, ok := cacheGetByString.Check(query, args[0].(string)); ok {
				return i.(Domain), nil
			}
			defer func() {
				cacheGetByString.Add(result, query, args[0].(string))
			}()
		}
	}
	conn, err := connection.Get()
	if err != nil {
		return Domain{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err = parseRow(row)
	return result, err
}

func GetList(query string, args ...interface{}) ([]Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Domain{}, err
	}
	rows, err := conn.Query(query, args...)
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

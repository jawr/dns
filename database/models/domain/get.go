package domain

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tld"
	"net/url"
	"strings"
)

const (
	SELECT string = "SELECT * FROM domain "
)

func GetByNameAndTLD() string {
	return SELECT + "WHERE name = $1 AND tld = $2"
}

func GetByUUID() string {
	return SELECT + "WHERE uuid = $1"
}

func GetAll() string {
	return SELECT
}

func Search(params url.Values, idx, limit int) ([]Domain, error) {
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
	query += fmt.Sprintf("LIMIT %d OFFSET %d", limit, idx)
	fmt.Println(query)
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
	d.UUID = uuid.Parse(dUUID)
	t, err := tld.Get(tld.GetByID(), tldId)
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

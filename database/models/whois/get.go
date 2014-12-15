package whois

import (
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"net/url"
	"strings"
)

const (
	SELECT string = "SELECT * FROM whois "
)

func GetAll() string {
	return SELECT
}

func GetByID() string {
	return SELECT + "WHERE id = $1"
}

func Search(params url.Values, idx, limit int) ([]Result, error) {
	query := GetAll()
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		// TODO: handle times and json
		case "id":
		case "domain":
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

func parseRow(row connection.Row) (Result, error) {
	w := Result{}
	var uuid string
	err := row.Scan(&w.ID, &uuid, &w.Data, &w.Added)
	if err != nil {
		return w, err
	}
	d, err := domain.Get(domain.GetByUUID(), uuid)
	if err != nil {
		return w, err
	}
	w.Domain = d
	return w, nil
}

func Get(query string, args ...interface{}) (Result, error) {
	conn, err := connection.Get()
	if err != nil {
		return Result{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Result, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Result{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Result{}, err
	}
	var list []Result
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

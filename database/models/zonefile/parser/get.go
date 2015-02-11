package parser

import (
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tlds"
	"net/url"
	"strings"
)

const (
	SELECT string = "SELECT * FROM parser "
)

func GetAll() string {
	return SELECT
}

func GetByID() string {
	return SELECT + "WHERE id = $1"
}

func parseRow(row connection.Row) (Parser, error) {
	p := Parser{}
	var tldID int32
	err := row.Scan(&p.ID, &p.Filename, &p.Started, p.Finished, &p.Date, &tldID, &p.Logs)
	if err != nil {
		return p, err
	}
	t, err := tlds.GetByID(tldID).One()
	p.TLD = t
	return p, err
}

func Search(params url.Values, idx, limit int) ([]Parser, error) {
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

func Get(query string, args ...interface{}) (Parser, error) {
	conn, err := connection.Get()
	if err != nil {
		return Parser{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Parser, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Parser{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Parser{}, err
	}
	var list []Parser
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

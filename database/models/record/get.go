package record

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record_type"
	"net/url"
	"strings"
)

const (
	SELECT string = "SELECT * FROM record "
)

func GetAll() string {
	return SELECT
}

func GetByUUID() string {
	return SELECT + "WHERE uuid = $1"
}

func Search(params url.Values, idx, limit int) ([]Record, error) {
	query := GetAll()
	var where []string
	var args []interface{}
	i := 1
	for k, _ := range params {
		switch k {
		// TODO: handle times and args
		case "name":
			where = append(where, fmt.Sprintf("name = $%d", i))
			args = append(args, strings.ToLower(params.Get(k)))
			i++
		case "domain":
			where = append(where, fmt.Sprintf("domain = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "type":
			where = append(where, fmt.Sprintf("record_type = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "uuid":
			where = append(where, fmt.Sprintf("uuid = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "tld":
			where = append(where, fmt.Sprintf("tld = $%d", i))
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

func parseRow(row connection.Row) (Record, error) {
	r := Record{}
	var rUUID, dUUID string
	var rtID int32
	var args []byte
	err := row.Scan(&rUUID, &dUUID, &r.Name, &args, &rtID, &r.Date, &r.Added)
	if err != nil {
		return r, err
	}
	r.UUID = uuid.Parse(rUUID)
	d, err := domain.Get(domain.GetByUUID(), dUUID)
	if err != nil {
		return r, err
	}
	r.Domain = d
	rt, err := record_type.Get(record_type.GetByID(), rtID)
	if err != nil {
		return r, err
	}
	r.RecordType = rt
	return r, nil
}

func Get(query string, args ...interface{}) (Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return Record{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Record{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Record{}, err
	}
	var list []Record
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

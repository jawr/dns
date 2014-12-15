package record

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record_type"
	"github.com/jawr/dns/database/models/tld"
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
		// TODO: handle times and json
		case "name", "domain", "uuid", "record_type":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			args = append(args, params.Get(k))
			i++
		case "tld":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			t, err := tld.Get(tld.GetByName(), params.Get(k))
			if err != nil {
				return []Record{}, err
			}
			args = append(args, t.ID)
			i++
		case "type":
			where = append(where, fmt.Sprintf(k+" = $%d", i))
			rt, err := record_type.Get(record_type.GetByName(), params.Get(k))
			if err != nil {
				return []Record{}, err
			}
			args = append(args, rt.ID)
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
	err := row.Scan(&rUUID, &dUUID, &r.Name, &r.Args, &rtID, &r.Date, &r.Added)
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

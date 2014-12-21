package record_type

import (
	"fmt"
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tld"
	"net/url"
	"reflect"
	"strings"
)

type RecordType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

var c = cache.New()
var cacheGetID = cache.NewCacheInt32()

func New(name string, t tld.TLD) (RecordType, error) {
	if rt, ok := c.Check(name); ok {
		return rt.(RecordType), nil
	}
	conn, err := connection.Get()
	if err != nil {
		return RecordType{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT ensure_record_table($1, $2)", name, t.ID).Scan(&id)
	rt := RecordType{
		ID:   id,
		Name: name,
	}
	c.Add(rt)
	return rt, err
}

func (rt RecordType) UID() string { return rt.Name }

const (
	SELECT string = "SELECT * FROM record_type "
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

func parseRow(row connection.Row) (RecordType, error) {
	rt := RecordType{}
	err := row.Scan(&rt.ID, &rt.Name)
	return rt, err
}

func Get(query string, args ...interface{}) (RecordType, error) {
	var result RecordType
	if len(args) == 1 {
		switch reflect.TypeOf(args[0]).Kind() {
		case reflect.Int32:
			if i, ok := cacheGetID.Check(query, args[0].(int32)); ok {
				return i.(RecordType), nil
			}
			defer func() {
				cacheGetID.Add(result, query, args[0].(int32))
			}()
		}
	}
	conn, err := connection.Get()
	if err != nil {
		return RecordType{}, err
	}
	row := conn.QueryRow(query, args...)
	result, err = parseRow(row)
	return result, err
}

func GetList(query string, args ...interface{}) ([]RecordType, error) {
	conn, err := connection.Get()
	if err != nil {
		return []RecordType{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []RecordType{}, err
	}
	var list []RecordType
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

func Search(params url.Values, idx, limit int) ([]RecordType, error) {
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

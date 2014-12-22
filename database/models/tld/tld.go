package tld

import (
	"fmt"
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"net/url"
	"reflect"
	"strings"
)

type TLD struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

var c = cache.New()
var cacheGetID = cache.NewCacheInt32()

func New(name string) (TLD, error) {
	if t, ok := c.Check(name); ok {
		return t.(TLD), nil
	}
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_tld($1)", name).Scan(&id)
	t := TLD{
		ID:   id,
		Name: name,
	}
	c.Add(t)
	return t, err
}

func (t TLD) UID() string { return t.Name }

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

func Search(params url.Values, idx, limit int) ([]TLD, error) {
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

func Detect(s string) (TLD, error) {
	t, err := Get(GetByName(), s)
	if err != nil {
		args := strings.Split(s, ".")
		if len(args) > 1 {
			return Detect(strings.Join(args[1:], "."))
		}
		return TLD{}, err
	}
	return t, err
}

func DetectDomainAndTLD(s string) (string, TLD, error) {
	t, err := Detect(s)
	s = strings.TrimSuffix(s, "."+t.Name)
	args := strings.Split(s, ".")
	return args[len(args)-1], t, err
}

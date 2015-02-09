package watcher

import (
	"fmt"
	"net/url"
	"strings"
)

/* TODO: replace all instances of search with the safer version, think 300 */
func Search(params url.Values, idx, limit int) ([]Watcher, error) {
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

	fmt.Println(query)
	fmt.Println(args)
	return GetList(query, args...)
}

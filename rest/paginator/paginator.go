package paginator

import (
	"github.com/gorilla/context"
	"net/http"
	"net/url"
	"strconv"
)

func getInt(s string, i int) int {
	if len(s) > 0 {
		if j, err := strconv.ParseInt(s, 10, 32); err == nil {
			i = int(j)
		}
	}
	return i
}

func Paginate(fn func(http.ResponseWriter, *http.Request, url.Values, int, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if duuid, ok := context.GetOk(r, "duuid"); ok {
			params["duuid"] = []string{duuid.(string)}
		}
		if domain, ok := context.GetOk(r, "domain"); ok {
			params["domain"] = []string{domain.(string)}
		}
		limit := getInt(params.Get("limit"), 15)
		if limit > 50 {
			limit = 50
		}
		offset := limit * getInt(params.Get("page"), 0)
		fn(w, r, params, limit, offset)
	}
}

/* Extend paginate to build the query using the below pattern
package main

import "fmt"

func main() {
	var results []string
	var err error
	fn := func(query string, args ...interface{}) {
		results, err = GetList(query, args...)
	}
	switchFn := func(query *string) {
		*query += " SOME GENERATED STUFF"
	}
	buildQuery(fn, switchFn)
	fmt.Println(results)
	fmt.Println(err)
}

func buildQuery(fn func(string, ...interface{}), switchFn func(*string)) {
	query := "SELECT * FROM "
	switchFn(&query)
	fn(query)
}

func GetList(query string, args ...interface{}) ([]string, error) {
	return []string{query}, nil
}

*/

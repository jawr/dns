package paginator

import (
	"github.com/gorilla/context"
	"net/http"
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

func Paginate(fn func(http.ResponseWriter, *http.Request, map[string][]string, int, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if domain, ok := context.GetOk(r, "domain"); ok {
			params["domain"] = []string{domain.(string)}
		}
		limit := getInt(params.Get("limit"), 15)
		if limit > 50 {
			// log?
			limit = 50
		}
		idx := limit * getInt(params.Get("page"), 0)
		fn(w, r, params, idx, limit)
	}
}

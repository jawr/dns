package util

import (
	"encoding/json"
	"net/http"
)

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func parseError(err error, w http.ResponseWriter) {
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func ToJSON(v interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		parseError(err, w)
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		parseError(err, w)
		return
	}
	w.Write(b)
}

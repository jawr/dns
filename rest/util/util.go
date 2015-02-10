package util

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func parseError(err error, w http.ResponseWriter) {
	// TODO: switch error to decide what error code
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func Error(err error, w http.ResponseWriter) {
	parseError(err, w)
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
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		if reflect.ValueOf(v).Len() == 0 {
			b = []byte("[]")
		}
	}
	w.Write(b)
}

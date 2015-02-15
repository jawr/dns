package util

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
)

const (
	UUID_REGEX string = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
)

func parseError(err error, w http.ResponseWriter) {
	// TODO: switch error to decide what error code
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Resource not found."}`))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	}
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

package util

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func SetupRouter(router *mux.Router, name string, routes Routes) {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(name + route.Name).
			Handler(route.HandlerFunc)
	}
}

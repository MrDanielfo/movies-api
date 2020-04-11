package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Name(route.Name).Methods(route.Method).Path(route.Pattern).Handler(route.HandleFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"Movies",
		"GET",
		"/movies",
		MoviesHandler,
	},
	Route{
		"Movie Single",
		"GET",
		"/movies/{id}",
		MoviesSingle,
	},
	Route{
		"Movie Add",
		"POST",
		"/add-movie",
		MovieAdd,
	},
	Route{
		"Movie Update",
		"PUT",
		"/update/{id}",
		MovieUpdate,
	},
	Route{
		"Movie Delete",
		"DELETE",
		"/delete/{id}",
		MovieRemove,
	},
}

package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/img/").Handler(http.StripPrefix("/img", http.FileServer(http.Dir("./web/teamplate/img"))))
	return router
}

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",
		index,
	},

	Route{
		"Admin_panel",
		"GET",
		"/admin_panel",
		admin_panel,
	},

	Route{
		"Admin_panel_set",
		"POST",
		"/admin_panel",
		admin_panel_post,
	},

	Route{
		"Login",
		"GET",
		"/login",
		login,
	},
}

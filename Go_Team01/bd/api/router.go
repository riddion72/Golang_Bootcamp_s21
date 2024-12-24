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

	return router
}

var routes = append(internalRoutes, publicRoutes...)

var internalRoutes = Routes{
	Route{"InternalRequstBDset", "POST", "/internal_set", handlerInternalSET},
	Route{"InternalRequstBDdelle", "POST", "/internal_delete", handlerInternalDELETE},
}

var publicRoutes = Routes{
	Route{"RequstBDS", "POST", "/req_get", handlerGET},
	Route{"RequstBDP", "POST", "/req_set", handlerSET},
	Route{"RequstBDD", "POST", "/req_delete", handlerDELETE},
	Route{"HeartBeat", "GET", "/HeartBeat", handleHeartBeatGet},
	Route{"InitBeat", "POST", "/HeartBeat", handleHeartBeatPost},
	Route{"Migrate", "GET", "/Migrate", handleMigrate},
}

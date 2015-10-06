package route

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

func CreateRoutes() Routes {
	return Routes{
		Route{
			"URLRoot",
			"GET",
			"/",
			URLRoot,
		},
		Route{
			"URLShow",
			"GET",
			"/{shorturl}",
			URLShow,
		},
		Route{
			"URLCreate",
			"POST",
			"/create",
			URLCreate,
		},
	}
}

func NewLinkShortenerRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}

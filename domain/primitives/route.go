package domain_primitives

import "strings"

type Router struct {
	historyRoute []string
}

func NewRouter(
	userRoute string,
) Router {

	routes := strings.Split(userRoute, ".")

	return Router{
		historyRoute: routes,
	}
}

func (r *Router) CurrentRoute() string {
	return r.historyRoute[len(r.historyRoute)-1]
}

func (r *Router) PreviousRoute() Router {
	return Router{
		historyRoute: r.historyRoute[:len(r.historyRoute)-1],
	}
}

func (r *Router) NextRoute(route string) Router {
	return Router{
		historyRoute: append(r.historyRoute, route),
	}
}

package domain_primitives

import "strings"

type Router struct {
	redirect     bool
	historyRoute []string
}

func NewRouter(
	redirect bool,
	userRoute string,
) Router {

	routes := strings.Split(userRoute, ".")

	return Router{
		redirect:     redirect,
		historyRoute: routes,
	}
}

func (r *Router) IsRedirect() bool {
	return r.redirect
}

func (r *Router) HistoryRoute() string {
	return strings.Join(r.historyRoute, ".")
}

func (r *Router) CurrentRoute() string {
	return r.historyRoute[len(r.historyRoute)-1]
}

func (r *Router) PreviousRoute(redirect bool) Router {
	var prevRoute string
	if len(r.historyRoute) > 1 {
		prevRoute = r.historyRoute[len(r.historyRoute)-2]
	} else {
		prevRoute = r.historyRoute[0]
	}

	return Router{
		redirect:     redirect,
		historyRoute: append(r.historyRoute, prevRoute),
	}
}

func (r *Router) NextRoute(redirect bool, route string) Router {
	return Router{
		redirect:     redirect,
		historyRoute: append(r.historyRoute, route),
	}
}

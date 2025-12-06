package service

import (
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	"strings"
)

// RegisterRoute registers a route handler for the specified route name.
// When a message is received and the user's current route matches the
// registered route, the corresponding handler will be invoked.
// The route name is normalized to lowercase and trimmed of whitespace.
// Optional RouterHandlerOptions can be provided to customize handler behavior;
// if not provided, default options are used.
// Panics if a route with the same name is already registered.
func (app *ChatbotApp[Obs]) RegisterRoute(
	route string,
	handler d_router.RouteHandler[Obs],
	options ...d_router.RouterHandlerOptions,
) {

	route = strings.TrimSpace(route)
	route = strings.ToLower(route)

	if _, exists := app.routes[route]; exists {
		panic("route already registered: " + route)
	}

	opts := app.defaultOptions
	if len(options) > 0 {
		d := options[0]
		opts.SetOps(d)
	}

	app.routes[route] = d_router.RouterHandlerAdmnistrator[Obs]{
		HandlerOptions: opts,
		Handler:        handler,
	}
}

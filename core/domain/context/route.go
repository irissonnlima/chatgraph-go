package d_context

import d_route "chatgraph/core/domain/route"

func (c *ChatContext[Obs]) GetRoute() d_route.Route {
	return c.UserState.Route
}

func (c *ChatContext[Obs]) NextRoute(routeName string) d_route.Route {
	if c.Context.Err() != nil {
		return c.UserState.Route
	}
	return c.UserState.Route.Next(routeName)
}

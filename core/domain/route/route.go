package d_route

import "strings"

// Route represents the chatbot navigation history.
// It keeps track of the path taken by the user through the routes.
type Route struct {
	// History contains the list of visited routes in navigation order.
	History []string
	// Separator is the character used to split routes in the full path.
	Separator rune
}

// NewRoute creates a new Route instance from a full path string.
// The path is split by the specified separator and each segment is trimmed of spaces.
//
// Example:
//
//	route := NewRoute("start.menu.options", '.')
//	// route.History = ["start", "menu", "options"]
func NewRoute(fullPath string, separator rune) Route {
	segments := strings.Split(fullPath, string(separator))
	for i, segment := range segments {
		segments[i] = strings.TrimSpace(segment)
	}
	return Route{History: segments, Separator: separator}
}

func (r Route) IsRouteReturn() {}

// CurrentRepeated returns how many times the current route appears consecutively
// at the end of the history.
//
// Example:
//
//	route := NewRoute("start.choice.choice.choice", '.')
//	route.CurrentRepeated() // returns 3
func (r Route) CurrentRepeated() int {
	count := 0
	current := r.Current()
	for i := len(r.History) - 1; i >= 0; i-- {
		count++
		if r.History[i] != current {
			return count
		}
	}
	return count
}

// HistoryDedup returns the route history with consecutive duplicates removed.
// It preserves the order of appearance, useful for navigation without loops.
//
// Example:
//
//	route := NewRoute("start.choice.choice.select_a.choice", '.')
//	route.HistoryDedup() // returns ["start", "choice", "select_a", "choice"]
func (r Route) HistoryDedup() []string {
	if len(r.History) == 0 {
		return r.History
	}

	result := make([]string, 0, len(r.History))
	result = append(result, r.History[0])

	for i := 1; i < len(r.History); i++ {
		if r.History[i] != r.History[i-1] {
			result = append(result, r.History[i])
		}
	}

	return result
}

// Current returns the current route (last route in history).
// Returns an empty string if history is empty.
func (r Route) Current() string {
	if len(r.History) == 0 {
		return ""
	}
	return r.History[len(r.History)-1]
}

// Previous returns a new Route with deduplicated history and without the last route.
// Useful for navigating back without getting stuck in repeated route loops.
//
// Example:
//
//	route := NewRoute("start.choice.choice.menu", '.')
//	prev := route.Previous()
//	// prev.History = ["start", "choice"] (deduplicated and without "menu")
func (r Route) Previous() Route {
	historyDedup := r.HistoryDedup()

	return Route{
		History:   historyDedup[:len(historyDedup)-1],
		Separator: r.Separator,
	}
}

// Next adds a new route to the history and returns a new Route instance.
// The original Route remains unchanged (immutable operation).
//
// Example:
//
//	route := NewRoute("start", '.')
//	newRoute := route.Next("menu")
//	// newRoute.History = ["start", "menu"]
//	// route.History = ["start"] (unchanged)
func (r Route) Next(route string) Route {
	newHistory := make([]string, len(r.History), len(r.History)+1)
	copy(newHistory, r.History)
	newHistory = append(newHistory, route)

	return Route{
		History:   newHistory,
		Separator: r.Separator,
	}
}

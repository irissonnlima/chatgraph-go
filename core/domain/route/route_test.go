package d_route

import (
	"reflect"
	"testing"
)

func TestNewRoute(t *testing.T) {
	tests := []struct {
		name      string
		fullPath  string
		separator rune
		want      Route
	}{
		{
			name:      "simple path with arrow separator",
			fullPath:  "start > menu > options",
			separator: '>',
			want: Route{
				History:   []string{"start", "menu", "options"},
				Separator: '>',
			},
		},
		{
			name:      "path with pipe separator",
			fullPath:  "home|settings|profile",
			separator: '|',
			want: Route{
				History:   []string{"home", "settings", "profile"},
				Separator: '|',
			},
		},
		{
			name:      "single route",
			fullPath:  "start",
			separator: '>',
			want: Route{
				History:   []string{"start"},
				Separator: '>',
			},
		},
		{
			name:      "empty path",
			fullPath:  "",
			separator: '>',
			want: Route{
				History:   []string{""},
				Separator: '>',
			},
		},
		{
			name:      "path with extra spaces",
			fullPath:  "  start  >  menu  >  options  ",
			separator: '>',
			want: Route{
				History:   []string{"start", "menu", "options"},
				Separator: '>',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRoute(tt.fullPath, tt.separator)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_Current(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		want  string
	}{
		{
			name: "returns last route in history",
			route: Route{
				History:   []string{"start", "menu", "options"},
				Separator: '>',
			},
			want: "options",
		},
		{
			name: "single route",
			route: Route{
				History:   []string{"start"},
				Separator: '>',
			},
			want: "start",
		},
		{
			name: "empty history returns empty string",
			route: Route{
				History:   []string{},
				Separator: '>',
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.Current(); got != tt.want {
				t.Errorf("Route.Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_CurrentRepeated(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		want  int
	}{
		{
			name:  "three consecutive repeats at end",
			route: NewRoute("start > choice > choice > choice", '>'),
			want:  3, // choice appears 3 times consecutively at end
		},
		{
			name:  "no repeats at end",
			route: NewRoute("start > menu > options", '>'),
			want:  1, // options appears only once
		},
		{
			name:  "all same routes",
			route: NewRoute("choice > choice > choice > choice", '>'),
			want:  4,
		},
		{
			name:  "repeat in middle only",
			route: NewRoute("start > choice > choice > end", '>'),
			want:  1, // end appears only once at the end
		},
		{
			name:  "single route",
			route: NewRoute("start", '>'),
			want:  1,
		},
		{
			name:  "two consecutive repeats at end",
			route: NewRoute("start > menu > end > end", '>'),
			want:  2, // end appears 2 times consecutively at end
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.CurrentRepeated(); got != tt.want {
				t.Errorf("Route.CurrentRepeated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_HistoryDedup(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		want  []string
	}{
		{
			name:  "removes consecutive duplicates",
			route: NewRoute("start > choice > choice > select_a > choice", '>'),
			want:  []string{"start", "choice", "select_a", "choice"},
		},
		{
			name:  "no duplicates",
			route: NewRoute("start > menu > options", '>'),
			want:  []string{"start", "menu", "options"},
		},
		{
			name:  "all same routes",
			route: NewRoute("choice > choice > choice", '>'),
			want:  []string{"choice"},
		},
		{
			name:  "empty history",
			route: Route{History: []string{}, Separator: '>'},
			want:  []string{},
		},
		{
			name:  "single route",
			route: NewRoute("start", '>'),
			want:  []string{"start"},
		},
		{
			name:  "multiple groups of duplicates",
			route: NewRoute("a > a > b > b > b > c > c > a > a", '>'),
			want:  []string{"a", "b", "c", "a"},
		},
		{
			name:  "complex navigation pattern",
			route: NewRoute("start > choice > select_a > choice > choice > choice > select_b > end", '>'),
			want:  []string{"start", "choice", "select_a", "choice", "select_b", "end"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.HistoryDedup()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.HistoryDedup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_Previous(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		want  Route
	}{
		{
			name:  "goes back with deduplication",
			route: NewRoute("start > choice > choice > menu", '>'),
			want: Route{
				History:   []string{"start", "choice"},
				Separator: '>',
			},
		},
		{
			name:  "simple back navigation",
			route: NewRoute("start > menu > options", '>'),
			want: Route{
				History:   []string{"start", "menu"},
				Separator: '>',
			},
		},
		{
			name:  "back from many duplicates",
			route: NewRoute("start > choice > choice > choice > choice", '>'),
			want: Route{
				History:   []string{"start"},
				Separator: '>',
			},
		},
		{
			name:  "preserves separator",
			route: NewRoute("start|menu|options", '|'),
			want: Route{
				History:   []string{"start", "menu"},
				Separator: '|',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.Previous()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Previous() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_Next(t *testing.T) {
	tests := []struct {
		name      string
		route     Route
		nextRoute string
		want      Route
	}{
		{
			name:      "adds route to history",
			route:     NewRoute("start", '.'),
			nextRoute: "menu",
			want: Route{
				History:   []string{"start", "menu"},
				Separator: '.',
			},
		},
		{
			name:      "adds to existing history",
			route:     NewRoute("start.menu", '.'),
			nextRoute: "options",
			want: Route{
				History:   []string{"start", "menu", "options"},
				Separator: '.',
			},
		},
		{
			name:      "allows duplicate routes",
			route:     NewRoute("start.choice", '.'),
			nextRoute: "choice",
			want: Route{
				History:   []string{"start", "choice", "choice"},
				Separator: '.',
			},
		},
		{
			name:      "preserves separator",
			route:     NewRoute("start|menu", '|'),
			nextRoute: "options",
			want: Route{
				History:   []string{"start", "menu", "options"},
				Separator: '|',
			},
		},
		{
			name:      "adds empty route",
			route:     NewRoute("start", '.'),
			nextRoute: "",
			want: Route{
				History:   []string{"start", ""},
				Separator: '.',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.Next(tt.nextRoute)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_Next_Immutability(t *testing.T) {
	t.Run("original route remains unchanged", func(t *testing.T) {
		original := NewRoute("start.menu", '.')
		originalHistory := make([]string, len(original.History))
		copy(originalHistory, original.History)

		// Call Next
		newRoute := original.Next("options")

		// Verify original is unchanged
		if !reflect.DeepEqual(original.History, originalHistory) {
			t.Errorf("Original route was modified: got %v, want %v", original.History, originalHistory)
		}

		// Verify new route has correct history
		expectedNew := []string{"start", "menu", "options"}
		if !reflect.DeepEqual(newRoute.History, expectedNew) {
			t.Errorf("New route incorrect: got %v, want %v", newRoute.History, expectedNew)
		}
	})

	t.Run("multiple Next calls do not affect each other", func(t *testing.T) {
		base := NewRoute("start", '.')

		route1 := base.Next("menu")
		route2 := base.Next("settings")

		// Verify base is unchanged
		if len(base.History) != 1 || base.History[0] != "start" {
			t.Errorf("Base route was modified: %v", base.History)
		}

		// Verify route1
		expected1 := []string{"start", "menu"}
		if !reflect.DeepEqual(route1.History, expected1) {
			t.Errorf("route1 incorrect: got %v, want %v", route1.History, expected1)
		}

		// Verify route2
		expected2 := []string{"start", "settings"}
		if !reflect.DeepEqual(route2.History, expected2) {
			t.Errorf("route2 incorrect: got %v, want %v", route2.History, expected2)
		}
	})

	t.Run("chained Next calls work correctly", func(t *testing.T) {
		route := NewRoute("start", '.').
			Next("menu").
			Next("options").
			Next("confirm")

		expected := []string{"start", "menu", "options", "confirm"}
		if !reflect.DeepEqual(route.History, expected) {
			t.Errorf("Chained Next() incorrect: got %v, want %v", route.History, expected)
		}
	})
}

func TestRoute_NavigationFlow(t *testing.T) {
	t.Run("complete navigation scenario using Next", func(t *testing.T) {
		// Start navigation
		route := NewRoute("start", '.')

		// Verify initial state
		if route.Current() != "start" {
			t.Errorf("Initial current should be 'start', got %v", route.Current())
		}

		// Navigate through choice multiple times (simulating user returning to choice)
		route = route.Next("choice")
		route = route.Next("select_a")
		route = route.Next("choice")
		route = route.Next("choice")
		route = route.Next("choice")
		route = route.Next("select_b")
		route = route.Next("end")

		// Verify current
		if route.Current() != "end" {
			t.Errorf("Current should be 'end', got %v", route.Current())
		}

		// Verify deduplication
		expectedDedup := []string{"start", "choice", "select_a", "choice", "select_b", "end"}
		if !reflect.DeepEqual(route.HistoryDedup(), expectedDedup) {
			t.Errorf("HistoryDedup() = %v, want %v", route.HistoryDedup(), expectedDedup)
		}

		// Navigate back
		prev := route.Previous()
		if prev.Current() != "select_b" {
			t.Errorf("After Previous(), current should be 'select_b', got %v", prev.Current())
		}

		// Navigate back again
		prev = prev.Previous()
		if prev.Current() != "choice" {
			t.Errorf("After second Previous(), current should be 'choice', got %v", prev.Current())
		}
	})

	t.Run("forward and backward navigation", func(t *testing.T) {
		route := NewRoute("home", '.').
			Next("products").
			Next("details").
			Next("checkout")

		// Go back
		route = route.Previous()
		if route.Current() != "details" {
			t.Errorf("Expected 'details', got %v", route.Current())
		}

		// Go forward again
		route = route.Next("payment")
		expected := []string{"home", "products", "details", "payment"}
		if !reflect.DeepEqual(route.History, expected) {
			t.Errorf("History incorrect: got %v, want %v", route.History, expected)
		}
	})
}

func TestRoute_EdgeCases(t *testing.T) {
	t.Run("empty separator behavior", func(t *testing.T) {
		route := NewRoute("start", 0)
		if route.Current() != "start" {
			t.Errorf("Should handle zero separator, got %v", route.Current())
		}
	})

	t.Run("special characters in route names", func(t *testing.T) {
		route := NewRoute("start > menu-item > sub_option > option.1", '>')
		expected := []string{"start", "menu-item", "sub_option", "option.1"}
		if !reflect.DeepEqual(route.History, expected) {
			t.Errorf("Should handle special characters, got %v", route.History)
		}
	})

	t.Run("unicode route names", func(t *testing.T) {
		route := NewRoute("início > menu > opções", '>')
		expected := []string{"início", "menu", "opções"}
		if !reflect.DeepEqual(route.History, expected) {
			t.Errorf("Should handle unicode, got %v", route.History)
		}
	})
}

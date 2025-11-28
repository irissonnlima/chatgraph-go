// Package d_router provides routing configuration types for the chatbot framework.
// It includes route handlers, options, and trigger configurations.
package d_router

// RouteTrigger defines a pattern-based automatic route redirection.
// When a user's message matches the Regex pattern, the conversation
// is automatically redirected to the specified Route.
//
// This is useful for implementing global commands like "help", "cancel",
// or "menu" that should work regardless of the current route.
type RouteTrigger struct {
	// Regex is the regular expression pattern to match against user messages.
	// The pattern is evaluated using Go's regexp package.
	Regex string
	// Route is the target route name to redirect to when the pattern matches.
	Route string
}

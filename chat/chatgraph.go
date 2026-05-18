// Package chat provides a simple framework for building chatbot applications.
// It offers routing, message handling, context management, and integrations with
// message queues and external APIs.
//
// Example usage:
//
//	import "github.com/irissonnlima/chatgraph-go/chat"
//
//	app := chat.NewApp(receiver, router)
//	app.RegisterRoute("start", func(ctx *chat.Context[MyObs]) chat.RouteReturn {
//	    ctx.SendTextMessage("Hello!")
//	    return ctx.NextRoute("next")
//	})
//	app.Start()
package chat

import (
	"testing"

	input_queue "github.com/irissonnlima/chatgraph-go/adapters/input/queue"
	output_router_api "github.com/irissonnlima/chatgraph-go/adapters/output/router_api"
	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_context "github.com/irissonnlima/chatgraph-go/core/domain/context"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_guard "github.com/irissonnlima/chatgraph-go/core/domain/guard"
	d_logger "github.com/irissonnlima/chatgraph-go/core/domain/logger"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	adapter_input "github.com/irissonnlima/chatgraph-go/core/ports/adapters/input"
	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
	"github.com/irissonnlima/chatgraph-go/core/service"
)

// ============================================================================
// Type Aliases - Core Types
// ============================================================================

// Context is the chat context passed to route handlers.
// It provides access to user state, message, and utility methods.
type Context[Obs any] = d_context.ChatContext[Obs]

// RouteReturn is the return type for route handlers.
type RouteReturn = route_return.RouteReturn

// App is the main chatbot application.
type App[Obs any] = service.ChatbotApp[Obs]

// Engine handles route registration and execution logic without I/O.
// It is designed to be testable in isolation without requiring external dependencies.
type Engine[Obs any] = service.Engine[Obs]

// EngineTester is a test helper for validating chatbot handler executions.
type EngineTester[Obs any] = service.EngineTester[Obs]

// ExpectedAction represents an expected action during handler execution.
type ExpectedAction = service.ExpectedAction

// ActionExecType represents the type of action executed by a handler.
type ActionExecType = service.ActionExecType

// Action execution type constants for use with EngineTester.
const (
	ExecSendMessage    = service.ExecSendMessage
	ExecSetObservation = service.ExecSetObservation
	ExecSetRoute       = service.ExecSetRoute
	ExecGetFile        = service.ExecGetFile
	ExecUploadFile     = service.ExecUploadFile
)

// ============================================================================
// Type Aliases - Actions
// ============================================================================

// EndAction signals the end of a conversation session.
type EndAction = d_action.EndAction

// RedirectResponse triggers an immediate redirect to another route.
type RedirectResponse = d_action.RedirectResponse

// TransferToMenu transfers the user to a different menu.
type TransferToMenu = d_action.TransferToMenu

// ============================================================================
// Type Aliases - Message Types
// ============================================================================

// Message represents a chat message with text, buttons, and files.
type Message = d_message.Message

// TextMessage represents the text content of a message.
type TextMessage = d_message.TextMessage

// Button represents an interactive button in a message.
type Button = d_message.Button

// ButtonType represents the type of button (POSTBACK or URL).
type ButtonType = d_message.ButtonType

// Button type constants.
const (
	POSTBACK = d_message.POSTBACK
	URL      = d_message.URL
)

// ============================================================================
// Type Aliases - File Types
// ============================================================================

// File represents a file attachment.
type File = d_file.File

// FileType represents the type of file.
type FileType = d_file.FileType

// ============================================================================
// Type Aliases - User Types
// ============================================================================

// User represents user information.
type User = d_user.User

// UserState represents the complete state of a user's chat session.
type UserState[Obs any] = d_user.UserState[Obs]

// ChatID identifies a user and company for a chat.
type ChatID = d_user.ChatID

// Menu represents the current menu context.
type Menu = d_user.Menu

// AuthLevel represents the user's authorization level.
type AuthLevel = d_user.AuthLevel

// AuthLevel constants.
const (
	AuthBlocked = d_user.AuthBlocked
	AuthUnknown = d_user.AuthUnknown
	AuthRead    = d_user.AuthRead
	AuthWrite   = d_user.AuthWrite
)

// UserIdentity represents the authenticated identity of the user.
type UserIdentity = d_user.UserIdentity

// UserInternal represents internal HR data for the user.
type UserInternal = d_user.UserInternal

// ============================================================================
// Type Aliases - Guard Types
// ============================================================================

// GuardFunc is the function signature for route authorization guards.
type GuardFunc[Obs any] = d_guard.GuardFunc[Obs]

// ============================================================================
// Type Aliases - Logger Types
// ============================================================================

// UserLogger provides a structured logger scoped to a specific user.
type UserLogger = d_logger.UserLogger

// UserLoggerManager manages per-user loggers with file-based logging.
type UserLoggerManager = d_logger.UserLoggerManager

// LogLevel represents the severity of a log message.
type LogLevel = d_logger.LogLevel

// Log level constants.
const (
	LevelDebug   = d_logger.LevelDebug
	LevelInfo    = d_logger.LevelInfo
	LevelWarning = d_logger.LevelWarning
	LevelError   = d_logger.LevelError
)

// ============================================================================
// Type Aliases - Route Types
// ============================================================================

// Route represents the chatbot navigation history.
type Route = d_route.Route

// RouteHandler is the function signature for route handlers.
type RouteHandler[Obs any] = d_router.RouteHandler[Obs]

// RouterHandlerOptions configures behavior for route handlers.
type RouterHandlerOptions = d_router.RouterHandlerOptions

// TimeoutRouteOps configures timeout behavior for route handlers.
type TimeoutRouteOps = d_router.TimeoutRouteOps

// LoopCountRouteOps configures loop protection for route handlers.
type LoopCountRouteOps = d_router.LoopCountRouteOps

// ProtectedRouteOps configures route protection settings.
type ProtectedRouteOps = d_router.ProtectedRouteOps

// RouteTrigger defines a regex-based trigger for automatic route changes.
type RouteTrigger = d_router.RouteTrigger

// ============================================================================
// Type Aliases - Adapter Interfaces
// ============================================================================

// MessageReceiver is the interface for message queue consumers.
type MessageReceiver[Obs any] = adapter_input.IMessageReceiver[Obs]

// RouterService is the interface for routing and messaging operations.
type RouterService = adapter_output.IBotExecutor

// ============================================================================
// Constructors - Adapters
// ============================================================================

// NewRabbitMQ creates a new RabbitMQ message receiver.
func NewRabbitMQ[Obs any](user, password, host, vhost, queue string) MessageReceiver[Obs] {
	return input_queue.NewRabbitMQ[Obs](user, password, host, vhost, queue)
}

// NewRouterApi creates a new Router API service.
func NewRouterApi(url, username, password string) RouterService {
	return output_router_api.NewRouterApi(url, username, password)
}

// ============================================================================
// Constructors - Application
// ============================================================================

// NewApp creates a new chatbot application with the provided engine and adapters.
func NewApp[Obs any](
	engine *Engine[Obs],
	receiver MessageReceiver[Obs],
	router RouterService,
) *App[Obs] {
	return service.NewChatbotApp(engine, receiver, router)
}

// NewRoute creates a new Route from a path string.
func NewRoute(fullPath string, separator rune) Route {
	return d_route.NewRoute(fullPath, separator)
}

// NewEngine creates a new Engine instance for route registration and execution.
// The engine handles route registration and can be used both for testing and production.
func NewEngine[Obs any](options ...RouterHandlerOptions) *Engine[Obs] {
	return service.NewEngine[Obs](options...)
}

// NewEngineTester creates a new EngineTester for testing chatbot handlers.
// Use this to validate handler actions and return values.
func NewEngineTester[Obs any](t *testing.T, engine *Engine[Obs]) *EngineTester[Obs] {
	return service.NewEngineTester(t, engine)
}

// ============================================================================
// Constructors - Guard
// ============================================================================

// DefaultGuard returns a guard function that checks user.Identity.AuthLevel
// against the route's required level. Denied users are redirected to deniedRoute.
func DefaultGuard[Obs any](deniedRoute string) GuardFunc[Obs] {
	return d_guard.DefaultGuard[Obs](deniedRoute)
}

// InternalGuard returns a guard function that checks whether the user has
// internal HR data. Denied users are redirected to deniedRoute.
func InternalGuard[Obs any](deniedRoute string) GuardFunc[Obs] {
	return d_guard.InternalGuard[Obs](deniedRoute)
}

// ============================================================================
// Constructors - Logger
// ============================================================================

// NewUserLoggerManager creates a new per-user logger manager.
// Log files are stored in the specified directory.
func NewUserLoggerManager(dir string, level LogLevel) *UserLoggerManager {
	return d_logger.NewUserLoggerManager(dir, level)
}

// ParseLogLevel converts a string (e.g. "DEBUG", "INFO") to a LogLevel.
func ParseLogLevel(s string) LogLevel {
	return d_logger.ParseLogLevel(s)
}

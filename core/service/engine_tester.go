package service

import (
	"reflect"
	"testing"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// ActionType represents the type of action executed by the handler.
type ActionExecType int

const (
	ExecSendMessage ActionExecType = iota
	ExecSetObservation
	ExecSetRoute
	ExecGetFile
	ExecUploadFile
)

// ExpectedAction represents an expected action during execution.
type ExpectedAction struct {
	Type ActionExecType

	Message *d_message.Message

	Observation string

	Route string

	FileID   string
	FilePath string
	FileName string
}

// EngineTester is a test helper class for validating Engine executions.
type EngineTester[Obs any] struct {
	t      *testing.T
	engine *Engine[Obs]
}

// NewEngineTester creates a new EngineTester.
func NewEngineTester[Obs any](t *testing.T, engine *Engine[Obs]) *EngineTester[Obs] {
	return &EngineTester[Obs]{
		t:      t,
		engine: engine,
	}
}

// Execute runs the handler and validates expected actions and return.
func (e *EngineTester[Obs]) Execute(
	userState d_user.UserState[Obs],
	message d_message.Message,
	expectedActions []ExpectedAction,
	expectedReturn route_return.RouteReturn,
) {
	e.t.Helper()

	mock := newMockExecutor()
	result, err := e.engine.Execute(userState, message, mock)
	if err != nil {
		e.t.Fatalf("Execute returned error: %v", err)
	}

	// Validate actions
	e.validateActions(mock.expectedExec, expectedActions)

	// Validate return
	e.validateReturn(result, expectedReturn)
}

// validateActions validates if executed actions match expected ones.
func (e *EngineTester[Obs]) validateActions(actual, expected []ExpectedAction) {
	e.t.Helper()

	if len(actual) != len(expected) {
		e.t.Errorf("Expected %d actions, got %d", len(expected), len(actual))
		e.t.Logf("Expected: %+v", expected)
		e.t.Logf("Got: %+v", actual)
		return
	}

	for i, exp := range expected {
		act := actual[i]

		if act.Type != exp.Type {
			e.t.Errorf("Action %d: expected type %d, got %d", i, exp.Type, act.Type)
			continue
		}

		switch exp.Type {
		case ExecSendMessage:
			if exp.Message != nil && !reflect.DeepEqual(act.Message, exp.Message) {
				e.t.Errorf("Action %d: message mismatch\nExpected: %+v\nGot: %+v", i, exp.Message, act.Message)
			}
		case ExecSetObservation:
			if exp.Observation != "" && act.Observation != exp.Observation {
				e.t.Errorf("Action %d: expected observation %q, got %q", i, exp.Observation, act.Observation)
			}
		case ExecSetRoute:
			if exp.Route != "" && act.Route != exp.Route {
				e.t.Errorf("Action %d: expected route %q, got %q", i, exp.Route, act.Route)
			}
		case ExecGetFile:
			if exp.FileID != "" && act.FileID != exp.FileID {
				e.t.Errorf("Action %d: expected fileID %q, got %q", i, exp.FileID, act.FileID)
			}
		case ExecUploadFile:
			if exp.FilePath != "" && act.FilePath != exp.FilePath {
				e.t.Errorf("Action %d: expected filePath %q, got %q", i, exp.FilePath, act.FilePath)
			}
		}
	}
}

// validateReturn validates if the return matches the expected one.
func (e *EngineTester[Obs]) validateReturn(result, expected route_return.RouteReturn) {
	e.t.Helper()

	if reflect.TypeOf(result) != reflect.TypeOf(expected) {
		e.t.Errorf("Expected return type %T, got %T", expected, result)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		e.t.Errorf("Return value mismatch\nExpected: %+v\nGot: %+v", expected, result)
	}
}

package service

import (
	"fmt"
	"log"
)

// checkHealthRoutes validates the registered routes before starting the application.
// It checks:
// - A "start" route must exist
// - All routes referenced in triggers must exist
func (app *ChatbotApp[Obs]) checkHealthRoutes() error {
	mapRoutes := app.routes

	// Check if "start" route exists
	if _, exists := mapRoutes["start"]; !exists {
		return fmt.Errorf("required route 'start' is not registered")
	}

	// Check if all trigger routes exist
	for _, trigger := range app.routeTriggers {
		if _, exists := mapRoutes[trigger.Route]; !exists {
			return fmt.Errorf("trigger route '%s' (regex: %s) is not registered", trigger.Route, trigger.Regex)
		}
	}

	// Also check triggers defined in individual route options
	for routeName, handler := range mapRoutes {
		for _, trigger := range handler.HandlerOptions.Triggers {
			if _, exists := mapRoutes[trigger.Route]; !exists {
				return fmt.Errorf("trigger route '%s' in route '%s' (regex: %s) is not registered",
					trigger.Route, routeName, trigger.Regex)
			}
		}
	}

	// Also check the default options triggers
	rhoRoutes := app.defaultOptions.GetRhoRoutes()
	for _, rhoRoute := range rhoRoutes {
		if _, exists := mapRoutes[rhoRoute]; !exists {
			return fmt.Errorf("default option route '%s' is not registered", rhoRoute)
		}
	}

	log.Printf("[INFO] Routes validated successfully. %d routes registered.", len(mapRoutes))
	return nil
}

// Start begins consuming messages from the message receiver in an infinite loop.
// It only returns an error in case of a critical failure.
// Non-critical errors from HandleMessage are logged but do not stop the consumer.
func (app *ChatbotApp[Obs]) Start() error {

	if err := app.checkHealthRoutes(); err != nil {
		log.Printf("[ERROR] Failed to setup routes: %v", err)
		return err
	}

	messages := app.messageReceiver.ConsumeMessage()

	for msg := range messages {
		err := app.HandleMessage(msg.UserState, msg.Message)
		if err != nil {
			log.Printf("[ERROR] Failed to handle message: %v", err)
		}
	}

	log.Println("[CRITICAL] Message channel closed unexpectedly")
	return nil
}

package controller

import (
	"net/http"
)

// ControllerBase provides common functionalities for all controllers
type ControllerBase struct{}

// NewControllerBase creates a new ControllerBase instance
func NewControllerBase() *ControllerBase {
	return &ControllerBase{}
}

// Execute allows controllers to execute an action and potentially include common functionalities.
func (c *ControllerBase) Execute(action func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Before action hook, for logging, authentication, etc.
		action(w, r)
		// After action hook, for cleanup, logging, etc.
	}
}

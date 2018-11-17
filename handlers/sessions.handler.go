// Package handlers is a list of handler functions.
package handlers

import (
	"net/http"

	sessions "../models/sessions"
)

// CleanSessions clears expired sessions from the database
func CleanSessions(res http.ResponseWriter, req *http.Request) {
	sessions.CleanSessions()
}

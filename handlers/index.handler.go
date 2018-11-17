// Package handlers is a list of handler functions.
package handlers

import (
	"fmt"
	"net/http"

	globals "../config/globals"
	sessions "../models/sessions"
	users "../models/users"
)

// Index displays the index page
func Index(res http.ResponseWriter, req *http.Request) {
	// User is logged in
	if sessions.AlreadyLoggedIn(res, req) {
		u := users.GetUser(res, req)

		err := globals.Tpl.ExecuteTemplate(res, "index.gohtml", u)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err := globals.Tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

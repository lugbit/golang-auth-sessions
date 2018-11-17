// Package handlers is a list of handler functions.
package handlers

import (
	"fmt"
	"net/http"
	"time"

	globals "../config/globals"
	sessions "../models/sessions"
	users "../models/users"
	validators "../models/validators"
)

// Login handler authenticates the user login request
func Login(res http.ResponseWriter, req *http.Request) {
	// If the user is already logged in, redirect to homepage
	if sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Form was submitted
	if req.Method == http.MethodPost {
		// Validate form fields
		login := &validators.LoginValidator{
			Email:    req.FormValue("frmEmail"),
			Password: req.FormValue("frmPassword"),
		}

		// If validation fails, show errors in template
		if login.Validate() == false {
			err := globals.Tpl.ExecuteTemplate(res, "login.gohtml", login)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Validation successful
		// Get user by email
		u := users.GetUserByEmail(login.Email)
		// Create a new cookie and sent to client (GenerateCookie generates a UUID)
		c := sessions.GenerateCookie()
		http.SetCookie(res, c)
		// Create a session for the user with the same UUID as the cookie
		sessions.CreateSession(c.Value, u.UserID)
		// Redirect to the profile page
		http.Redirect(res, req, "/my-profile", http.StatusSeeOther)

		return
	}

	err := globals.Tpl.ExecuteTemplate(res, "login.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// Logout handler logs the user out
func Logout(res http.ResponseWriter, req *http.Request) {
	if !sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	//Retrieve Session cookie from the request
	c, _ := req.Cookie("session")
	// Delete the session(s) for this user in the database
	sessions.DeleteSessions(users.SessionUserID(c.Value))
	// Create a new session cookie and set to expire immeidately
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1, // Set to delete straight away
	}
	// Send cookie which is set to expire immediately
	http.SetCookie(res, c)

	// Clean up sessions every time a user logs out only if the last time
	// the sessions was cleaned is greater than 30 minutes
	if time.Now().Sub(sessions.DBSessionsCleaned) > (time.Second * 1800) {
		go sessions.CleanSessions()
	}

	// Redirect to login page
	http.Redirect(res, req, "/login", http.StatusSeeOther)
}

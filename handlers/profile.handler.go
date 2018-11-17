// Package handlers is a list of handler functions.
package handlers

import (
	"fmt"
	"log"
	"net/http"

	globals "../config/globals"
	sessions "../models/sessions"
	users "../models/users"
	validators "../models/validators"
	"golang.org/x/crypto/bcrypt"
)

// MyProfile handler displays the user's profile
func MyProfile(res http.ResponseWriter, req *http.Request) {
	// If the user is already logged in, redirect to homepage
	if !sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/login", http.StatusForbidden)
		return
	}

	// Get the user based on the session cookie stored on their browser
	u := users.GetUser(res, req)
	// Execute the profile page and pass the user
	err := globals.Tpl.ExecuteTemplate(res, "my-profile.gohtml", u)
	if err != nil {
		log.Fatalln(err)
	}
}

// UpdatePassword updates logged in user's password
func UpdatePassword(res http.ResponseWriter, req *http.Request) {
	if !sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/login", http.StatusForbidden)
		return
	}

	if req.Method == http.MethodPost {
		updatePW := &validators.UpdatePasswordValidator{
			CurrentPassword:    req.FormValue("frmCurrentPassword"),
			NewPassword:        req.FormValue("frmNewPassword"),
			ConfirmNewPassword: req.FormValue("frmConfirmNewPassword"),
		}

		// Get the user based on the session cookie stored on their browser
		u := users.GetUser(res, req)
		if updatePW.Validate(u.Password) == false {
			err := globals.Tpl.ExecuteTemplate(res, "update-password.gohtml", updatePW)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Proceed to update password
		// Hash the new passwords
		hashedNewPW, err := bcrypt.GenerateFromPassword([]byte(updatePW.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		}

		// Update the user's password in the database
		users.UpdatePassword(string(hashedNewPW), u.UserID)

		// Log user out
		http.Redirect(res, req, "/logout", http.StatusSeeOther)

		return
	}

	err := globals.Tpl.ExecuteTemplate(res, "update-password.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// Package handlers is a list of handler functions.
package handlers

import (
	"fmt"
	"log"
	"net/http"

	globals "../config/globals"
	smtp "../config/smtp"
	sessions "../models/sessions"
	tokens "../models/tokens"
	users "../models/users"
	validators "../models/validators"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register handler handles user registration
func Register(res http.ResponseWriter, req *http.Request) {
	// If the user is already logged in, redirect to homepage
	if sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Form was submitted
	if req.Method == http.MethodPost {
		// Validate form fields
		rv := &validators.RegisterValidator{
			FirstName:     req.FormValue("frmFirstName"),
			LastName:      req.FormValue("frmLastName"),
			Email:         req.FormValue("frmEmail"),
			Password:      req.FormValue("frmPassword"),
			PasswordAgain: req.FormValue("frmPasswordAgain"),
		}

		// If validation fails, show errors in template
		if rv.Validate() == false {
			err := globals.Tpl.ExecuteTemplate(res, "register.gohtml", rv)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Validation successful
		// Hash password using default cost
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(rv.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		}

		// Collect form fields and create a user struct
		u := users.User{
			FirstName: rv.FirstName,
			LastName:  rv.LastName,
			Email:     rv.Email,
			Password:  hashedPass,
		}
		// Generate new UUID token for activation token
		token, err := uuid.NewV4()
		if err != nil {
			log.Fatalln(err)
		}
		// Insert the user and the activated token into the database
		users.InsertUser(u, token.String())
		// Send verification link to users email
		smtp.SendEmail(token.String(), rv.Email)
		// Redirect user using status see other

		err = globals.Tpl.ExecuteTemplate(res, "limbo.gohtml", nil)
		if err != nil {
			fmt.Println(err)
		}
		// Redirect user to the account activation reminder page
		http.Redirect(res, req, "/activation-reminder", http.StatusSeeOther)

		return
	}

	err := globals.Tpl.ExecuteTemplate(res, "register.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// ActivateAccount handler is the endpoint for when a user tries to activate their account
func ActivateAccount(res http.ResponseWriter, req *http.Request) {
	if sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Grab the token parameter value from URL
	t := req.FormValue("token")
	// Check if token is already used. If already used, display link expired page
	if tokens.TokenAlreadyUsed(t) {
		err := globals.Tpl.ExecuteTemplate(res, "token-expired.gohtml", nil)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Check if the token is expired
	if tokens.TokenExpired(t) {
		err := globals.Tpl.ExecuteTemplate(res, "token-expired.gohtml", nil)
		if err != nil {
			fmt.Println(err)
		}
		// Mark token as used even if token is expired
		tokens.MarkTokenUsed(t)
		return
	}

	// Token is not used or expired, proceed to activate the user
	// Mark user as activated
	users.ActivateUser(t)
	// Mark token as used
	tokens.MarkTokenUsed(t)
	// Redirect to the login page
	http.Redirect(res, req, "/login", http.StatusSeeOther)
}

// SendActivationLink handler generates and sends new activation links
func SendActivationLink(res http.ResponseWriter, req *http.Request) {
	if sessions.AlreadyLoggedIn(res, req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		sav := &validators.SendActivationValidator{
			Email: req.FormValue("frmEmail"),
		}
		// Validate form
		if sav.Validate() == false {
			err := globals.Tpl.ExecuteTemplate(res, "send-activation.gohtml", sav)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Send new link
		// Generate new UUID token for new link
		token, err := uuid.NewV4()
		if err != nil {
			log.Fatalln(err)
		}

		// Update users token in the database with the new token
		tokens.UpdateToken(token.String(), sav.Email)
		// Send verification email to user with new token embedded in link
		smtp.SendEmail(token.String(), sav.Email)
		// Redirect user to the account activation reminder page
		http.Redirect(res, req, "/activation-reminder", http.StatusSeeOther)

		return
	}

	err := globals.Tpl.ExecuteTemplate(res, "send-activation.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// ActivationMessage displays the template reminding the user to activate their account
func ActivationMessage(res http.ResponseWriter, req *http.Request) {
	err := globals.Tpl.ExecuteTemplate(res, "limbo.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// Package validators is used to validate forms
package validators

import (
	"regexp"
	"strings"

	users "../users"
	"golang.org/x/crypto/bcrypt"
)

// RegisterValidator struct to store fields from the register page
type RegisterValidator struct {
	FirstName     string
	LastName      string
	Email         string
	Password      string
	PasswordAgain string
	Errors        map[string]string
}

// Validate function validates each fields and returns true if there are no errors
func (rv *RegisterValidator) Validate() bool {
	// Errors map to store the errors collected if any
	rv.Errors = make(map[string]string)

	// Check if the first name field is empty. If empty, add top errors map.
	if strings.TrimSpace(rv.FirstName) == "" {
		rv.Errors["FirstName"] = "First name cannot be empty"
	}
	if strings.TrimSpace(rv.LastName) == "" {
		rv.Errors["LastName"] = "Last name cannot be empty"
	}
	// Check if email is not empty, check it's the correct format and if it already exists.
	if strings.TrimSpace(rv.Email) != "" {
		re := regexp.MustCompile(".+@.+\\..+")
		matched := re.Match([]byte(rv.Email))

		if matched == false {
			rv.Errors["EmailFormat"] = "Email is not valid"
		}
		if users.EmailExists(rv.Email) {
			rv.Errors["EmailExists"] = "That email already exists"
		}
	} else {
		rv.Errors["Email"] = "Email cannot be empty"
	}
	// Check if the password field is not empty and if it matches the confirm password field
	if strings.TrimSpace(rv.Password) != "" {
		if strings.TrimSpace(rv.Password) != strings.TrimSpace(rv.PasswordAgain) {
			rv.Errors["PasswordMatch"] = "Passwords do not match"
		}
	} else {
		rv.Errors["Password"] = "Password cannot be empty"
	}

	// Returns true if there's no errrors recorded.
	return len(rv.Errors) == 0
}

// SendActivationValidator struct to store fields from the send activation page
type SendActivationValidator struct {
	Email  string
	Errors map[string]string
}

// Validate function validates each fields and returns true if there are no errors
func (sav *SendActivationValidator) Validate() bool {
	sav.Errors = make(map[string]string)

	if strings.TrimSpace(sav.Email) != "" {
		re := regexp.MustCompile(".+@.+\\..+")
		matched := re.Match([]byte(sav.Email))

		if matched == false {
			sav.Errors["EmailFormat"] = "Email is not valid"
		}
		// Check if email doesn't already exist
		if !users.EmailExists(sav.Email) {
			sav.Errors["EmailNotFound"] = "That email doesn't exists"
		}
		// Check if user is already activated
		if users.UserActivated(sav.Email) {
			sav.Errors["AccountAlreadyActivated"] = "That account is already activated. No link sent."
		}
	} else {
		sav.Errors["Email"] = "Email cannot be empty"
	}

	return len(sav.Errors) == 0
}

// LoginValidator struct to store fields from the login page
type LoginValidator struct {
	Email    string
	Password string
	Errors   map[string]string
}

// Validate function validates each fields and returns true if there are no errors
func (lv *LoginValidator) Validate() bool {
	lv.Errors = make(map[string]string)

	if strings.TrimSpace(lv.Email) != "" {
		// Check if email exists
		if !users.EmailExists(lv.Email) {
			lv.Errors["EmailDoesNotExist"] = "Email does not exist"
		} else {
			// Email exists
			// Check that the user account is activated
			if !users.UserActivated(lv.Email) {
				lv.Errors["NotActive"] = "Account not activated. Please activate your account first."
			} else {
				// Get user info by username
				u := users.GetUserByEmail(lv.Email)
				// Compare the user's hashed password from the database with the provided password
				err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(lv.Password))
				// Password doesn't match the hash
				if err != nil {
					lv.Errors["Email"] = "Email or Password is incorrect"
				}
			}
		}
	} else {
		lv.Errors["Email"] = "Email cannot be empty"
	}
	if strings.TrimSpace(lv.Password) == "" {
		lv.Errors["Password"] = "Password cannot be empty"
	}

	return len(lv.Errors) == 0
}

// UpdatePasswordValidator struct to store fields from update password page
type UpdatePasswordValidator struct {
	CurrentPassword    string
	NewPassword        string
	ConfirmNewPassword string
	Errors             map[string]string
}

// Validate function validates each fields and returns true if there are no errors
func (up *UpdatePasswordValidator) Validate(userPW []byte) bool {
	up.Errors = make(map[string]string)

	if strings.TrimSpace(up.CurrentPassword) != "" {
		// Compare the user's hashed password from the database with the provided password
		err := bcrypt.CompareHashAndPassword(userPW, []byte(up.CurrentPassword))
		// Password doesn't match the hash
		if err != nil {
			up.Errors["InvalidCurrentPassword"] = "Current password is incorrect"
		}
	} else {
		up.Errors["CurrentPassword"] = "Current password is required"
	}
	if strings.TrimSpace(up.NewPassword) != "" {
		if strings.TrimSpace(up.NewPassword) != strings.TrimSpace(up.ConfirmNewPassword) {
			up.Errors["NewPasswordMatch"] = "New passwords do not match"
		}
	} else {
		up.Errors["NewPasswordEmpty"] = "New password is required"
	}

	return len(up.Errors) == 0
}

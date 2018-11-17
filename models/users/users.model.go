// Package users is used for user related functions
package users

import (
	"database/sql"
	"log"
	"net/http"

	globals "../../config/globals"
	sessions "../sessions"
)

// User struct
type User struct {
	UserID    int
	FirstName string
	LastName  string
	Email     string
	Password  []byte
}

// GetUser returns a user based on the cookie UUID
func GetUser(res http.ResponseWriter, req *http.Request) User {
	// Request for the 'session' cookie
	c, err := req.Cookie("session")
	if err != nil {
		// Return empty user
		return User{}
	}
	// Set cookie path to index
	c.Path = "/"
	// Reset cookie age to the session length everytime this function runs
	c.MaxAge = sessions.SessionLength
	http.SetCookie(res, c)

	// If the cookie exists, get user from the database
	var u User
	sqlStmt := `SELECT
					tblUsers.fldID,
    				tblUsers.fldFirstName,
    				tblUsers.fldLastName,
    				tblUsers.fldEmail,
    				tblUsers.fldPassword
				FROM tblUsers
				INNER JOIN tblSessions
					ON tblSessions.fldFKUserID = tblUsers.fldID
				WHERE tblSessions.fldSessionID = ?`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.Value).Scan(&u.UserID, &u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return u
		}
		// Something else went wrong
		log.Fatal(err)
	}
	// Update the users last active
	sessions.UpdateLastActive(c.Value)

	return u
}

// EmailExists returns true if the email address argument exists in the database
func EmailExists(e string) bool {
	sqlStmt := `SELECT 
					fldEmail 
				FROM tblUsers 
				WHERE fldEmail = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var email string
	err = stmt.QueryRow(e).Scan(&email)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return false
		}
		// Something else went wrong
		log.Fatal(err)
	}

	return true
}

// InsertUser inserts argument user and activation token in the database
func InsertUser(u User, t string) {
	// Begin transaction
	tx, err := globals.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// Prepare user insertion and execute
	stmt, err := tx.Prepare("INSERT INTO tblUsers(fldFirstName, fldLastName, fldEmail, fldPassword) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		log.Fatal(err)
	}

	// Get the last inserted row's ID
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare token insertion and execute
	stmt, err = tx.Prepare("INSERT INTO tblActivationTokens (fldToken,fldFKUserID) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err = stmt.Exec(t, lastID)
	if err != nil {
		log.Fatal(err)
	}

	// Commit query
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// GetUserByEmail returns a user based on the argument email
func GetUserByEmail(e string) User {
	sqlStmt := `SELECT 
					fldID,
					fldFirstName,
    				fldLastName,
					fldEmail,
					fldPassword
				FROM tblUsers
				WHERE fldEmail = ?;`

	// Prepare statement
	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var u User
	err = stmt.QueryRow(e).Scan(&u.UserID, &u.FirstName, &u.LastName, &u.LastName, &u.Password)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			// Return empty user
			return u
		}
		// Something else went wrong
		log.Fatal(err)
	}

	return u
}

// ActivateUser marks the user activated in the database based on the activation token argument
func ActivateUser(t string) {
	sqlStmt := `UPDATE tblActivationTokens
				SET fldIsActivated = true
				WHERE fldToken = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// Execute insert and replace with actual values.
	_, err = stmt.Exec(t)
	if err != nil {
		log.Fatalln(err)
	}
}

// UserActivated returns true if a user is already activated based on the email argument
func UserActivated(e string) bool {
	sqlStmt := `SELECT 
					fldIsActivated
				FROM tblUsers
				INNER JOIN tblActivationTokens
				ON tblActivationTokens.fldFKUserID = tblUsers.fldID
				WHERE fldEmail = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var isActivated int
	err = stmt.QueryRow(e).Scan(&isActivated)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return false
		}
		// Something else went wrong
		log.Fatal(err)
	}
	// User is not activated
	if isActivated == 0 {
		return false
	}

	return true
}

// SessionUserID returns the user's ID based on the session ID argument
func SessionUserID(sID string) int {
	sqlStmt := `SELECT
					fldFKUserID
				FROM tblSessions
				WHERE fldSessionID = ?`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var userID int
	err = stmt.QueryRow(sID).Scan(&userID)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			log.Fatal(err)
		}
		// Something else went wrong
		log.Fatal(err)
	}

	return userID
}

// UpdatePassword updates the users password based on the user's id
func UpdatePassword(newPW string, userID int) {
	sqlStmt := `UPDATE tblUsers
				SET fldPassword = ?
				WHERE fldID = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// Execute insert and replace with actual values.
	_, err = stmt.Exec(newPW, userID)
	if err != nil {
		log.Fatalln(err)
	}
}

// Package sessions is for session related functions
package sessions

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	globals "../../config/globals"
	uuid "github.com/satori/go.uuid"
)

// Session struct
type session struct {
	SID        string
	CreatedAt  time.Time
	LastActive time.Time
	UserID     int
}

// DBSessionsCleaned keeps track of the last time sessions was cleaned
var DBSessionsCleaned = time.Now()

// SessionLength is the length of the max life of a cookie in seconds
// 3600s = 1hr
const SessionLength int = 3600

// AlreadyLoggedIn returns true if the requesting user already has a session
func AlreadyLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	var s session
	// Attempts to retrieve the 'session' cookie from the request
	c, err := req.Cookie("session")
	// No session cookie found
	if err != nil {
		return false
	}

	// Session cookie exists, query the database for the session details
	sqlStmt := `SELECT 
					fldSessionID,
					fldCreatedAt,
					fldLastActive,
					fldFKUserID 
				FROM tblSessions 
				WHERE fldSessionID = ?`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// Execute query and place result in container struct s
	err = stmt.QueryRow(c.Value).Scan(&s.SID, &s.CreatedAt, &s.LastActive, &s.UserID)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return false
		}
		// Something else went wrong
		log.Fatal(err)
	}

	// User has a session, update last active to now
	UpdateLastActive(c.Value)
	// Refresh session cookie
	c.Path = "/"
	c.MaxAge = int(SessionLength)
	http.SetCookie(res, c)

	return true
}

// UpdateLastActive is a helper function to update the users last active field
func UpdateLastActive(sID string) {
	sqlStmt := `UPDATE tblSessions
				SET fldLastActive = CURRENT_TIMESTAMP()
				WHERE fldSessionID = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sID)
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateSession creates a new session entry
func CreateSession(sID string, userID int) {
	sqlStmt := `INSERT INTO tblSessions(fldSessionID, fldFKUserID) 
				VALUES (?, ?)`
	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sID, userID)
	if err != nil {
		log.Fatalln(err)
	}
}

// DeleteSessions deletes all session for a user
func DeleteSessions(userID int) {
	sqlStmt := `DELETE FROM tblSessions 
				WHERE fldFKUserID = ?;`
	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		log.Fatalln(err)
	}
}

// GenerateCookie generates cookie with a UUID
func GenerateCookie() *http.Cookie {
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}

	return c
}

// CleanSessions clears all sessions from the database that are older than the allowed session length
func CleanSessions() {
	sqlStmt := `SELECT 
					fldSessionID,
					fldCreatedAt,
					fldLastActive,
					fldFKUserID 
				FROM tblSessions;`

	// Prepare statment
	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// rows will hold all found rows
	rows, err := stmt.Query()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	// Loop through each row
	for rows.Next() {
		// Scan each field into a struct field of s
		var s session
		err := rows.Scan(&s.SID, &s.CreatedAt, &s.LastActive, &s.UserID)
		if err != nil {
			log.Fatalln(err)
		}

		//Check if session is expired
		if time.Now().Sub(s.LastActive) > (time.Second * 3600) {
			// Delete expired session from the database
			clearSession(s.SID)
		}
	}

	// Check for errors during row scan
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}

	// Set DBSessionsCleaned last cleaned to now
	DBSessionsCleaned = time.Now()
}

// clearSession is a helper function that clears a session from the database
func clearSession(sID string) {
	sqlStmt := `DELETE FROM tblSessions
				WHERE fldSessionID = ?`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sID)
	if err != nil {
		log.Fatalln(err)
	}
}

// Package tokens is used for activation token related function
package tokens

import (
	"database/sql"
	"log"

	globals "../../config/globals"
)

const (
	// Set the token expiration time in minutes
	// 1440m = 24hr
	tExpMinutes = 1440
)

// TokenExpired returns true if argument token is expired
func TokenExpired(t string) bool {
	sqlStmt := `SELECT 
					TIMESTAMPDIFF(MINUTE, fldTokenDateIssued, NOW()) AS tokenActiveInMinutes
				FROM tblActivationTokens
				WHERE fldToken = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Container to hold the token life in minutes
	var tActiveMinutes int
	err = stmt.QueryRow(t).Scan(&tActiveMinutes)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return true
		}
		log.Fatal(err)
	}

	// If the token is less than the set token expiry date, token is not expired.
	if tActiveMinutes < tExpMinutes {
		return false
	}

	// The token is expired
	return true
}

// MarkTokenUsed marks the argument token as used in the database
func MarkTokenUsed(t string) {
	sqlStmt := `UPDATE tblUsers
				INNER JOIN tblActivationTokens
					ON tblActivationTokens.fldFKUserID = tblUsers.fldID
				SET fldTokenDateUsed = NOW()
				WHERE fldToken = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// Execute insert and replace with actual value(s).
	_, err = stmt.Exec(t)
	if err != nil {
		log.Fatalln(err)
	}
}

// TokenAlreadyUsed returns true if the argument token has already been used
func TokenAlreadyUsed(t string) bool {
	sqlStmt := `SELECT 
					fldTokenDateUsed
				FROM tblActivationTokens
				WHERE fldToken = ?`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	var dateUsed sql.NullString
	err = stmt.QueryRow(t).Scan(&dateUsed)
	if err != nil {
		// No rows found
		if err == sql.ErrNoRows {
			return true
		}
		log.Fatalln(err)
	}

	// If fldTokenDateused is NULL, token has not been used
	if !dateUsed.Valid {
		return false
	}

	// fldTokenDateused is not empty, token has been used
	return true
}

// UpdateToken updates the token in the database
func UpdateToken(t string, e string) {
	sqlStmt := `UPDATE tblUsers
				INNER JOIN tblActivationTokens
					ON tblActivationTokens.fldFKUserID = tblUsers.fldID
				SET fldToken = ?, fldTokenDateIssued = NOW(), fldTokenDateUsed = NULL
				WHERE fldEmail = ?;`

	stmt, err := globals.Db.Prepare(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	// Execute insert and replace with actual values.
	_, err = stmt.Exec(t, e)
	if err != nil {
		log.Fatalln(err)
	}
}

// Package globals stores the global variables
package globals

import (
	"database/sql"
	"html/template"
)

var (
	// Tpl is container of parsed templates
	Tpl *template.Template
	// Db is a database handler which manages the connection pool
	Db *sql.DB
	// Err for errors
	Err error
)

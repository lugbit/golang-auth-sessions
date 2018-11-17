// Package db connects to the database
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// OpenDB sets up the connection to the database and pings the mysql server
func OpenDB() *sql.DB {
	// Build datasource from environmental variables
	dataSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// Open databse using the mysql driver and the configuration above
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalln(err)
	}
	// Ping the database to test connection
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	// Return databse handler struct

	return db
}

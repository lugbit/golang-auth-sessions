package main

import (
	"html/template"
	"net/http"

	db "./config/db"
	globals "./config/globals"
	handler "./handlers"
	_ "github.com/go-sql-driver/mysql"
	gotenv "github.com/subosito/gotenv"
)

// Parse all .gohtml files in the templates folder
func init() {
	globals.Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	// Load env variables
	gotenv.Load()
}

func main() {
	// Open database handler and defer close until main goroutine exits
	globals.Db = db.OpenDB()
	defer globals.Db.Close()

	// Routes
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/activate-account", handler.ActivateAccount)
	http.HandleFunc("/send-activation", handler.SendActivationLink)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/my-profile", handler.MyProfile)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/my-profile/update-password", handler.UpdatePassword)
	http.HandleFunc("/activation-reminder", handler.ActivationMessage)
	// Force clean sessions for testing
	http.HandleFunc("/clean-sessions", handler.CleanSessions)
	// Route to serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Listen on port 8080 and use default multiplexer
	http.ListenAndServe(":8080", nil)
}

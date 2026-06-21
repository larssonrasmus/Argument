package app

import (
	"html/template"
	"net/http"

	"argument/internal/auth"
	"argument/internal/config"
	"argument/internal/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Page contains generic information for all pages
type Page struct {
	Title string
}

// PostIndex is used for the individual post pages
type PostIndex struct {
	Page Page
	Post database.Post
}

// PostsIndex is used for the index page
type PostsIndex struct {
	Page  Page
	Posts []database.Post
}

var storage = sessions.NewCookieStore([]byte(config.Conf.Session.Key))

// SetupHandler prepares the route handler
func SetupHandler() *mux.Router {
	r := mux.NewRouter()
	t := template.Must(template.ParseGlob(config.Conf.Root + "templates/*.tmpl"))

	page := Page{
		Title: config.Conf.Title,
	}

	// Index page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.tmpl", PostsIndex{
			Page:  page,
			Posts: database.GetPosts(),
		})
	})

	// Creating a post
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := database.Post{Title: r.FormValue("title")}

		database.SavePost(p)

		t.ExecuteTemplate(w, "index.tmpl", PostsIndex{
			Page:  page,
			Posts: database.GetPosts(),
		})
	}).Methods(http.MethodPost)

	// Registration
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		failed := false

		if r.Method == http.MethodPost {
			p := auth.HashPassword(r.FormValue("password"))
			u := auth.NormalizeUsername(r.FormValue("username"))
			user := database.User{Username: u, Password: p}

			id, _ := database.CreateUser(user)

			if id > 0 {
				// Set the user as logged in
				// https://github.com/gorilla/sessions

				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

				return
			}

			failed = true
		}

		t.ExecuteTemplate(w, "register.tmpl", struct {
			Page   Page
			Failed bool
		}{page, failed})
	}).Methods(http.MethodGet, http.MethodPost)

	// Login
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, err := database.VerifyLogin(r.FormValue("username"), r.FormValue("password"))

			if err == nil {
				// Set the user as logged in
				// https://github.com/gorilla/sessions

				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

				return
			}

			// Show message about user not found
		}

		t.ExecuteTemplate(w, "register.tmpl", PostsIndex{
			Page:  page,
			Posts: database.GetPosts(),
		})
	}).Methods(http.MethodGet, http.MethodPost)

	return r
}

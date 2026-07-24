package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/pelletier/go-toml"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB        *gorm.DB
	Sessions  *scs.SessionManager
	Config    config
	Templates map[string]*template.Template
}

type User struct {
	gorm.Model
	Handle   string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type arguments struct {
	config string
	port   int
}

type config struct {
	Title string `toml:"title"`
	Logo  string `toml:"logo"`
}

var args arguments

func (app *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.Sessions.GetInt(r.Context(), "userID") == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (app *App) handleSplash(w http.ResponseWriter, r *http.Request) {
	app.Templates["splash"].Execute(w, app.Config)
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	app.Templates["splash"].Execute(w, app.Config)
}

func (app *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	handle := r.FormValue("handle")
	password := r.FormValue("password")

	var user User
	if err := app.DB.Where("handle = ?", handle).First(&user).Error; err != nil {
		// user not found
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// invalid password
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	app.Sessions.Put(r.Context(), "userID", int(user.ID))

	http.Redirect(w, r, "/bbs", http.StatusSeeOther)
}

func (app *App) handleRegistration(w http.ResponseWriter, r *http.Request) {
	handle := r.FormValue("handle")
	password := r.FormValue("password")
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := User{Handle: handle, Password: string(hashed)}

	app.DB.Create(&user)
	app.Sessions.Put(r.Context(), "userID", int(user.ID))

	http.Redirect(w, r, "/bbs", http.StatusSeeOther)
}

func loadConfig(file string) config {
	cfg := config{}
	f, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("could not read config file:", err)
	}

	if err := toml.Unmarshal(f, &cfg); err != nil {
		log.Fatal("could not parse config:", err)
	}

	return cfg
}

func main() {
	flag.StringVar(&args.config, "config", "config.toml", "config file")
	flag.IntVar(&args.port, "port", 8080, "server port")
	flag.Parse()

	cfg := loadConfig(args.config)

	db, _ := gorm.Open(sqlite.Open("argument.db"), &gorm.Config{})
	db.AutoMigrate(&User{})

	sm := scs.New()
	sm.Store, _ = gormstore.New(db)
	sm.Lifetime = 24 * time.Hour

	app := &App{
		Config:   cfg,
		DB:       db,
		Sessions: sm,
		Templates: map[string]*template.Template{
			"splash": template.Must(template.ParseGlob("templates/splash.html")),
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.handleSplash)
	mux.HandleFunc("POST /login", app.handleLogin)
	mux.HandleFunc("POST /register", app.handleRegistration)
	mux.HandleFunc("POST /logout", app.handleLogout)

	http.ListenAndServe(":"+strconv.Itoa(args.port), sm.LoadAndSave(mux))
}

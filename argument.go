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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB        *gorm.DB
	Sessions  *scs.SessionManager
	Config    config
	Templates map[string]*template.Template
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

func (app *App) handleSplash(w http.ResponseWriter, r *http.Request) {
	app.Templates["splash"].Execute(w, app.Config)
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	app.Templates["login"].Execute(w, app.Config)
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

	sm := scs.New()
	sm.Store, _ = gormstore.New(db)
	sm.Lifetime = 24 * time.Hour

	app := &App{
		Config:   cfg,
		DB:       db,
		Sessions: sm,
		Templates: map[string]*template.Template{
			"splash": template.Must(template.ParseGlob("templates/splash.html")),
			"login":  template.Must(template.ParseGlob("templates/login.html")),
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.handleSplash)

	http.ListenAndServe(":"+strconv.Itoa(args.port), sm.LoadAndSave(mux))
}

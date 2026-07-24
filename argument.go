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

type arguments struct {
	config string
	port   int
}

type config struct {
	Title string `toml:"title"`
	Logo  string `toml:"logo"`
}

var args arguments
var cfg config
var db gorm.DB
var sessionManager *scs.SessionManager

func handler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/splash.html")

	if err := t.Execute(w, cfg); err != nil {
		log.Println(err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func loadConfig(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("could not read config file:", err)
	}

	if err := toml.Unmarshal(f, &cfg); err != nil {
		log.Fatal("could not parse config:", err)
	}

	log.Println("loaded config", file)
}

func main() {
	flag.StringVar(&args.config, "config", "config.toml", "config file")
	flag.IntVar(&args.port, "port", 8080, "server port")
	flag.Parse()

	log.Println("read config, starting up")

	db, err := gorm.Open(sqlite.Open("argument.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	log.Println("connected to SQLite database")

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	if sessionManager.Store, err = gormstore.New(db); err != nil {
		log.Fatal(err)
	}

	log.Println("session manager initialized")
	log.Println("starting server on port", args.port)

	loadConfig(args.config)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	http.ListenAndServe(":"+strconv.Itoa(args.port), sessionManager.LoadAndSave(mux))
}

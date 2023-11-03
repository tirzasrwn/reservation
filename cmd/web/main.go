package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/driver"
	"github.com/tirzasrwn/reservation/internal/handlers"
	"github.com/tirzasrwn/reservation/internal/helpers"
	"github.com/tirzasrwn/reservation/internal/models"
	"github.com/tirzasrwn/reservation/internal/render"
)

const portNumber = ":4545"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// What am I going to put in the session.
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database.
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=reservation user=postgres password=postgres")
	if err != nil {
		fmt.Println("Cannot connect to the database!")
		return nil, err
	}
	log.Println("Connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("Cannot create template cache!")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	return db, nil
}

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/spf13/viper"
	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/driver"
	"github.com/tirzasrwn/reservation/internal/handlers"
	"github.com/tirzasrwn/reservation/internal/helpers"
	"github.com/tirzasrwn/reservation/internal/models"
	"github.com/tirzasrwn/reservation/internal/render"
)

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
)

func main() {
	err := initializeAppConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()
	fmt.Println("Start mail listener...")

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", app.Port),
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
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		app.DBHost, app.DBPort, app.DBName, app.DBUser, app.DBPassword))
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

	fmt.Printf("Starting application on port %d\n", app.Port)

	return db, nil
}

func initializeAppConfig() error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AllowEmptyEnv(false)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	app.Port = viper.GetInt("PORT")
	app.DBUser = viper.GetString("DB_USER")
	app.DBName = viper.GetString("DB_NAME")
	app.DBPassword = viper.GetString("DB_PASSWORD")
	app.DBHost = viper.GetString("DB_HOST")
	app.DBPort = viper.GetInt("DB_PORT")

	log.Println("[INIT] configuration loaded")

	return nil
}

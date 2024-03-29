package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/tirzasrwn/reservation/internal/models"
)

// AppConfig holds the application config.
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData

	Port       int
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int
}

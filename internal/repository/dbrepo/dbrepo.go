package dbrepo

import (
	"database/sql"

	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

package repository

import "github.com/tirzasrwn/reservation/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}

package repository

import (
	"time"

	"github.com/tirzasrwn/reservation/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(reservation models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, peocessed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockByID(id int) error
	CountAllReservations() (int, error)
	AllReservationsPagination(page, limitPerPage int) (models.Pagination, error)
	CountAllNewReservations() (int, error)
	AllNewReservationsPagination(page, limitPerPage int) (models.Pagination, error)
}

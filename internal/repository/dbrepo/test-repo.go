package dbrepo

import (
	"time"

	"github.com/tirzasrwn/reservation/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reseration into database.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets room by ID.
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}

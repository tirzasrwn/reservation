package dbrepo

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/tirzasrwn/reservation/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reseration into database.
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into database.
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (
    start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id
  ) values ($1, $2, $3, $4, $5, $6, $7);`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exist for roomID.
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;`
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	var numRows int
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range.
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select r.id, r.room_name from rooms r where r.id not in (
    select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date
  );`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// GetRoomByID gets room by ID.
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, room_name, created_at, updated_at from rooms where id = $1;`
	var room models.Room
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID gets user by ID.
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1;`
	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates a user in the database.
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
  where id = $6`
	_, err := m.DB.ExecContext(ctx, query,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.AccessLevel,
		time.Now(),
		&user.ID,
	)
	return err
}

// Authenticate authenticate the user.
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hasedPassowrd string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hasedPassowrd)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hasedPassowrd), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hasedPassowrd, nil
}

// AllReservation returs a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
  `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// AllNewReservations returs a slice of all new reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
    where processed = 0
		order by r.start_date asc
  `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// GetRoomByID returns a reservation by id.
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
  select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
  r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
  from reservations r
  left join rooms rm on (r.room_id = rm.id)
  where r.id = $1
  `

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

// UpdateReservation updates a reservation.
func (m *postgresDBRepo) UpdateReservation(reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
  where id = $6`
	_, err := m.DB.ExecContext(ctx, query,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		time.Now(),
		&reservation.ID,
	)
	return err
}

// DeleteReservation deletes a reservation by id.
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProcessedForReservation updates processed reservation by id.
func (m *postgresDBRepo) UpdateProcessedForReservation(id, peocessed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
  update reservations set processed = $1 where id = $2
  `
	_, err := m.DB.ExecContext(ctx, query, peocessed, id)
	if err != nil {
		return err
	}
	return nil
}

// AllRooms returns all rooms.
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Room
		err := rows.Scan(
			&r.ID,
			&r.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate resturns slice of room restrictions by date.
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
		select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
		from room_restrictions where $1 < end_date and $2 >= start_date
		and room_id = $3
`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction.
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	// restriction_id is 2 because it is a block.
	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// DeleteBlockByID deletes a room restriction.
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`
	// restriction_id is 2 because it is a block.
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) CountAllReservations() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) from reservations`
	row := m.DB.QueryRowContext(ctx, query)

	result := 0
	err := row.Scan(
		&result,
	)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// AllReservation returs a slice of all reservations
func (m *postgresDBRepo) AllReservationsPagination(page, limitPerPage int) (models.Pagination, error) {
	var p models.Pagination
	p.CurrentPage = page
	p.LimitPerPage = limitPerPage

	count, err := m.CountAllReservations()
	if err != nil {
		return p, err
	}
	p.TotalItems = count
	offset := (p.CurrentPage - 1) * p.LimitPerPage

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
    limit $1
    offset $2
  `

	rows, err := m.DB.QueryContext(ctx, query, p.LimitPerPage, offset)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	countRows := 0
	for rows.Next() {
		countRows++
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return p, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return p, err
	}

	p.Items = reservations
	p.TotalPage = int(math.Ceil(float64(p.TotalItems) / float64(p.LimitPerPage)))
	p.TotalRows = countRows

	return p, nil
}

// CountAllNewReservations returns the number of new reservations.
func (m *postgresDBRepo) CountAllNewReservations() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) from reservations where processed = 0`
	row := m.DB.QueryRowContext(ctx, query)

	result := 0
	err := row.Scan(
		&result,
	)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// AllNewReservationsPagination returns pagination for all new reservations.
func (m *postgresDBRepo) AllNewReservationsPagination(page, limitPerPage int) (models.Pagination, error) {
	var p models.Pagination
	p.CurrentPage = page
	p.LimitPerPage = limitPerPage

	count, err := m.CountAllNewReservations()
	if err != nil {
		return p, err
	}
	p.TotalItems = count
	offset := (p.CurrentPage - 1) * p.LimitPerPage

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
    where processed = 0
		order by r.start_date asc
    limit $1
    offset $2
  `

	rows, err := m.DB.QueryContext(ctx, query, p.LimitPerPage, offset)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	countRows := 0
	for rows.Next() {
		countRows++
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return p, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return p, err
	}

	p.Items = reservations
	p.TotalPage = int(math.Ceil(float64(p.TotalItems) / float64(p.LimitPerPage)))
	p.TotalRows = countRows

	return p, nil
}

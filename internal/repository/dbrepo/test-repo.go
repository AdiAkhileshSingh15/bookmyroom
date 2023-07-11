package dbrepo

import (
	"time"

	"github.com/AdiAkhileshSingh15/bookmyroom/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

func (m *testDBRepo) SearchAvailabilityForDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	return []models.Room{}, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	return models.Room{}, nil
}

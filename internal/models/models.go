package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	ID        int       `json:"id"`
	RoomName  string    `json:"room_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Restriction struct {
	ID              int       `json:"id"`
	RestrictionName string    `json:"restriction_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Reservation struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RoomID    int       `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Room      Room      `json:"room"`
}

type RoomRestriction struct {
	ID            int       `json:"id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	RoomID        int       `json:"room_id"`
	ReservationID int       `json:"reservation_id"`
	RestrictionID int       `json:"restriction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Room          Room      `json:"room"`
	Reservation   Reservation
	Restriction   Restriction
}

type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	RoomID   int
	Template string
}

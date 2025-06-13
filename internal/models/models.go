// Package models
package models

import "gorm.io/gorm"

// --- STRUCT UNTUK DATABASE ---

type User struct {
	ID             string        `gorm:"primaryKey" json:"id"` // NIM
	Name           string        `json:"name"`
	Faculty        string        `json:"faculty"`
	Major          string        `json:"major"`
	CreditScore    int           `json:"creditScore"`
	ProfilePictURL string        `json:"profilePictUrl"`
	PasswordHash   string        `json:"-"` // Tidak akan pernah dikirim dalam JSON
	Reservations   []Reservation `gorm:"foreignKey:BorrowerID" json:"-"`
}

type TVInfo struct {
	ID           int           `gorm:"primaryKey;autoIncrement" json:"id"`
	ConsoleType  string        `json:"consoleType"`
	Games        []*Game       `gorm:"many2many:tv_info_games;" json:"gameList"`
	Reservations []Reservation `gorm:"foreignKey:TVID" json:"-"`
}

type Game struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string    `json:"title"`
	CoverPictURL string    `json:"coverPictUrl"`
	TVs          []*TVInfo `gorm:"many2many:tv_info_games;" json:"-"`
}

type Reservation struct {
	gorm.Model
	TVID       int
	BorrowerID string
	TimeSlot   string
}

// --- STRUCT UNTUK REQUEST & RESPONSE BODY ---

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type ErrorData struct {
	ErrorMsg string `json:"errorMsg"`
}

type ErrorResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   ErrorData `json:"data"`
}

type LoginBody struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type TokenCarrier struct {
	AuthToken  string `json:"authToken"`
	UserID     string `json:"userId"`
	ExpireDate string `json:"expireDate"`
}

type History struct {
	ID                  int    `json:"id"`
	TVID                int    `json:"tvId"`
	ReservationDateTime string `json:"reservationDateTime"`
	TVPictURL           string `json:"tvPictUrl"`
}

type TimeSlot struct {
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Availability string `json:"availability"`
}

type TV struct {
	ID          int        `json:"id"`
	ConsoleType string     `json:"consoleType"`
	TimeSlots   []TimeSlot `json:"timeSlots"`
}

type ReservationBody struct {
	TVID       int    `json:"tvId"`
	BorrowerID string `json:"borrowerId"`
	Timeslot   string `json:"timeslot"`
}

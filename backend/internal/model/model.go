package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusName string

const (
	Success  StatusName = "success"
	Failure  StatusName = "failure"
	Occupied StatusName = "occupied"
)

// Record описывает одну сессию парковки.
// Поля соответствуют документу коллекции parking_records.
type Record struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Token        string             `bson:"token"`
	// PassCode     string     `bson:"pass_code" json:"pass_code"`
	ClientName   string    `bson:"client_name" json:"client_name"`
	PhoneNumber  string    `bson:"phone_number" json:"phone_number"`
	LicensePlate string    `bson:"license_plate" json:"license_plate"`
	SpotNumber   int       `bson:"spot_number" json:"spot_number"`
	StartTime    time.Time `bson:"start_time" json:"start_time"`
	EndTime      time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Status       string    `bson:"status" json:"status"`
	// Role         Role       `bson:"role" json:"role"`
}

type RecordDto struct {
	ClientName   string     `json:"client_name"`
	PhoneNumber  string     `json:"phone_number"`
	LicensePlate string     `json:"license_plate"`
	SpotNumber   int        `json:"spot_number"`
	Duration     *string    `json:"duration,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
}

type Response struct {
	// Token    string     `json:"token"`
	// PassCode string     `json:"pass_code"`
	Status StatusName `json:"status"`
}

// type AuthDto struct {
// PhoneNumber string `json:"phone_number"`
// PassCode    string `json:"pass_code"`
// }

type ProlongSessionDto struct {
	Duration string `json:"duration"`
}

type PaymentReceipt struct {
	Cost      int        `json:"cost"`
	Duration  *string    `json:"duration"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

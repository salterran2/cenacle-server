package person

import "time"

//Person object
type Person struct {
	ID           int64     `json:"id"`
	Firstname    string    `json:"firstname" validate:"required"`
	Lastname     string    `json:"lastname" validate:"required"`
	EmailAddress string    `json:"email_address" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

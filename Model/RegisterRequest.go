package Model

import "time"

type RegisterRequest struct {
	ID     uint   `xml:"-" json:"-"`
	Token  string `xml:"-" json:"-"`
	Email  string
	Expiry time.Time
}

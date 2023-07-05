package Model

import "time"

type RequestAffair int

const (
	Register      RequestAffair = 0
	ResetPassword RequestAffair = 1
	ChangeEmail   RequestAffair = 2
)

type RegisterRequest struct {
	ID     uint   `xml:"-" json:"-"`
	Token  string `xml:"-" json:"-"`
	Affair RequestAffair
	Expiry time.Time

	UserID uint   `xml:",omitempty" json:",omitempty"`
	Email  string `xml:",omitempty" json:",omitempty"`
}

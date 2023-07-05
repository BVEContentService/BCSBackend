package Model

import "github.com/jinzhu/gorm"

type Privilege int

const (
	Normal    Privilege = 0
	Validator Privilege = 10
	Moderator Privilege = 50
	SiteAdmin Privilege = 100
)

type Uploader struct {
	gorm.Model
	Developer
	Validated bool
	//Username    string
	Password    string `xml:"-" json:"-"`
	Description string `gorm:"type:text"`
	Privilege   Privilege
	Packages    []Package `xml:">Package,omitempty" json:",omitempty"`
}

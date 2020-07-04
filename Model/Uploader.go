package Model

import "github.com/jinzhu/gorm"

type Privilege int
const (
    Normal          Privilege = 0
    Validator       Privilege = 10
    Moderator       Privilege = 50
    SiteAdmin       Privilege = 100
)

type Uploader struct {
    gorm.Model
    Developer
    Validated       bool
    Username        string
    Password        []byte         `xml:",omitempty" json:",omitempty"`
    Description     NullString     `gorm:"type:text"`
    Privilege       Privilege
    Packages        []Package
}

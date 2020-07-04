package Model

import (
	"github.com/jinzhu/gorm"
)

type Package struct {
	gorm.Model
	Identifier      string
	GUID            NullString      `xml:",omitempty" json:",omitempty"`
	Name            String3         `gorm:"embedded;embedded_prefix:name_"`
	UploaderID      uint
	Uploader        *Uploader       `xml:",omitempty" json:",omitempty"`
	Author          *Developer      `gorm:"embedded;embedded_prefix:author_" xml:",omitempty" json:",omitempty"`
	Homepage        NullString      `xml:",omitempty" json:",omitempty"`
	Thumbnail       NullString      `xml:",omitempty" json:",omitempty"`
	ThumbnailLQ     NullString      `xml:",omitempty" json:",omitempty"`
	Description     NullString      `xml:",omitempty" json:",omitempty" gorm:"type:text"`
	IsRepost        bool            `xml:",omitempty" json:",omitempty" gorm:"-"` // For POST requests only
	Files           []File
}

func (p *Package) AfterFind() (err error) {
	if *p.Author == NULL_DEVELOPER {
		p.Author = nil
	}
	return
}
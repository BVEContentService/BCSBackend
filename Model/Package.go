package Model

import (
	"github.com/jinzhu/gorm"
)

type Package struct {
	gorm.Model
	Identifier  string
	GUID        string  `xml:",omitempty" json:",omitempty"`
	Name        String3 `gorm:"embedded;embedded_prefix:name_"`
	UploaderID  uint
	Uploader    *Uploader      `xml:",omitempty" json:",omitempty"`
	Author      *Developer     `xml:",omitempty" json:",omitempty" gorm:"embedded;embedded_prefix:author_"`
	Homepage    string         `xml:",omitempty" json:",omitempty"`
	Thumbnail   string         `xml:",omitempty" json:",omitempty"`
	ThumbnailLQ string         `xml:",omitempty" json:",omitempty"`
	Description string         `xml:",omitempty" json:",omitempty" gorm:"type:text"`
	Files       []File         `xml:",omitempty" json:",omitempty"`
	Platforms   []PlatformType `xml:",omitempty" json:",omitempty" gorm:"-"` // Collected on the fly
}

func (p *Package) AfterFind() (err error) {
	if *p.Author == NULL_DEVELOPER {
		p.Author = nil
	}
	return
}

package Model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Package struct {
	gorm.Model
	Identifier  string
	GUID        string  `xml:",omitempty" json:",omitempty"`
	Name        String3 `gorm:"embedded;embedded_prefix:name_"`
	UploaderID  uint
	Uploader    *Uploader  `xml:",omitempty" json:",omitempty"`
	Author      *Developer `xml:",omitempty" json:",omitempty" gorm:"embedded;embedded_prefix:author_"`
	Homepage    string     `xml:",omitempty" json:",omitempty"`
	ForcePopup  bool
	Thumbnail   string         `xml:",omitempty" json:",omitempty"`
	ThumbnailLQ string         `xml:",omitempty" json:",omitempty"`
	Description string         `xml:",omitempty" json:",omitempty" gorm:"type:text"`
	Files       []File         `xml:">File"`                                              /*`xml:",omitempty" json:",omitempty"`*/
	Platforms   []PlatformType `xml:">PlatformType,omitempty" json:",omitempty" gorm:"-"` // Collected on the fly
}

func (p *Package) AfterFind() (err error) {
	if *p.Author == NULL_DEVELOPER {
		p.Author = nil
	}
	return
}

func (p *Package) MatchKeyword(keyword string) bool {
	if keyword == "" {
		return true
	} else {
		nameToSearch := strings.ToLower(fmt.Sprintf("%s %s %s", p.Name.Local, p.Name.English, p.Name.Tag))
		return strings.Contains(nameToSearch, keyword)
	}
}

func (p *Package) MatchPlatform(platform PlatformType) bool {
	if platform == 0 {
		return true
	} else if platform == -1 {
		return len(p.Platforms) > 0
	} else {
		return inSlice(p.Platforms, platform)
	}
}

func inSlice(haystack []PlatformType, needle PlatformType) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}
	return false
}

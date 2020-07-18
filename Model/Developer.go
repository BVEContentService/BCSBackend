package Model

type Developer struct {
	Name     String3 `gorm:"embedded;embedded_prefix:name_"`
	Email    string  `xml:",omitempty" json:",omitempty"`
	Homepage string  `xml:",omitempty" json:",omitempty"`
}

var NULL_DEVELOPER = Developer{}

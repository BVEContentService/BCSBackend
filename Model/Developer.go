package Model

type Developer struct {
    Name            String3         `gorm:"embedded;embedded_prefix:name_"`
    Email           NullString      `xml:",omitempty" json:",omitempty"`
    Homepage        NullString      `xml:",omitempty" json:",omitempty"`
}

var NULL_DEVELOPER = Developer {}
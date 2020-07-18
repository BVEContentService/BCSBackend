package Model

import "strings"

type String3 struct {
	Local   string
	English string
	Tag     string
}

func (s3 *String3) TrimNames() {
	s3.Local = strings.TrimSpace(s3.Local)
	s3.English = strings.TrimSpace(s3.English)
	s3.Tag = strings.TrimSpace(s3.Tag)
}

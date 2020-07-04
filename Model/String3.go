package Model

type String3 struct {
    Local       NullString
    English     NullString
    Tag         NullString
}

func (s3 *String3) TrimNames() {
    s3.Local.TrimSpace()
    s3.English.TrimSpace()
    s3.Tag.TrimSpace()
}
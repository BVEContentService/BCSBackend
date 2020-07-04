package Model

import (
    "database/sql"
    "database/sql/driver"
    "encoding/json"
    "encoding/xml"
    "strings"
)

type NullString sql.NullString
var NullStringNull = NullString{ String:"", Valid:false }

func (ns *NullString) Scan(value interface{}) error {
    var rns sql.NullString
    err := rns.Scan(value);
    if err != nil {
        return err
    }
    ns.Valid = rns.Valid
    ns.String = rns.String
    return nil
}

func (ns NullString) Value() (driver.Value, error) {
    rns := sql.NullString{ String: ns.String, Valid: ns.Valid }
    return rns.Value()
}

func (x NullString) MarshalJSON() ([]byte, error) {
    if !x.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(x.String)
}
func (x *NullString) UnmarshalJSON(data []byte) error {
    var vstr *string
    if err := json.Unmarshal(data, &vstr); err != nil {
        return err
    }
    if vstr != nil {
        x.Valid = true
        x.String = *vstr
    } else {
        x.Valid = false
        x.String = ""
    }
    return nil
}
func (x NullString) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    if !x.Valid {
        return nil
    }
    return e.EncodeElement(x.String, start)
}
func (x *NullString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var s string
    if err := d.DecodeElement(&s, &start); err != nil { return err }
    x.String = s
    x.Valid = true
    return nil
}

func (x *NullString) TrimSpace() {
    if !x.Valid { return }
    x.String = strings.TrimSpace(x.String)
}

func (x *NullString) NotEmpty() bool {
    return x.Valid && x.String != ""
}
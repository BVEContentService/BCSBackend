package Model

import (
    "OBPkg/Config"
    "OBPkg/Utility"
    "encoding/json"
    "encoding/xml"
)

type PlatformType int
func (pt PlatformType) String() string {
    return Config.CurrentConfig.Platform[int(pt)];
}
func (pt *PlatformType) MarshalJSON() ([]byte, error) {
    return json.Marshal(pt.String())
}
func (pt *PlatformType) UnmarshalJSON(data []byte) error {
    ptInt, ok := mapkey(Config.CurrentConfig.Platform, string(data))
    if !ok {
        return Utility.ERR_BAD_PARAMETER
    }
    *pt = PlatformType(ptInt)
    return nil
}
func (pt *PlatformType) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    return e.EncodeElement(pt.String(), start)
}
func (pt *PlatformType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var s string
    if err := d.DecodeElement(&s, &start); err != nil { return err }
    ptInt, ok := mapkey(Config.CurrentConfig.Platform, s)
    if !ok {
        return Utility.ERR_BAD_PARAMETER
    }
    *pt = PlatformType(ptInt)
    return nil
}

type ServiceType int
func (st ServiceType) String() string {
    return Config.CurrentConfig.FileService.DatabaseMap[int(st)];
}
func (st *ServiceType) MarshalJSON() ([]byte, error) {
    return json.Marshal(st.String())
}
func (st *ServiceType) UnmarshalJSON(data []byte) error {
    ptInt, ok := mapkey(Config.CurrentConfig.FileService.DatabaseMap, string(data))
    if !ok {
        return Utility.ERR_BAD_PARAMETER
    }
    *st = ServiceType(ptInt)
    return nil
}
func (st *ServiceType) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
    return e.EncodeElement(st.String(), start)
}
func (st *ServiceType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var s string
    if err := d.DecodeElement(&s, &start); err != nil { return err }
    stInt, ok := mapkey(Config.CurrentConfig.FileService.DatabaseMap, s)
    if !ok {
        return Utility.ERR_BAD_PARAMETER
    }
    *st = ServiceType(stInt)
    return nil
}

func mapkey(m map[int]string, value string) (key int, ok bool) {
    for k, v := range m {
        if v == value {
            key = k
            ok = true
            return
        }
    }
    return
}

type File struct {
    ID              uint
    PackageID       uint
    Package         *Package        `xml:",omitempty" json:",omitempty"`
    Platform        *PlatformType
    Validated       *bool
    Version         string
    Service         *ServiceType    `xml:",omitempty" json:",omitempty"`
    URLParam        NullString      `xml:",omitempty" json:",omitempty"`
    AuthParam       NullString      `xml:",omitempty" json:",omitempty"`
    FetchURL        string          `xml:",omitempty" json:",omitempty" gorm:"-"`
}
package Model

import (
	"OBPkg/Config"
	"OBPkg/Utility"
	"encoding/json"
	"encoding/xml"
)

type PlatformType int

func (pt PlatformType) String() string {
	if int(pt) == -1 {
		return "+"
	}
	ptName, ok := mapKey(Config.CurrentConfig.Platform.DatabaseMap, int(pt))
	if !ok {
		return ""
	}
	return ptName
}
func (pt PlatformType) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}
func (pt *PlatformType) UnmarshalJSON(data []byte) error {
	var str string
	if json.Unmarshal(data, &str) != nil {
		return Utility.ERR_BAD_PARAMETER.WithData("Unmarshal, PlatformType")
	}
	if str == "+" {
		*pt = PlatformType(-1)
		return nil
	}
	ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[str]
	if !ok {
		return Utility.ERR_BAD_PARAMETER.WithData("PlatformType not in DatabaseMap")
	}
	*pt = PlatformType(ptInt)
	return nil
}
func (pt PlatformType) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return e.EncodeElement(pt.String(), start)
}
func (pt *PlatformType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return Utility.ERR_BAD_PARAMETER.WithData("Unmarshal, PlatformType")
	}
	if s == "+" {
		*pt = PlatformType(-1)
		return nil
	}
	ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[s]
	if !ok {
		return Utility.ERR_BAD_PARAMETER.WithData("PlatformType not in DatabaseMap")
	}
	*pt = PlatformType(ptInt)
	return nil
}

type ServiceType int

func (st ServiceType) String() string {
	stName, ok := mapKey(Config.CurrentConfig.FileService.DatabaseMap, int(st))
	if !ok {
		return ""
	}
	return stName
}
func (st ServiceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(st.String())
}
func (st *ServiceType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return Utility.ERR_BAD_PARAMETER.WithData("Unmarshal, ServiceType")
	}
	ptInt, ok := Config.CurrentConfig.FileService.DatabaseMap[str]
	if !ok {
		return Utility.ERR_BAD_PARAMETER.WithData("ServiceType not in DatabaseMap")
	}
	*st = ServiceType(ptInt)
	return nil
}
func (st ServiceType) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return e.EncodeElement(st.String(), start)
}
func (st *ServiceType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return Utility.ERR_BAD_PARAMETER.WithData("Unmarshal, ServiceType")
	}
	stInt, ok := Config.CurrentConfig.FileService.DatabaseMap[s]
	if !ok {
		return Utility.ERR_BAD_PARAMETER.WithData("ServiceType not in DatabaseMap")
	}
	*st = ServiceType(stInt)
	return nil
}

type File struct {
	ID             uint
	PackageID      uint
	Package        *Package `xml:"-" json:"-"`
	Platform       PlatformType
	Validated      bool
	NeedValidation bool   `gorm:"-"`
	RejectReason   string `xml:",omitempty" json:",omitempty"`
	Version        string
	Size           string
	Service        ServiceType `xml:",omitempty" json:",omitempty"`
	URLParam       string      `xml:",omitempty" json:",omitempty"`
	AuthParam      string      `xml:",omitempty" json:",omitempty"`
	FetchURL       string      `xml:",omitempty" json:",omitempty" gorm:"-"`
}

func (f *File) AfterFind() (err error) {
	f.NeedValidation = Config.CurrentConfig.Platform.NeedValidation[f.Platform.String()]
	return
}

func mapKey(m map[string]int, value int) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

package Config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

type ApiConfig struct {
	Gin           GinConfig
	MySql         DatabaseConfig
	JWT           JWTConfig
	SMTP          SMTPConfig
	GeoIPDatabase string `json:"geoip_database"`
	Platform      PlatformConfig
	FileService   FileServiceConfig `json:"file_service"`
}
type GinConfig struct {
	Debug       bool
	Address     string
	TLS         bool
	CertFile    string   `json:"cert_file"`
	KeyFile     string   `json:"key_file"`
	AllowOrigin []string `json:"allow_origin"`
}
type DatabaseConfig struct {
	Debug    bool
	Host     string
	Database string
	Username string
	Password string
}
type JWTConfig struct {
	SecretKey  string `json:"secret_key"`
	Timeout    Duration
	MaxRefresh Duration `json:"max_refresh"`
}
type SMTPConfig struct {
	Host          string
	Username      string
	Password      string
	TokenDuration Duration `json:"token_duration"`
}
type PlatformConfig struct {
	DatabaseMap    map[string]int  `json:"database_map"`
	NeedValidation map[string]bool `json:"need_validation"`
}
type FileServiceConfig struct {
	DatabaseMap map[string]int    `json:"database_map"`
	URLMap      map[string]string `json:"url_map"`
}

var CurrentConfig ApiConfig

func InitConfig(ConfigFile string) {
	rawData, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic("Failed reading config file")
	}
	err = json.Unmarshal(rawData, &CurrentConfig)
	if err != nil {
		panic("Failed parsing config file: " + err.Error())
	}
}

const DbDialect = "mysql"

func GetConnectionString() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		CurrentConfig.MySql.Username, CurrentConfig.MySql.Password,
		CurrentConfig.MySql.Host, CurrentConfig.MySql.Database)
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		dp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

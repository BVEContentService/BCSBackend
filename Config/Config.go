package Config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ApiConfig struct {
	Gin             GinConfig
	MySql           DatabaseConfig
	JWT             JWTConfig
	GeoIPDatabase   string              `json:"geoip_database"`
	Platform        map[int]string
	FileService     FileServiceConfig   `json:"file_service"`
}
type GinConfig struct {
	Debug           bool
	Address         string
	TLS             bool
	CertFile        string              `json:"cert_file"`
	KeyFile         string              `json:"key_file"`
	AllowOrigin     []string            `json:"allow_origin"`
}
type DatabaseConfig struct {
	Debug           bool
	Host            string
	Database        string
	Username        string
	Password        string
}
type JWTConfig struct {
	SecretKey       string              `json:"secret_key"`
}
type FileServiceConfig struct {
	DatabaseMap     map[int]string      `json:"database_map"`
	URLMap          map[string]string   `json:"url_map"`
}

var CurrentConfig ApiConfig

func InitConfig(ConfigFile string) {
	rawData, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic("Failed reading config file")
	}
	err = json.Unmarshal(rawData, &CurrentConfig);
	if err != nil {
		panic("Failed parsing config file")
	}
	fmt.Print(CurrentConfig)
}

const DbDialect = "mysql"

func GetConnectionString() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		CurrentConfig.MySql.Username, CurrentConfig.MySql.Password,
		CurrentConfig.MySql.Host, CurrentConfig.MySql.Database)
}

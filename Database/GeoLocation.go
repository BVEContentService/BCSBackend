package Database
import (
    "OBPkg/Config"
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "net"
)

func GetGeoDB() *geoip2.Reader {
    /*if GeoDB != nil {
        var err error
        GeoDB, err = geoip2.Open(Config.CurrentConfig.GeoIPDatabase)
        if err != nil {
            panic(err)
        }
    }
    defer GeoDB.Close()
    return GeoDB*/
    geoDB, err := geoip2.Open(Config.CurrentConfig.GeoIPDatabase)
    if err != nil {
        panic(err)
    }
    defer geoDB.Close()
    return geoDB
}

func GetIPCountryCode(IP string) string {
    ip := net.ParseIP(IP)
    fmt.Println(IP)
    record, err := GetGeoDB().Country(ip)
    if err != nil {
        return ""
    } else {
        return record.Country.IsoCode
    }
}
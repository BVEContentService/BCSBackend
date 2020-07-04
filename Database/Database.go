package Database

import (
	"OBPkg/Config"
	"OBPkg/Model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var (
	DB *gorm.DB
)

func GetDB() *gorm.DB {
	if DB == nil {
		log.Println("Establishing Database Connection...")
		var err error
		DB, err = gorm.Open(Config.DbDialect, Config.GetConnectionString())
		if err != nil {
			DB = nil
			log.Println(err)
			return nil
		}
		DB.LogMode(Config.CurrentConfig.MySql.Debug)
		DB.AutoMigrate(&Model.Uploader{}, &Model.File{}, &Model.Package{})
		log.Println("Database Connection has been Established!")
	}
	return DB
}

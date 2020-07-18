package main

import (
	"OBPkg/Config"
	"OBPkg/Router"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ymzuiku/hit"
)

func main() {
	cfgFile := flag.String("c", "./obpkg.json", "Config file used by the server")
	flag.Parse()
	Config.InitConfig(*cfgFile)
	gin.SetMode(hit.If(Config.CurrentConfig.Gin.Debug, gin.DebugMode, gin.ReleaseMode).(string))
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	if len(Config.CurrentConfig.Gin.AllowOrigin) == 0 ||
		(len(Config.CurrentConfig.Gin.AllowOrigin) == 1 && Config.CurrentConfig.Gin.AllowOrigin[0] == "*") {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = Config.CurrentConfig.Gin.AllowOrigin
	}
	corsConfig.AllowMethods = []string{"OPTION", "HEAD", "GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Range"}
	corsConfig.ExposeHeaders = []string{"Accept-Ranges", "Content-Range"}
	router.Use(cors.New(corsConfig))
	Router.SetRoutes(router)
	var err error
	if Config.CurrentConfig.Gin.TLS {
		err = router.RunTLS(Config.CurrentConfig.Gin.Address, Config.CurrentConfig.Gin.CertFile, Config.CurrentConfig.Gin.KeyFile)
	} else {
		err = router.Run(Config.CurrentConfig.Gin.Address)
	}
	if err != nil {
		panic(err)
	}
}

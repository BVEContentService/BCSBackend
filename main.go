package main

import (
	"OBPkg/Config"
	"OBPkg/Router"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ymzuiku/hit"
	"net/http"
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
		print("All !\n");
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = Config.CurrentConfig.Gin.AllowOrigin
	}
	corsConfig.AllowMethods = []string {"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string {"Origin", "Accept", "Authorization", "Content-Type"}
	router.Use(Cors())
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


// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

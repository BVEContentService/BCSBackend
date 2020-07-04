package Router

import (
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
)

func NoRouteHandler(c *gin.Context) {
    Utility.ERR_NOT_FOUND_URL.PrintResponse(c)
}

func NoMethodHandler(c *gin.Context) {
    Utility.ERR_NOT_ALLOWED.PrintResponse(c)
}

func SetRoutes(router *gin.Engine) {
    router.NoRoute(NoRouteHandler)
    router.NoMethod(NoMethodHandler)
    SetAuthorizationRoutes(router)
    SetPackageRoutes(router)
    SetUploaderRoutes(router)
}
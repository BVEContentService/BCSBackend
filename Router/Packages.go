package Router

import (
    "OBPkg/Controller"
    "OBPkg/Middleware"
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
)

func SetPackageRoutes(router *gin.Engine) {
    router.GET("/packages", Utility.Wrapper(Controller.PackageList))
    router.GET("/packages/:id", Utility.Wrapper(Controller.PackageGet))
    var middleware = Middleware.GetAuthMiddleware()
    router.PUT("/packages", middleware.MiddlewareFunc(), Utility.Wrapper(Controller.PackagePut))
    router.POST("/packages/:id", middleware.MiddlewareFunc(), Utility.Wrapper(Controller.PackagePost))
    router.DELETE("/packages/:id", middleware.MiddlewareFunc(), Utility.Wrapper(Controller.PackageDelete))
}

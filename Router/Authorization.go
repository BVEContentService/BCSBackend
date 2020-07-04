package Router

import (
    "OBPkg/Middleware"
    "github.com/gin-gonic/gin"
)

func SetAuthorizationRoutes(router *gin.Engine) {
    middleware := Middleware.GetAuthMiddleware()
    router.POST("/auth/login", middleware.LoginHandler)
    router.POST("/auth/refresh_token", middleware.RefreshHandler)
}

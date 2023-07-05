package Router

import (
	"OBPkg/Controller"
	"OBPkg/Middleware"
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
	var authMiddleware = Middleware.GetAuthMiddleware().MiddlewareFunc()
	var optionalMiddleware = func(c *gin.Context) {
		if header := c.Request.Header.Get("Authorization"); header != "" {
			authMiddleware(c)
		}
	}
	router.NoRoute(NoRouteHandler)
	router.NoMethod(NoMethodHandler)

	router.POST("/auth/login", Middleware.GetAuthMiddleware().LoginHandler)
	router.GET("/auth/refresh_token", Middleware.GetAuthMiddleware().RefreshHandler)
	router.POST("/auth/register", Utility.Wrapper(Controller.AuthRegister))
	router.POST("/auth/check_token", Utility.Wrapper(Controller.AuthCheckToken))
	router.POST("/auth/activate", Utility.Wrapper(Controller.AuthActivate))
	router.GET("/auth/activate/:token", Utility.Wrapper(Controller.AuthEmailAffair))
	router.POST("/auth/change_password", authMiddleware, Utility.Wrapper(Controller.AuthChangePassword))

	router.HEAD("/packages", Utility.Wrapper(Controller.PackageHead))
	router.GET("/packages", Utility.Wrapper(Controller.PackageList))
	router.GET("/packages/:id", optionalMiddleware, Utility.Wrapper(Controller.PackageGet))
	router.PUT("/packages", authMiddleware, Utility.Wrapper(Controller.PackagePut))
	router.POST("/packages/:id", authMiddleware, Utility.Wrapper(Controller.PackagePost))
	router.DELETE("/packages/:id", authMiddleware, Utility.Wrapper(Controller.PackageDelete))

	router.HEAD("/files", Utility.Wrapper(Controller.FileHead))
	router.GET("/files", optionalMiddleware, Utility.Wrapper(Controller.FileList))
	router.GET("/files/:id", authMiddleware, Utility.Wrapper(Controller.FileGet))
	router.PUT("/files", authMiddleware, Utility.Wrapper(Controller.FilePut))
	router.POST("/files/:id", authMiddleware, Utility.Wrapper(Controller.FilePost))
	router.DELETE("/files/:id", authMiddleware, Utility.Wrapper(Controller.FileDelete))

	router.GET("/uploaders/:id", Utility.Wrapper(Controller.UploaderGet))
	router.POST("/uploaders/:id", authMiddleware, Utility.Wrapper(Controller.UploaderPost))
}

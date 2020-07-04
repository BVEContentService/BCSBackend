package Router

import (
    "OBPkg/Controller"
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
)

func SetUploaderRoutes(router *gin.Engine) {
    router.GET("/uploaders/:id", Utility.Wrapper(Controller.UploaderGet))
}
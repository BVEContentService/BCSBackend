package Controller

import (
    "OBPkg/Database"
    "OBPkg/Model"
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
    "strconv"
)

func UploaderGet(c *gin.Context) error {
    db := Database.GetDB()
    if db == nil { return Utility.ERR_DATABASE }
    upid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var uploaderModel Model.Uploader
    if db.Preload("Packages").First(&uploaderModel, uint(upid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    } else {
        Utility.MarshalResponse(c, 200, uploaderModel)
    }
    return nil
}
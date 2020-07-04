package Controller

import (
    "OBPkg/Config"
    "OBPkg/Database"
    "OBPkg/Middleware"
    "OBPkg/Model"
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
    "strconv"
    "strings"
)

func FileGet(c *gin.Context) error {
    db := Database.GetDB()
    if db == nil { return Utility.ERR_DATABASE }
    fid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var fileModel Model.File;
    var privilege = c.MustGet("user").(*Middleware.JWTUser).Privilege
    var canViewParameter = privilege >= Model.SiteAdmin ||
        fileModel.Package.UploaderID == c.MustGet("user").(*Middleware.JWTUser).UID
    if db.Preload("Package").First(&fileModel, uint(fid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    } else {
        FileCreateFetchURL(&fileModel, c, canViewParameter)
        Utility.MarshalResponse(c, 200, fileModel)
    }
    return nil
}

func FilePost(c *gin.Context) error {
    db := Database.GetDB()
    if db == nil { return Utility.ERR_DATABASE }
    var newFileModel Model.File
    if Utility.UnMarshalBody(c, &newFileModel) != nil {
        return Utility.ERR_BAD_PARAMETER
    }
    fid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var oldFileModel Model.File
    var privilege = c.MustGet("user").(*Middleware.JWTUser).Privilege
    if db.Preload("Package").First(&oldFileModel, uint(fid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    }
    newFileModel.ID = uint(fid)
    if privilege < Model.Moderator {
        if oldFileModel.Package.UploaderID != c.MustGet("user").(*Middleware.JWTUser).UID {
            return Utility.ERR_LOW_PRIV
        }
        newFileModel.PackageID = 0
    }
    return nil
}

func FileCreateFetchURL(f *Model.File, c *gin.Context, hideParameter bool) {
    countryCode := Database.GetIPCountryCode(c.ClientIP())
    print(countryCode)
    keyWithCountry := f.Service.String() + ":" + strings.ToLower(countryCode)
    var urlTemplate string
    if val, ok := Config.CurrentConfig.FileService.URLMap[keyWithCountry]; ok {
        urlTemplate = val
    } else {
        urlTemplate = Config.CurrentConfig.FileService.URLMap[f.Service.String()]
    }
    urlTemplate = strings.Replace(urlTemplate, "{FILE_ID}", strconv.Itoa(int(f.ID)), -1)
    urlTemplate = strings.Replace(urlTemplate, "{URL_PARAM}", f.URLParam.String, -1)
    urlTemplate = strings.Replace(urlTemplate, "{AUTH_PARAM}", f.AuthParam.String, -1)
    f.FetchURL = urlTemplate
    f.Service = nil
    if hideParameter {
        f.URLParam = Model.NullStringNull
        f.AuthParam = Model.NullStringNull
    }
}
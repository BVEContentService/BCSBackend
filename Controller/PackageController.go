package Controller

import (
    "OBPkg/Database"
    "OBPkg/Middleware"
    "OBPkg/Model"
    "OBPkg/Utility"
    "github.com/gin-gonic/gin"
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "time"
)

func PackageList(c *gin.Context) error {
    db := Database.GetDB()
    if db == nil { return Utility.ERR_DATABASE }
    var packages []Model.Package
    db.Find(&packages)
    Utility.MarshalResponse(c, 200, packages)
    return nil
}

func PackageGet(c *gin.Context) error {
    db := Database.GetDB()
    if db == nil { return Utility.ERR_DATABASE }
    pkgid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var packageModel Model.Package
    if db.Preload("Uploader").Preload("Files").First(&packageModel, uint(pkgid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    } else {
        for index := range packageModel.Files {
            FileCreateFetchURL(&packageModel.Files[index], c, true)
        }
        Utility.MarshalResponse(c, 200, packageModel)
    }
    return nil
}

func PackagePost(c *gin.Context) error {
    db := Database.GetDB();
    if db == nil { return Utility.ERR_DATABASE }
    var newPackageModel Model.Package
    if Utility.UnMarshalBody(c, &newPackageModel) != nil { return Utility.ERR_BAD_PARAMETER }
    pkgid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var oldPackageModel Model.Package
    var privilege = c.MustGet("user").(*Middleware.JWTUser).Privilege
    if db.First(&oldPackageModel, uint(pkgid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    }
    newPackageModel.ID = uint(pkgid)
    if privilege < Model.Moderator {
        if oldPackageModel.UploaderID != c.MustGet("user").(*Middleware.JWTUser).UID {
            return Utility.ERR_LOW_PRIV
        }
        newPackageModel.UploaderID = 0
        newPackageModel.GUID = Model.NullStringNull
        newPackageModel.Identifier = ""
    }
    cleanRequestPackage(&newPackageModel, privilege)
    if !validateRequestPackage(&newPackageModel, false) { return Utility.ERR_FORM_VALIDATE }
    db.Model(&oldPackageModel).Updates(&newPackageModel)
    if (!newPackageModel.IsRepost) {
        // Because the absence of Author field in the request body cannot tell if
        // the Author information is not to be changed or is to be removed
        // IsRepost internal field is used to indicate this operation.
        // If it is set to false, the Author information is to be removed.
        db.Model(&oldPackageModel).Updates(map[string]interface{}{
            "author_name_local":nil,
            "author_name_english":nil,
            "author_name_tag":nil,
            "author_email":nil,
            "author_homepage":nil,
        });
    }
    return PackageGet(c)
}

func PackagePut(c *gin.Context) error {
    db := Database.GetDB();
    if db == nil { return Utility.ERR_DATABASE }
    var newPackageModel Model.Package
    if Utility.UnMarshalBody(c, &newPackageModel) != nil { return Utility.ERR_BAD_PARAMETER }
    var privilege = c.MustGet("user").(*Middleware.JWTUser).Privilege
    newPackageModel.ID = 0
    if newPackageModel.UploaderID == 0 || privilege < Model.Moderator {
        newPackageModel.UploaderID = c.MustGet("user").(*Middleware.JWTUser).UID
    }
    cleanRequestPackage(&newPackageModel, privilege)
    if !validateRequestPackage(&newPackageModel, true) { return Utility.ERR_FORM_VALIDATE }
    var bufferPackageModel Model.Package
    if !db.Where("Identifier = ?", newPackageModel.Identifier).First(&bufferPackageModel).RecordNotFound() { return Utility.ERR_NAME_TAKEN }
    if !db.Where("GUID = ?", newPackageModel.GUID).First(&bufferPackageModel).RecordNotFound() { return Utility.ERR_GUID_TAKEN }
    db.Create(&newPackageModel)
    db.Preload("Uploader").Preload("Files").First(&newPackageModel, newPackageModel.ID)
    Utility.MarshalResponse(c, 201, newPackageModel)
    return nil
}

func PackageDelete(c *gin.Context) error {
    db := Database.GetDB();
    if db == nil { return Utility.ERR_DATABASE }
    pkgid, err := strconv.Atoi(c.Param("id"))
    if err != nil { return Utility.ERR_BAD_PARAMETER }
    var packageModel Model.Package
    if db.First(&packageModel, uint(pkgid)).RecordNotFound() {
        return Utility.ERR_DATA_NOT_FOUND
    }
    var privilege = c.MustGet("user").(*Middleware.JWTUser).Privilege
    if privilege < Model.Moderator {
        if packageModel.UploaderID != c.MustGet("user").(*Middleware.JWTUser).UID {
            return Utility.ERR_LOW_PRIV
        }
    }
    db.Delete(&packageModel)
    return Utility.SUCCESS_REMOVED
}

func cleanRequestPackage(newPackageModel *Model.Package, privilege Model.Privilege) {
    newPackageModel.Uploader = nil
    newPackageModel.Files = nil
    if privilege < Model.Moderator {
        newPackageModel.CreatedAt = time.Time{}
        newPackageModel.DeletedAt = nil
        newPackageModel.UpdatedAt = time.Time{}
    }
    newPackageModel.Homepage.TrimSpace()
    if newPackageModel.GUID.NotEmpty() {
        newPackageModel.GUID.String = strings.TrimSpace(strings.Replace(strings.ToLower(newPackageModel.GUID.String), "-", "", -1))
    }
    newPackageModel.Identifier = strings.TrimSpace(newPackageModel.Identifier)
    newPackageModel.Thumbnail.TrimSpace()
    newPackageModel.ThumbnailLQ.TrimSpace()
    newPackageModel.Name.TrimNames()
    if newPackageModel.Author != nil {
        newPackageModel.Author.Homepage.TrimSpace()
        newPackageModel.Author.Email.TrimSpace()
        newPackageModel.Author.Name.TrimNames()
        if *newPackageModel.Author == Model.NULL_DEVELOPER {
            newPackageModel.Author = nil
        }
    }
}

func validateRequestPackage(newPackageModel *Model.Package, isCreate bool) bool {
    emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    guidRegex := regexp.MustCompile("^[0-9a-f]{32}$")
    if newPackageModel.GUID.NotEmpty() && !guidRegex.MatchString(newPackageModel.GUID.String) { return false };
    if newPackageModel.Homepage.NotEmpty() && !isUrl(newPackageModel.Homepage.String) { return false }
    if isCreate && !newPackageModel.Name.Local.NotEmpty() { return false }
    if newPackageModel.ThumbnailLQ.NotEmpty() && !newPackageModel.Thumbnail.NotEmpty() { return false }
    if newPackageModel.Author != nil {
        if !newPackageModel.Author.Name.Local.NotEmpty() { return false }
        if newPackageModel.Author.Email.NotEmpty() && !emailRegex.MatchString(newPackageModel.Author.Email.String) { return false }
        if newPackageModel.Author.Homepage.NotEmpty() && !isUrl(newPackageModel.Author.Homepage.String) { return false }
    }
    return true;
}

func isUrl(str string) bool {
    u, err := url.Parse(str)
    return err == nil && u.Scheme != "" && u.Host != ""
}
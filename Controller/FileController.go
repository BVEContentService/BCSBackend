package Controller

import (
	"OBPkg/Config"
	"OBPkg/Database"
	"OBPkg/Middleware"
	"OBPkg/Model"
	"OBPkg/Utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func fileGetFullList(c *gin.Context, platform Model.PlatformType, validated *bool) []Model.File {
	var db = Database.GetDB()
	var requestUser, _ = c.Get("user")
	var returningModels, filteredModels []Model.File
	tx := db.Order("id desc")
	if platform != 0 {
		tx = tx.Where("platform = ?", platform)
	}
	tx.Find(&returningModels)
	canViewParameter := requestUser != nil && requestUser.(*Middleware.JWTUser).Privilege >= Model.SiteAdmin
	for index := range returningModels {
		fileCreateFetchURL(&returningModels[index], c, !canViewParameter)
	}
	if validated != nil {
		for _, model := range returningModels {
			var actuallyValidated = !model.NeedValidation || model.Validated
			if actuallyValidated == *validated {
				filteredModels = append(filteredModels, model)
			}
		}
	} else {
		filteredModels = returningModels
	}
	return filteredModels
}

func FileHead(c *gin.Context) error {
	var filterPlatform Model.PlatformType = 0
	var filterValidated *bool = nil
	if platform, ok := c.GetQuery("platform"); ok {
		ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[platform]
		if !ok {
			return Utility.ERR_BAD_PARAMETER
		}
		filterPlatform = Model.PlatformType(ptInt)
	}
	if validated, ok := c.GetQuery("validated"); ok {
		validBool, ok := strconv.ParseBool(validated)
		if ok != nil {
			return Utility.ERR_BAD_PARAMETER
		}
		filterValidated = &validBool
	}
	var returningModels = fileGetFullList(c, filterPlatform, filterValidated)
	c.Header("Accept-Ranges", "files")
	c.Header("Content-Range",
		fmt.Sprintf("files %d-%d/%d", 0, len(returningModels)-1, len(returningModels)))
	return nil
}

func FileList(c *gin.Context) error {
	var filterPlatform Model.PlatformType = 0
	var filterValidated *bool = nil
	if platform, ok := c.GetQuery("platform"); ok {
		ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[platform]
		if !ok {
			return Utility.ERR_BAD_PARAMETER
		}
		filterPlatform = Model.PlatformType(ptInt)
	}
	if validated, ok := c.GetQuery("validated"); ok {
		validBool, ok := strconv.ParseBool(validated)
		if ok != nil {
			return Utility.ERR_BAD_PARAMETER
		}
		filterValidated = &validBool
	}
	var returningModels = fileGetFullList(c, filterPlatform, filterValidated)
	begin, end, err := Utility.ParseRange(c, len(returningModels))
	if err != nil {
		return err
	}
	c.Header("Accept-Ranges", "files")
	c.Header("Content-Range", fmt.Sprintf("files %d-%d/%d", begin, end-1, len(returningModels)))
	Utility.MarshalResponse(c, 206, returningModels[begin:end])
	return nil
}

func FileGet(c *gin.Context) error {
	db := Database.GetDB()
	var returningModel Model.File
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if db.Preload("Package").First(&returningModel, uint(fileID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}

	var canViewParameter = requestUser.Privilege >= Model.SiteAdmin || returningModel.Package.UploaderID == requestUser.UID
	fileCreateFetchURL(&returningModel, c, !canViewParameter)
	returningModel.Package = nil

	Utility.MarshalResponse(c, 200, returningModel)
	return nil
}

func FilePost(c *gin.Context) error {
	db := Database.GetDB()
	var requestModel, originalModel, modifyingModel Model.File
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if db.First(&originalModel, uint(fileID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	modifyingModel = originalModel
	if Utility.UnMarshalBody(c, &requestModel) != nil ||
		Utility.UnMarshalBody(c, &modifyingModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	fileClean(&modifyingModel)
	modifyingModel.ID = uint(fileID)
	if requestUser.Privilege < Model.Moderator {
		if originalModel.Package.UploaderID != c.MustGet("user").(*Middleware.JWTUser).UID {
			return Utility.ERR_LOW_PRIV
		}
		modifyingModel.PackageID = originalModel.PackageID
		if Config.CurrentConfig.Platform.NeedValidation[modifyingModel.Platform.String()] {
			modifyingModel.Validated = false
		}
	}
	if err := fileValidate(&modifyingModel); err != "" {
		return Utility.ERR_FORM_VALIDATE.WithData(err)
	}
	db.Save(&modifyingModel)

	return FileGet(c)
}

func FilePut(c *gin.Context) error {
	var db = Database.GetDB()
	var creatingModel, checkModel Model.File
	var parentPackageModel Model.Package
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	if err := Utility.UnMarshalBody(c, &creatingModel); err != nil {
		return err
	}
	if db.First(&parentPackageModel, creatingModel.PackageID).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}

	fileClean(&creatingModel)
	creatingModel.ID = 0
	if err := fileValidate(&creatingModel); err != "" {
		return Utility.ERR_FORM_VALIDATE.WithData(err)
	}
	if parentPackageModel.UploaderID != c.MustGet("user").(*Middleware.JWTUser).UID &&
		requestUser.Privilege < Model.Moderator {
		return Utility.ERR_LOW_PRIV
	}
	if !db.Where(&Model.File{
		PackageID: creatingModel.PackageID,
		Platform:  creatingModel.Platform,
		Version:   creatingModel.Version,
	}).Find(&checkModel).RecordNotFound() {
		return Utility.ERR_PLAT_VER_TAKEN
	}
	db.Create(&creatingModel)

	db.First(&creatingModel, creatingModel.ID)
	Utility.MarshalResponse(c, 201, creatingModel)
	return nil
}

func FileDelete(c *gin.Context) error {
	var db = Database.GetDB()
	var removingModel Model.File
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	if db.Preload("Package").First(&removingModel, uint(fileID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	if requestUser.Privilege < Model.Moderator && removingModel.Package.UploaderID != requestUser.UID {
		return Utility.ERR_LOW_PRIV
	}
	db.Delete(&removingModel)

	return Utility.SUCCESS_NO_BODY
}

func fileClean(newFileModel *Model.File) {
	newFileModel.Package = nil
	newFileModel.FetchURL = ""
}

func fileValidate(newFileModel *Model.File) string {
	if newFileModel.Platform.String() == "" {
		return "Platform"
	}
	if newFileModel.Service.String() == "" {
		return "Service"
	}
	var versionParts = strings.Split(newFileModel.Version, ".")
	if len(versionParts) < 2 || len(versionParts) > 4 {
		return "Version:Part"
	}
	for _, v := range versionParts {
		if v != strings.TrimSpace(v) {
			return "Version:Trim"
		}
		if _, err := strconv.Atoi(v); err != nil {
			return "Version:NAN"
		}
	}
	var sizeParts = strings.Split(newFileModel.Size, " ")
	if len(sizeParts) != 2 {
		return "Size:Part"
	}
	if _, err := strconv.ParseFloat(sizeParts[0], 64); err != nil {
		return "Size:NAN"
	}
	switch sizeParts[1] {
	case "B", "KB", "MB", "GB", "TB":
		break
	default:
		return "Size:Unit"
	}
	return ""
}

func fileCreateFetchURL(f *Model.File, c *gin.Context, hideParameter bool) {
	var countryCode string
	if ccode, ok := c.Get("countryCode"); ok {
		countryCode = ccode.(string)
	} else {
		countryCode = Database.GetIPCountryCode(c.ClientIP())
		c.Set("countryCode", countryCode)
	}
	keyWithCountry := f.Service.String() + ":" + strings.ToLower(countryCode)
	var urlTemplate string
	if val, ok := Config.CurrentConfig.FileService.URLMap[keyWithCountry]; ok {
		urlTemplate = val
	} else {
		urlTemplate = Config.CurrentConfig.FileService.URLMap[f.Service.String()]
	}
	urlTemplate = strings.Replace(urlTemplate, "{FILE_ID}", strconv.Itoa(int(f.ID)), -1)
	urlTemplate = strings.Replace(urlTemplate, "{URL_PARAM}", f.URLParam, -1)
	urlTemplate = strings.Replace(urlTemplate, "{AUTH_PARAM}", f.AuthParam, -1)
	urlTemplate = strings.Replace(urlTemplate, "{EMAIL_DOT}", strings.Replace(f.AuthParam, "@", ".", -1), -1)
	f.FetchURL = urlTemplate
	if hideParameter {
		f.Service = 0
		f.URLParam = ""
		f.AuthParam = ""
	}
}

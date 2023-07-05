package Controller

import (
	"OBPkg/Config"
	"OBPkg/Database"
	"OBPkg/Middleware"
	"OBPkg/Model"
	"OBPkg/Utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func packageGetFullList(platform Model.PlatformType, keyword string) []Model.Package {
	var db = Database.GetDB()
	var returningModels []Model.Package
	db.Preload("Files").Order("id desc").Find(&returningModels)
	for index := range returningModels {
		packageCollectPlatforms(&returningModels[index], false)
		// Reduce response size
		returningModels[index].Description = ""
		returningModels[index].Files = nil
	}
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if platform != 0 || keyword != "" {
		var filteredModels []Model.Package
		for _, model := range returningModels {
			if model.MatchPlatform(platform) && model.MatchKeyword(keyword) {
				filteredModels = append(filteredModels, model)
			}
		}
		return filteredModels
	} else {
		return returningModels
	}
}

func PackageHead(c *gin.Context) error {
	var filterPlatform Model.PlatformType = 0
	var filterKeyword string = ""
	if platform, ok := c.GetQuery("platform"); ok {
		if platform == "+" {
			filterPlatform = Model.PlatformType(-1)
		} else {
			ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[platform]
			if !ok {
				return Utility.ERR_BAD_PARAMETER
			}
			filterPlatform = Model.PlatformType(ptInt)
		}
	}
	if keyword, ok := c.GetQuery("keyword"); ok {
		filterKeyword = keyword
	}
	var returningModels = packageGetFullList(filterPlatform, filterKeyword)
	c.Header("Accept-Ranges", "packages")
	c.Header("Content-Range",
		fmt.Sprintf("packages %d-%d/%d", 0, len(returningModels)-1, len(returningModels)))
	return nil
}

func PackageList(c *gin.Context) error {
	var filterPlatform Model.PlatformType = 0
	var filterKeyword string = ""
	if platform, ok := c.GetQuery("platform"); ok {
		if platform == "+" {
			filterPlatform = Model.PlatformType(-1)
		} else {
			ptInt, ok := Config.CurrentConfig.Platform.DatabaseMap[platform]
			if !ok {
				return Utility.ERR_BAD_PARAMETER
			}
			filterPlatform = Model.PlatformType(ptInt)
		}
	}
	if keyword, ok := c.GetQuery("keyword"); ok {
		filterKeyword = keyword
	}
	var returningModels = packageGetFullList(filterPlatform, filterKeyword)
	begin, end, err := Utility.ParseRange(c, len(returningModels))
	if err != nil {
		return err
	}
	c.Header("Accept-Ranges", "packages")
	c.Header("Content-Range", fmt.Sprintf("packages %d-%d/%d", begin, end-1, len(returningModels)))
	if len(returningModels) == 0 {
		Utility.MarshalResponse(c, 206, []Model.Package{})
	} else {
		Utility.MarshalResponse(c, 206, returningModels[begin:end])
	}
	return nil
}

func queryPackage1b1(db *gorm.DB, id string) *Model.Package {
	var returningModel Model.Package
	idAsGUID := Utility.NormalizeGUID(id)
	idAsNumber, numConvErr := strconv.Atoi(id)
	if len(idAsGUID) == 32 {
		if !db.Preload("Uploader").Preload("Files").Where("guid = ?", idAsGUID).First(&returningModel).RecordNotFound() {
			return &returningModel
		}
	} else if numConvErr == nil {
		if !db.Preload("Uploader").Preload("Files").First(&returningModel, uint(idAsNumber)).RecordNotFound() {
			return &returningModel
		}
	} else {
		print(id)
		if !db.Preload("Uploader").Preload("Files").Where("identifier = ?", id).First(&returningModel).RecordNotFound() {
			return &returningModel
		}
	}
	return nil
}

func PackageGet(c *gin.Context) error {
	var db = Database.GetDB()
	var requestUser, _ = c.Get("user")
	returningModel := queryPackage1b1(db, c.Param("id"))
	if returningModel == nil {
		return Utility.ERR_DATA_NOT_FOUND
	}
	var canViewParameter = requestUser != nil &&
		(requestUser.(*Middleware.JWTUser).Privilege >= Model.Moderator ||
			returningModel.UploaderID == requestUser.(*Middleware.JWTUser).UID)
	for index := range returningModel.Files {
		fileCreateFetchURL(&returningModel.Files[index], c, !canViewParameter)
	}
	packageCollectPlatforms(returningModel, false)
	if returningModel.Author != nil && *returningModel.Author == Model.NULL_DEVELOPER {
		returningModel.Author = nil
	}
	Utility.MarshalResponse(c, 200, returningModel)
	return nil
}

func PackagePost(c *gin.Context) error {
	var db = Database.GetDB()
	var requestModel, originalModel, modifyingModel, checkModel Model.Package
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	packageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if db.First(&originalModel, uint(packageID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	modifyingModel = originalModel
	if Utility.UnMarshalBody(c, &requestModel) != nil ||
		Utility.UnMarshalBody(c, &modifyingModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	packageClean(&modifyingModel)
	modifyingModel.ID = uint(packageID)
	if requestModel.Author != nil && requestModel.Author.Name.Local == "" {
		modifyingModel.Author = &Model.Developer{}
	}
	if requestUser.Privilege < Model.Moderator {
		if originalModel.UploaderID != requestUser.UID {
			return Utility.ERR_LOW_PRIV
		}
		modifyingModel.UploaderID = originalModel.UploaderID
		modifyingModel.Identifier = originalModel.Identifier
		modifyingModel.CreatedAt = originalModel.CreatedAt
		modifyingModel.UpdatedAt = originalModel.UpdatedAt
		modifyingModel.DeletedAt = originalModel.DeletedAt
	}
	if err := packageValidate(&modifyingModel); err != "" {
		return Utility.ERR_FORM_VALIDATE.WithData(err)
	}
	if modifyingModel.GUID != "" {
		if !db.Where("GUID = ?", modifyingModel.GUID).First(&checkModel).RecordNotFound() {
			return Utility.ERR_GUID_TAKEN
		}
	}
	db.Save(&modifyingModel)

	return PackageGet(c)
}

func PackagePut(c *gin.Context) error {
	var db = Database.GetDB()
	var creatingModel, checkModel Model.Package
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	if Utility.UnMarshalBody(c, &creatingModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	packageClean(&creatingModel)
	creatingModel.ID = 0
	if creatingModel.UploaderID == 0 || requestUser.Privilege < Model.Moderator {
		creatingModel.UploaderID = requestUser.UID
	}
	if err := packageValidate(&creatingModel); err != "" {
		return Utility.ERR_FORM_VALIDATE.WithData(err)
	}
	if !db.Where("Identifier = ?", creatingModel.Identifier).First(&checkModel).RecordNotFound() {
		return Utility.ERR_NAME_TAKEN
	}
	if creatingModel.GUID != "" {
		if !db.Where("GUID = ?", creatingModel.GUID).First(&checkModel).RecordNotFound() {
			return Utility.ERR_GUID_TAKEN
		}
	}
	db.Create(&creatingModel)

	// This seems to cause null pointer dereference, I guess maybe there's some kind of cache
	//db.Preload("Uploader").Preload("Files").First(&creatingModel, creatingModel.ID)
	//if *creatingModel.Author == Model.NULL_DEVELOPER {
	//	creatingModel.Author = nil
	//}
	Utility.MarshalResponse(c, 201, creatingModel)
	return nil
}

func PackageDelete(c *gin.Context) error {
	var db = Database.GetDB()
	var removingModel Model.Package
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	packageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	if db.First(&removingModel, uint(packageID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	if requestUser.Privilege < Model.Moderator && removingModel.UploaderID != requestUser.UID {
		return Utility.ERR_LOW_PRIV
	}
	db.Delete(&removingModel)

	return Utility.SUCCESS_NO_BODY
}

func packageClean(newPackageModel *Model.Package) {
	newPackageModel.Uploader = nil
	newPackageModel.Files = nil
	if newPackageModel.GUID != "" {
		newPackageModel.GUID = Utility.NormalizeGUID(newPackageModel.GUID)
	}
	newPackageModel.Identifier = strings.TrimSpace(newPackageModel.Identifier)
	newPackageModel.Name.TrimNames()
	if newPackageModel.Author != nil {
		newPackageModel.Author.Name.TrimNames()
		if *newPackageModel.Author == Model.NULL_DEVELOPER {
			newPackageModel.Author = nil
		}
	}
}

func packageValidate(newPackageModel *Model.Package) string {
	if newPackageModel.GUID != "" && !Utility.REGEX_GUID.MatchString(newPackageModel.GUID) {
		return "GUID"
	}
	if newPackageModel.Homepage != "" && !Utility.IsUrl(newPackageModel.Homepage) {
		return "Homepage"
	}
	if newPackageModel.Name.Local == "" {
		return "Name.Local"
	}
	if newPackageModel.ThumbnailLQ != "" && newPackageModel.Thumbnail == "" {
		return "ThumbLQ-Thumb"
	}
	if newPackageModel.Thumbnail != "" && !Utility.IsUrl(newPackageModel.Thumbnail) {
		return "Thumb"
	}
	if newPackageModel.ThumbnailLQ != "" && !Utility.IsUrl(newPackageModel.ThumbnailLQ) {
		return "ThumbLQ"
	}
	if newPackageModel.Author != nil {
		if newPackageModel.Author.Name.Local == "" {
			return "Author.Name.Local"
		}
		if newPackageModel.Author.Email != "" &&
			!Utility.REGEX_EMAIL.MatchString(newPackageModel.Author.Email) {
			return "Author.Email"
		}
		if newPackageModel.Author.Homepage != "" &&
			!Utility.IsUrl(newPackageModel.Author.Homepage) {
			return "Author.Homepage"
		}
	}
	return ""
}

func packageCollectPlatforms(packageModel *Model.Package, collectUnvalidated bool) {
	var files []Model.File
	var platforms []Model.PlatformType
	if packageModel.Files == nil {
		var db = Database.GetDB()
		db.Model(&packageModel).Association("Files").Find(&files)
	} else {
		files = packageModel.Files
	}
	for _, file := range files {
		if Config.CurrentConfig.Platform.NeedValidation[file.Platform.String()] &&
			!file.Validated && !collectUnvalidated {
			continue
		}
		if !inSlice(platforms, file.Platform) {
			platforms = append(platforms, file.Platform)
		}
	}
	packageModel.Platforms = platforms
}

func inSlice(haystack []Model.PlatformType, needle Model.PlatformType) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}
	return false
}

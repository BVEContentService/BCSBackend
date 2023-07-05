package Controller

import (
	"OBPkg/Database"
	"OBPkg/Middleware"
	"OBPkg/Model"
	"OBPkg/Utility"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UploaderGet(c *gin.Context) error {
	db := Database.GetDB()
	upid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	var uploaderModel Model.Uploader
	if db.Preload("Packages").First(&uploaderModel, uint(upid)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	for index := range uploaderModel.Packages {
		packageCollectPlatforms(&uploaderModel.Packages[index], false)
		// Reduce Response Size
		uploaderModel.Packages[index].Description = ""
	}
	Utility.MarshalResponse(c, 200, uploaderModel)
	return nil
}

func UploaderPost(c *gin.Context) error {
	var db = Database.GetDB()
	var originalModel, modifyingModel Model.Uploader
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	uploaderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if db.First(&originalModel, uint(uploaderID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	modifyingModel = originalModel
	if Utility.UnMarshalBody(c, &modifyingModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}

	modifyingModel.ID = uint(uploaderID)
	//modifyingModel.Username = originalModel.Username
	modifyingModel.Password = originalModel.Password
	if requestUser.Privilege < Model.Moderator {
		modifyingModel.Email = originalModel.Email
		modifyingModel.Validated = originalModel.Validated
		modifyingModel.Privilege = originalModel.Privilege
		modifyingModel.CreatedAt = originalModel.CreatedAt
		modifyingModel.UpdatedAt = originalModel.UpdatedAt
		modifyingModel.DeletedAt = originalModel.DeletedAt
	}
	if err := uploaderValidate(&modifyingModel); err != "" {
		return Utility.ERR_FORM_VALIDATE.WithData(err)
	}
	db.Save(&modifyingModel)

	return UploaderGet(c)
}

func uploaderValidate(newModel *Model.Uploader) string {
	if newModel.Email != "" && !Utility.REGEX_EMAIL.MatchString(newModel.Email) {
		return "Email"
	}
	if newModel.Homepage != "" && !Utility.IsUrl(newModel.Homepage) {
		return "Homepage"
	}
	if newModel.Name.Local == "" {
		return "Name.Local"
	}
	//if newModel.Username == "" {
	//	return "Username"
	//}
	switch newModel.Privilege {
	case Model.Normal, Model.Validator, Model.Moderator, Model.SiteAdmin:
		break
	default:
		return "Privilege"
	}
	return ""
}

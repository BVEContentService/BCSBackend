package Controller

import (
	"OBPkg/Config"
	"OBPkg/Database"
	"OBPkg/Middleware"
	"OBPkg/Model"
	"OBPkg/Utility"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"time"
)

type activateRequestModel struct {
	Token string
	Name  Model.String3
	//Username string
	Password string
}

func AuthChangePassword(c *gin.Context) error {
	var db = Database.GetDB()
	var modifyingModel Model.Uploader
	var requestUser = c.MustGet("user").(*Middleware.JWTUser)
	uploaderID := requestUser.UID
	if db.First(&modifyingModel, uint(uploaderID)).RecordNotFound() {
		return Utility.ERR_DATA_NOT_FOUND
	}
	if modifyingModel.ID != requestUser.UID && requestUser.Privilege < Model.SiteAdmin {
		return Utility.ERR_LOW_PRIV
	}
	var changeRequest struct {
		PreviousPassword string
		NewPassword      string
	}
	if Utility.UnMarshalBody(c, &changeRequest) != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if !Utility.BCryptValidateHash(changeRequest.PreviousPassword, modifyingModel.Password) &&
		requestUser.Privilege < Model.SiteAdmin {
		return Utility.ERR_BAD_CERT
	}
	if changeRequest.NewPassword == "" {
		modifyingModel.Password = ""
	} else {
		modifyingModel.Password = Utility.BCryptCalculateHash(changeRequest.NewPassword)
	}
	db.Save(&modifyingModel)

	return Utility.SUCCESS_NO_BODY
}

func AuthRegister(c *gin.Context) error {
	// The "Register" process only files a request
	var db = Database.GetDB()
	var requestModel, requestCheckModel Model.RegisterRequest
	var uploaderCheckModel Model.Uploader
	if Utility.UnMarshalBody(c, &requestModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if requestModel.Email == "" {
		return Utility.ERR_BAD_PARAMETER
	}
	if !db.Where("email = ?", requestModel.Email).First(&uploaderCheckModel).RecordNotFound() {
		return Utility.ERR_EMAIL_TAKEN
	}
	requestExist := false
	if !db.Where("email = ?", requestModel.Email).First(&requestCheckModel).RecordNotFound() {
		requestExist = true
		if time.Now().Before(requestCheckModel.Expiry) {
			return Utility.ERR_EMAIL_TMR_EARLY.WithData(
				requestCheckModel.Expiry.UTC().Format("2006-01-02 15:00 UTC"))
		}
	} else {
		requestCheckModel.Email = requestModel.Email
	}
	if err := checkmail.ValidateFormat(requestModel.Email); err != nil {
		errorTemplate := Utility.ERR_BAD_PARAMETER
		errorTemplate.Data = err.Error()
		return errorTemplate
	}
	requestCheckModel.Token = shortuuid.New()
	requestCheckModel.Expiry = time.Now().Add(time.Duration(Config.CurrentConfig.SMTP.TokenDuration))
	requestCheckModel.Affair = Model.Register
	if err := Utility.EmailSendConfirmation(requestCheckModel.Email, requestCheckModel.Token,
		requestCheckModel.Expiry.UTC().Format("2006-01-02 15:00 UTC")); err != nil {
		return err
	}
	if requestExist {
		db.Save(&requestCheckModel)
	} else {
		requestCheckModel.ID = 0
		requestCheckModel.Email = requestModel.Email
		db.Create(&requestCheckModel)
	}
	Utility.MarshalResponse(c, 200, requestCheckModel)
	return nil
}

func AuthCheckToken(c *gin.Context) error {
	// The "Activate" process actually creates the account
	var db = Database.GetDB()
	var requestModel activateRequestModel
	if Utility.UnMarshalBody(c, &requestModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	if requestModel.Token == "" {
		return Utility.ERR_BAD_PARAMETER
	}
	var requestCheckModel Model.RegisterRequest
	if !db.Where("token = ?", requestModel.Token).First(&requestCheckModel).RecordNotFound() {
		if time.Now().After(requestCheckModel.Expiry) {
			return Utility.ERR_EMAIL_TMR_LATE
		}
	} else {
		return Utility.ERR_EMAIL_TMR_LATE
	}
	Utility.MarshalResponse(c, 200, requestCheckModel)
	//return Utility.SUCCESS_NO_BODY
	return nil
}

func AuthActivate(c *gin.Context) error {
	// The "Activate" process actually creates the account
	var db = Database.GetDB()
	var requestCheckModel Model.RegisterRequest
	var requestModel activateRequestModel
	var uploaderCheckModel Model.Uploader
	if Utility.UnMarshalBody(c, &requestModel) != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	//if requestModel.Token == "" || requestModel.Username == "" || requestModel.Password == "" ||
	if requestModel.Token == "" || requestModel.Password == "" ||
		requestModel.Name.Local == "" {
		return Utility.ERR_BAD_PARAMETER
	}
	//if requestModel.Username == "" {
	//	return Utility.ERR_BAD_PARAMETER
	//}
	if !db.Where("token = ?", requestModel.Token).First(&requestCheckModel).RecordNotFound() {
		if time.Now().After(requestCheckModel.Expiry) {
			return Utility.ERR_EMAIL_TMR_LATE
		}
	} else {
		return Utility.ERR_EMAIL_TMR_LATE
	}
	if checkmail.ValidateFormat(requestCheckModel.Email) != nil {
		return Utility.ERR_BAD_PARAMETER
	}
	//if !regexp.MustCompile("^[a-z][a-z0-9-_]*$").MatchString(requestModel.Username) {
	//	return Utility.ERR_BAD_PARAMETER
	//}
	if !db.Where("email = ?", requestCheckModel.Email).First(&uploaderCheckModel).RecordNotFound() {
		return Utility.ERR_EMAIL_TAKEN
	}
	//if !db.Where("username = ?", requestModel.Username).First(&uploaderCheckModel).RecordNotFound() {
	//	return Utility.ERR_NAME_TAKEN
	//}
	uploaderCheckModel = Model.Uploader{
		Developer: Model.Developer{
			Name:  requestModel.Name,
			Email: requestCheckModel.Email,
		},
		Validated: false,
		//Username:    requestModel.Username,
		Password:    Utility.BCryptCalculateHash(requestModel.Password),
		Description: "",
		Privilege:   Model.Normal,
	}
	db.Create(&uploaderCheckModel)
	db.Delete(&requestCheckModel)
	Utility.MarshalResponse(c, 200, uploaderCheckModel)
	return nil
}

func AuthEmailAffair(c *gin.Context) error {
	var db = Database.GetDB()
	var requestCheckModel Model.RegisterRequest
	if !db.Where("token = ?", c.Param("token")).First(&requestCheckModel).RecordNotFound() {
		if time.Now().After(requestCheckModel.Expiry) {
			return Utility.ERR_EMAIL_TMR_LATE
		}
	} else {
		return Utility.ERR_EMAIL_TMR_LATE
	}
	if requestCheckModel.Affair == Model.ResetPassword {
		var modifyingModel Model.Uploader
		uploaderID := requestCheckModel.UserID
		if db.First(&modifyingModel, uint(uploaderID)).RecordNotFound() {
			return Utility.ERR_DATA_NOT_FOUND
		}
		modifyingModel.Password = Utility.BCryptCalculateHash(requestCheckModel.Email)
		db.Save(&modifyingModel)
		db.Delete(&requestCheckModel)
		return Utility.SUCCESS
	} else if requestCheckModel.Affair == Model.ChangeEmail {
		var uploaderCheckModel Model.Uploader
		if checkmail.ValidateFormat(requestCheckModel.Email) != nil {
			return Utility.ERR_BAD_PARAMETER
		}
		if !db.Where("email = ?", requestCheckModel.Email).First(&uploaderCheckModel).RecordNotFound() {
			return Utility.ERR_EMAIL_TAKEN
		}
		var modifyingModel Model.Uploader
		uploaderID := requestCheckModel.UserID
		if db.First(&modifyingModel, uint(uploaderID)).RecordNotFound() {
			return Utility.ERR_DATA_NOT_FOUND
		}
		modifyingModel.Email = requestCheckModel.Email
		db.Save(&modifyingModel)
		db.Delete(&requestCheckModel)
		return Utility.SUCCESS
	}
	return Utility.ERR_BAD_PARAMETER
}

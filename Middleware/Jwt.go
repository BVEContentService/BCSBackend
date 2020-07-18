package Middleware

import (
	"OBPkg/Config"
	"OBPkg/Database"
	"OBPkg/Model"
	"OBPkg/Utility"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password"`
}

type JWTUser struct {
	UID       uint
	Privilege Model.Privilege
}

var AuthMiddleware *jwt.GinJWTMiddleware

func GetAuthMiddleware() *jwt.GinJWTMiddleware {
	if AuthMiddleware != nil {
		return AuthMiddleware
	}
	db := Database.GetDB()
	if db == nil {
		panic("Authorization Setup Failed")
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "obpkg",
		Key:         []byte(Config.CurrentConfig.JWT.SecretKey),
		Timeout:     time.Duration(Config.CurrentConfig.JWT.Timeout),
		MaxRefresh:  time.Duration(Config.CurrentConfig.JWT.MaxRefresh),
		IdentityKey: "uid",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*JWTUser); ok {
				return jwt.MapClaims{
					"uid":       v.UID,
					"privilege": v.Privilege,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &JWTUser{
				UID:       uint(claims["uid"].(float64)),
				Privilege: Model.Privilege(claims["privilege"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", Utility.ERR_BAD_PARAMETER
			}
			var u Model.Uploader
			if db.Where("username = ?", loginVals.Username).First(&u).RecordNotFound() ||
				!Utility.BCryptValidateHash(loginVals.Password, u.Password) {
				return nil, Utility.ERR_BAD_CERT
			} else {
				return &JWTUser{UID: u.ID, Privilege: u.Privilege}, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			v, ok := data.(*JWTUser)
			if !ok {
				return false
			}
			c.Set("user", v)
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if message == Utility.ERR_BAD_CERT.Msg {
				Utility.MarshalResponse(c, code, Utility.ERR_BAD_CERT)
			} else {
				Utility.MarshalResponse(c, code, Utility.JWTError(code, message))
			}
		},
	})
	if err != nil {
		panic(err)
	}
	AuthMiddleware = authMiddleware
	return authMiddleware
}

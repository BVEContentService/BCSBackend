package Utility

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
)

func MarshalResponse(c *gin.Context, code int, content interface{}) {
	var accept = strings.Split(c.Request.Header.Get("Accept"), ",")
	for _, contentType := range accept {
		if contentType == "application/xml" {
			c.XML(code, content)
			return
		} else if contentType == "application/json" {
			c.JSON(code, content)
			return
		}
	}
	c.XML(code, content)
}

func UnMarshalBody(c *gin.Context, out interface{}) error {
	var contentType = strings.TrimSpace(strings.Split(c.Request.Header.Get("Content-Type"), ";")[0])
	if contentType == "application/json" {
		return c.ShouldBindBodyWith(&out, binding.JSON)
	} else if contentType == "application/xml" || contentType == "" {
		return c.ShouldBindBodyWith(&out, binding.XML)
	} else {
		return ERR_BAD_PARAMETER
	}
}

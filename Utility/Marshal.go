package Utility

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"reflect"
	"strings"
)

func MarshalResponse(c *gin.Context, code int, content interface{}) {
	var accept = strings.Split(c.Request.Header.Get("Accept"), ",")
	for _, contentType := range accept {
		if contentType == "application/json" {
			c.JSON(code, content)
			return
		}
	}
	if reflect.TypeOf(content).Kind() == reflect.Slice {
		c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Writer.WriteHeader(code)
		_, _ = c.Writer.WriteString("<Slice>")
		c.XML(code, content)
		_, _ = c.Writer.WriteString("</Slice>")
	} else {
		c.XML(code, content)
	}
}

func UnMarshalBody(c *gin.Context, out interface{}) error {
	var contentType = strings.TrimSpace(strings.Split(c.Request.Header.Get("Content-Type"), ";")[0])
	if contentType == "application/json" {
		err := c.ShouldBindBodyWith(&out, binding.JSON)
		return err
	} else if contentType == "application/xml" || contentType == "" {
		err := c.ShouldBindBodyWith(&out, binding.XML)
		return err
	} else {
		return ERR_BAD_PARAMETER
	}
}

package Utility

import (
    "encoding/json"
    "encoding/xml"
    "github.com/gin-gonic/gin"
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
    var rawBody []byte
    var err error
    if rawBody, err = c.GetRawData(); err != nil { return ERR_BAD_PARAMETER }
    if contentType == "application/json" {
        if json.Unmarshal(rawBody, &out) != nil { return ERR_BAD_PARAMETER }
    } else if contentType == "application/xml" || contentType == "" {
        if xml.Unmarshal(rawBody, &out) != nil { return ERR_BAD_PARAMETER }
    } else {
        return ERR_BAD_PARAMETER
    }
    return nil
}
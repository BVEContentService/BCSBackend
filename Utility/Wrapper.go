package Utility

import "github.com/gin-gonic/gin"

type HandlerFunc func(c *gin.Context) error

func (e *Error) PrintResponse(c *gin.Context){
    MarshalResponse(c, e.Code, e)
}

func Wrapper(handler HandlerFunc) func(c *gin.Context){
    return func(c *gin.Context){
        var err = handler(c)
        if err != nil {
            var apiException *Error
            if h, ok := err.(*Error); ok {
                apiException = h
            } else if e, ok := err.(error); ok {
                apiException = UnknownError(e.Error())
            } else {
                apiException = UnknownError("Unknown Failure")
            }
            apiException.Request = c.Request.Method + " " + c.Request.URL.String()
            apiException.PrintResponse(c)
            return
        }
    }
}

package Utility

type Error struct {
    Code        int     `xml:"-" json:"-"`
    ErrorCode   int
    Msg         string
    Request     string
}

var (
    SUCCESS_REMOVED     = &Error { Code: 204, ErrorCode: 001, Msg: "Resource Removed"}
    ERR_NOT_FOUND_URL   = &Error { Code: 404, ErrorCode: 101, Msg: "API Endpoint not Found" }
    ERR_NOT_ALLOWED     = &Error { Code: 405, ErrorCode: 102, Msg: "Method not Allowed" }
    ERR_NO_PARAMETER    = &Error { Code: 400, ErrorCode: 111, Msg: "Required Parameter not Given" }
    ERR_BAD_PARAMETER   = &Error { Code: 400, ErrorCode: 112, Msg: "Parameter not in Correct Format" }
    ERR_FORM_VALIDATE   = &Error { Code: 422, ErrorCode: 113, Msg: "Request Data Failed Validation" }
    ERR_DATA_NOT_FOUND  = &Error { Code: 404, ErrorCode: 201, Msg: "Requested Resource not Found" }
    ERR_NAME_TAKEN      = &Error { Code: 422, ErrorCode: 202, Msg: "Identifier Name Already Taken"}
    ERR_GUID_TAKEN      = &Error { Code: 422, ErrorCode: 203, Msg: "Identifier GUID Already Taken"}
    ERR_JWT             = &Error { Code: 401, ErrorCode: 211, Msg: "Authorization Problem: " }
    ERR_BAD_CERT        = &Error { Code: 403, ErrorCode: 212, Msg: "Incorrect Username or Password" }
    ERR_LOW_PRIV        = &Error { Code: 403, ErrorCode: 221, Msg: "Manipulating Resources of Others is not Allowed"}
    ERR_DATABASE        = &Error { Code: 500, ErrorCode: 302, Msg: "Error during Database Operation" }
)

func (e *Error) Error() string {
    return  e.Msg
}

func UnknownError(msg string) *Error {
    return &Error { Code: 500, ErrorCode: 301, Msg: msg }
}

func JWTError(code int, msg string) *Error {
    return &Error { Code: code, ErrorCode: ERR_JWT.ErrorCode, Msg: ERR_JWT.Msg + msg }
}
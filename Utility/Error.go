package Utility

type Error struct {
	Code      int `xml:"-" json:"-"`
	ErrorCode int
	Msg       string
	Data      string `xml:",omitempty" json:",omitempty"`
	Request   string
}

var (
	SUCCESS_NO_BODY     = &Error{Code: 204, ErrorCode: 001, Msg: ""}
	ERR_NOT_FOUND_URL   = &Error{Code: 404, ErrorCode: 101, Msg: "API Endpoint not Found"}
	ERR_NOT_ALLOWED     = &Error{Code: 405, ErrorCode: 102, Msg: "Method not Allowed"}
	ERR_BAD_PARAMETER   = &Error{Code: 400, ErrorCode: 111, Msg: "Parameter not in Correct Format"}
	ERR_FORM_VALIDATE   = &Error{Code: 422, ErrorCode: 112, Msg: "Request Data Failed Validation"}
	ERR_BAD_RANGE       = &Error{Code: 400, ErrorCode: 113, Msg: "Range Header Syntax is Invalid"}
	ERR_NONEXIST_RANGE  = &Error{Code: 416, ErrorCode: 114, Msg: "Requested Range not Satisfiable"}
	ERR_DATA_NOT_FOUND  = &Error{Code: 404, ErrorCode: 201, Msg: "Requested Resource not Found"}
	ERR_NAME_TAKEN      = &Error{Code: 422, ErrorCode: 202, Msg: "Identifier Name already Taken"}
	ERR_GUID_TAKEN      = &Error{Code: 422, ErrorCode: 203, Msg: "Identifier GUID already Taken"}
	ERR_PLAT_VER_TAKEN  = &Error{Code: 422, ErrorCode: 204, Msg: "Already Registered a File with the Same Platform and Version"}
	ERR_EMAIL_TAKEN     = &Error{Code: 422, ErrorCode: 205, Msg: "This Email has already been Registered"}
	ERR_JWT             = &Error{Code: 401, ErrorCode: 211, Msg: "Authorization Problem: "}
	ERR_BAD_CERT        = &Error{Code: 403, ErrorCode: 212, Msg: "Incorrect Username or Password"}
	ERR_EMAIL_TMR_EARLY = &Error{Code: 429, ErrorCode: 213, Msg: "Cannot Send next Email before this Time:"}
	ERR_EMAIL_TMR_LATE  = &Error{Code: 401, ErrorCode: 214, Msg: "The Token has already Expired"}
	ERR_LOW_PRIV        = &Error{Code: 403, ErrorCode: 221, Msg: "Manipulating Resources of Others is not Allowed"}
)

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) WithData(data string) *Error {
	return &Error{e.Code, e.ErrorCode, e.Msg, data, e.Request}
}

func UnknownError(msg string) *Error {
	return &Error{Code: 500, ErrorCode: 301, Msg: msg}
}

func JWTError(code int, msg string) *Error {
	return &Error{Code: code, ErrorCode: ERR_JWT.ErrorCode, Msg: ERR_JWT.Msg + msg}
}

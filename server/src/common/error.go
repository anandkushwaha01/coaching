package common
import(
	"errors"
)
const(
	ERROR_TYPE_DB int = 700+iota
	ERROR_TYPE_OAUTH int =751
	ERROR_TYPE_COMMON int = 851
	ERROR_TYPE_CLIENT int = 901
	ERROR_TYPE_SESSION int = 1001
)
const(
	DB_ERROR_NO_CONNECTION int = ERROR_TYPE_DB + iota
	DB_ERROR_STMT_PREP
	DB_ERROR_STMT_EXEC
	DB_ERROR_QUERY_EXEC
)

const(
	OAUTH_ERROR_NO_USER int = iota+ERROR_TYPE_OAUTH
	OAUTH_ERROR_USER_EXIST
	OAUTH_ERROR_NO_STAGE_DATA
	OAUTH_ERROR_SECRET_NULL
	OAUTH_ERROR_NULL_USER_OR_PASS
	OAUTH_ERROR_AUTHENTICATION_FAILED
	OAUTH_ERROR_USERNAME_PASS_MISMATCH
	OAUTH_ERROR_INVALID_SIGNUPDATA
	OAUTH_ERROR_EMAIL_TYPE_OR_SECRET_NULL
	OAUTH_ERROR_IN_SECRET_GEN
)
const(
	ERROR_INVALID_REDIRECT_URI int = ERROR_TYPE_COMMON + iota
	ERROR_INVALID_CONTENT_TYPE
	ERROR_INVALID_METHOD_TYPE
	ERROR_INVALID_PARAM
	ERROR_INVALID_EMAIL
	ERROR_INVALID_NAME
	ERROR_INVALID_PHONE_NO
	ERROR_INVALID_LINK
	ERROR_INVALID_VERIFICATION_LINK
	ERROR_JSON_BODY_NOT_READABLE
	ERROR_JSON_BODY_NOT_PARSABLE
	ERROR_TEMPLATE_IS_NULL
)

var ErrorMsgMap map[int]string

func GetErrorMessage(err_code int) (err_msg string, err error){
	var ok bool
	if err_msg, ok = ErrorMsgMap[err_code]; !ok{
		err = errors.New("Invalid error code")
		return
	}
	return
}

type AppError struct{
	Code int
	Msg  string
}
func InitError(){
	ErrorMsgMap = make(map[int]string, 0)
	ErrorMsgMap[DB_ERROR_NO_CONNECTION] = "unknown error"
	ErrorMsgMap[DB_ERROR_QUERY_EXEC] = "unknown error"
	ErrorMsgMap[DB_ERROR_STMT_EXEC] = "unknown error"
	ErrorMsgMap[DB_ERROR_STMT_PREP] = "unknown error"
	ErrorMsgMap[OAUTH_ERROR_NO_USER] = "user is not registered with us"
	ErrorMsgMap[OAUTH_ERROR_INVALID_SIGNUPDATA] = "all the fields are compulsory. please fill it correctly"
	ErrorMsgMap[OAUTH_ERROR_AUTHENTICATION_FAILED]="authentication failed"
	ErrorMsgMap[OAUTH_ERROR_SECRET_NULL]="invlaid link"
	ErrorMsgMap[OAUTH_ERROR_USER_EXIST]="user already registered with this email id. please login"
	ErrorMsgMap[OAUTH_ERROR_USERNAME_PASS_MISMATCH]="username password mismatch"
	ErrorMsgMap[OAUTH_ERROR_NO_STAGE_DATA]="invalid link"
	ErrorMsgMap[OAUTH_ERROR_NULL_USER_OR_PASS]="username or password can not be empty"
	ErrorMsgMap[OAUTH_ERROR_EMAIL_TYPE_OR_SECRET_NULL]="invalid link. click here to resend"
	ErrorMsgMap[ERROR_JSON_BODY_NOT_READABLE]="invalid request. body is not readable"
	ErrorMsgMap[ERROR_JSON_BODY_NOT_PARSABLE]="invalid request. body is not readable"
	ErrorMsgMap[ERROR_INVALID_LINK]="invlaid link."
	ErrorMsgMap[ERROR_INVALID_NAME]="invalid name. name must contain alphabate, numericals, space( ), dot(.) only"
	ErrorMsgMap[ERROR_INVALID_EMAIL]="invalid email"
	ErrorMsgMap[ERROR_INVALID_PARAM]="invalid params for the request"
	ErrorMsgMap[ERROR_INVALID_VERIFICATION_LINK]="invalid verification link"
	ErrorMsgMap[ERROR_INVALID_PHONE_NO]="invalid phone number"
	ErrorMsgMap[ERROR_INVALID_METHOD_TYPE]="invlaid method type"
	ErrorMsgMap[ERROR_INVALID_CONTENT_TYPE]="invalid content type"
	ErrorMsgMap[ERROR_INVALID_REDIRECT_URI]="invalid redirect uri."
	ErrorMsgMap[ERROR_TEMPLATE_IS_NULL]="internal error"
}
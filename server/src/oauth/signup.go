package oauth

import (
	"build-concept/server/src/config"
	"build-concept/server/src/common"
	"log"
	"net/http"
	"strconv"
)

type UserData struct {
	Name     string
	Email    string
	Phno     string
	Password string
}
const(
	SIGNUP_EMAIL_VERIFICATION int = 1+ iota
	RESTE_PASSWD_VERIFICATION
)

func CheckUserValidation(email, phno string) (fn_err common.AppError){
	fn_err = common.AppError{Code:200, Msg:"",}
	filter := map[string]interface{}{
		"Email":email,
		"Phno":phno,
	}
	data, err := oauth_db.DbSelect("user", filter, []string{"id"})
	if err != nil{
		log.Println("CheckEmailValidation error in fetching data for user: ",email)
		fn_err.Code = common.DB_ERROR_STMT_EXEC
		fn_err.Msg = "db error"
		return
	}
	if len(data) > 0{
		fn_err.Code = common.OAUTH_ERROR_USER_EXIST
		fn_err.Msg = "user is already registered"
		return
	}
	return
}
func SignupAndSendEmail(user_data map[string]interface{})(fn_err common.AppError){
	fn_err = common.AppError{Code:200,}
	if len(user_data) < 0{
		log.Println("SignupAndSendEmail invalid data")
		fn_err.Code = common.ERROR_INVALID_PARAM
		return
	}
	secret, err := common.GenSecret()
	if err != nil{
		log.Println("Error in generating secret: ",secret)
		fn_err.Code = common.OAUTH_ERROR_IN_SECRET_GEN
		fn_err.Msg = "sign up error. try again"
		return
	}
	user_data["secret"] = secret
	pass, salt := GenPassHash(user_data["password"].(string))
	user_data["password"]= pass
	user_data["salt"]=salt
	_, err = oauth_db.DbInsert("staged", user_data)
	if err != nil{
		log.Println("Error in staging user data for email: ", user_data["email"])
		fn_err.Code = common.DB_ERROR_QUERY_EXEC
		return
	}
	e_data := common.EmailData{
		Email: 	user_data["email"].(string),
		Name:	user_data["first_name"].(string),
		Secret:	secret,
	}
	common.SendEmail(e_data, common.TMPL_TYPE_SIGNUP_VERIFICATION)
	return
}
func GetSignupHandler(cfg *config.Config) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Signup start...")
		resp := common.NewResponse()
		data, err := common.ParseJsonRequest(r)	
		if err.Code != 200{
			resp.ErrorStatusCode = err.Code
			resp.WriteError(w)
			return
		}
		log.Println("GetSignupHandler JOSNBody: ", data)
		if data["email"] == nil || data["first_name"] == nil ||  (data["phno"] ==nil) || (data["password"] == nil){
			log.Println("GetSignupHandler invalid param")
			resp.ErrorStatusCode = common.ERROR_INVALID_PARAM
			resp.WriteError(w)
			return
		}
		//data validation...
		if ok := common.EmailRegexValidation(data["email"].(string)); !ok{
			log.Println("invalid email. ", data["email"])
			resp.ErrorStatusCode = common.ERROR_INVALID_EMAIL
			resp.WriteError(w)
			return
		}
		if ok := common.PhoneRegexValidation(data["phno"].(string)); !ok{
			log.Println("invalid email. ", data["email"])
			resp.ErrorStatusCode = common.ERROR_INVALID_PHONE_NO
			resp.WriteError(w)
			return	
		}
		if ok := common.NameRegexValidation(data["first_name"].(string)); !ok{
			log.Println("invalid email. ", data["email"])
			resp.ErrorStatusCode = common.ERROR_INVALID_NAME
			resp.WriteError(w)
			return	
		}
		if data["last_name"] != nil{
			if ok := common.NameRegexValidation(data["last_name"].(string)); !ok{
				log.Println("invalid email. ", data["email"])
				resp.ErrorStatusCode = common.ERROR_INVALID_NAME
				resp.WriteError(w)
				return	
			}
		}
		//email and phone no validation...
		err = CheckUserValidation(data["email"].(string), data["phno"].(string))
		log.Println("CheckEmailValidation status:", err.Code)
		if err.Code != 200{
			resp.ErrorStatusCode = err.Code
			resp.WriteError(w)
			return
		}

		//register user and sending email data...
		log.Println("GetSignupHandler staging the data...")
		err = SignupAndSendEmail(data)
		if err.Code != 200{
			resp.ErrorStatusCode = err.Code
			resp.WriteError(w)
			return
		}
		log.Println("Signup success for email: ", data["email"])
		resp.Output["message"] = "Signup sucess. please verify to complete signup"
		resp.StatusCode = 200
		resp.WriteJson(w)
	}
}

func GetEmailVerificationHandler(cfg *config.Config) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request){
		fn_err := common.AppError{Code:200,}
		resp := common.NewResponse()
		fn_err = common.ParseGetRequest(r)
		if fn_err.Code != 200{
			log.Println("Invalid request: ", fn_err.Code)
			resp.ErrorStatusCode = fn_err.Code
			resp.WriteError(w)
			return
		}
		r.ParseForm()
		secret := r.Form.Get("secret")
		if secret == ""{
			log.Println("GetEmailVerificationHandler secret is null")
			resp.ErrorStatusCode = common.ERROR_INVALID_LINK
			resp.WriteError(w)
			return
		}
		email_type := r.Form.Get("email_type")
		if email_type == ""{
			log.Println("GetEmailVerificationHandler email type is null")
			resp.ErrorStatusCode = common.ERROR_INVALID_LINK
			resp.WriteError(w)
			return	
		}
		fn_err = EmailVerify(secret, email_type)
		if fn_err.Code != 200 {
			log.Println("EmailVerify failed: ", fn_err.Code)
			resp.ErrorStatusCode = fn_err.Code
			resp.WriteError(w)
			return
		}
		log.Println("Email verified for ", secret)
		resp.Output["message"]=fn_err.Msg
		resp.StatusCode = 200
		resp.WriteJson(w)
	}
}
func EmailVerify(secret string, email_type string)(fn_err common.AppError){
	fn_err = common.AppError{Code:200,}
	if email_type == "" || secret == ""{
		log.Println("Invalid data to verify")
		fn_err.Code = common.OAUTH_ERROR_EMAIL_TYPE_OR_SECRET_NULL
		return
	}
	udata, err := FetchStageDataBySec(secret)
	if err.Code != 200{
		fn_err.Code = err.Code
		return
	}
	if len(udata) <= 0{
		fn_err.Code = common.OAUTH_ERROR_NO_STAGE_DATA
		return
	}
	log.Println("stage data found for the secret:", secret)
	vtype, err1 := strconv.Atoi(email_type)
	if err1 != nil{
		log.Println("Email type conversion error: ", err1)
		fn_err.Code = common.ERROR_INVALID_VERIFICATION_LINK
		return
	}
	switch(vtype){
	case SIGNUP_EMAIL_VERIFICATION:
		mdata := map[string]interface{}{
			"email":udata["email"],
			"first_name":udata["first_name"],
			"last_name":udata["last_name"],
			"phno":udata["phno"],
			"city":udata["city"],
			"password":udata["password"],
			"salt":udata["salt"],
		}
		fn_err = RegisterUser(mdata)
	case RESTE_PASSWD_VERIFICATION:
		mdata := map[string]interface{}{
			"password":udata["password"],
		}
		fn_err = ResetPassword(mdata)
	}
	return
}
func RegisterUser(data map[string]interface{})(fn_err common.AppError){
	fn_err = common.AppError{Code:200,}
	if len(data) <= 0{
		log.Println("RegisterUser nill data to process", data)
		fn_err.Code = common.OAUTH_ERROR_INVALID_SIGNUPDATA
		return
	}
	log.Println("data:", data)
	if data["email"] == nil || data["first_name"] == nil ||  (data["phno"] == nil) || (data["password"] == nil){
		log.Println("RegisterUser invalid data to process")
		fn_err.Code = common.OAUTH_ERROR_INVALID_SIGNUPDATA
		return
	}
	log.Println("Registering user Email:", data["email"])
	fn_err = CheckUserValidation(data["email"].(string), data["phno"].(string))
	log.Println("CheckEmailValidation status:", fn_err.Code)
	if fn_err.Code != 200{
		return
	}
	//add entry in user table
	mdata := map[string]interface{}{
		"email":		data["email"],
		"first_name":	data["first_name"],
		"last_name":	data["last_name"],
		"phno":			data["phno"],
		"city":			data["city"],
	}
	result,err := oauth_db.DbInsert("user", mdata)
	if err != nil{
		log.Println("RegisterUser Db error: ", err)
		fn_err.Code =common.DB_ERROR_QUERY_EXEC
		return
	}
	uid,_ := result.LastInsertId()
	// Oauth data insertion
	log.Println("RegisterUser user created with ID: ", uid)
	odata := map[string]interface{}{
		"id":uid,
		"password":data["password"],
		"salt":data["salt"],
	}
	_, err = oauth_db.DbInsert("oauth", odata)
	if err != nil{
		log.Println("RegisterUser Db error: ", err)
		fn_err.Code =common.DB_ERROR_QUERY_EXEC
		return
	}
	log.Println("User has been registered")
	return
}

func ResetPassword(data map[string]interface{})(fn_err common.AppError){
	//TODO:
	fn_err = common.AppError{Code:200,}
	return

}
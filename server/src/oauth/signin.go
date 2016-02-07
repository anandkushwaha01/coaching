package oauth

import (
	"concept-build/server/src/config"
	"concept-build/server/src/session"
	"concept-build/server/src/common"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Authorize(username, password string) (fn_err common.AppError){
	fn_err = common.AppError{Code:200,}
	if username == "" || password == ""{
		log.Println("Authorize failed username and password empty")
		fn_err.Code = common.OAUTH_ERROR_NULL_USER_OR_PASS
		return
	}
	filter := map[string]interface{}{
		"email":username,
	}
	udata, err := oauth_db.DbSelect("user", filter, []string{"id"})
	if err != nil{
		log.Println("Authorize error in fetching user details.", err)
		fn_err.Code = common.DB_ERROR_QUERY_EXEC
		return
	}
	if len(udata) <= 0{
		log.Println("Authorize use not found ", username)
		fn_err.Code = common.OAUTH_ERROR_NO_USER
		return
	}
	var id int
	if udata[0]["id"] != ""{
		id, err = strconv.Atoi(udata[0]["id"])
		if err != nil{
			log.Println("Authorize error string id conversion", err)
		}
	}
	filter = map[string]interface{}{
		"id":id,
	}
	odata, err := oauth_db.DbSelect("oauth", filter, []string{"password", "salt"})
	if err != nil{
		log.Println("Authorize error in fetching password details", err)
		fn_err.Code = common.DB_ERROR_QUERY_EXEC
		return
	}
	if len(odata) <= 0{
		log.Println("Authorize password not found ", username)
		fn_err.Code = common.OAUTH_ERROR_NO_USER
		return
	}
	hashPass := GenPassHashWithSalt(password, odata[0]["salt"])
	if len(hashPass) <= 0{
		log.Println("Authorize error in getting hash for the given password")
		fn_err.Code = common.OAUTH_ERROR_AUTHENTICATION_FAILED
		return
	}
	if odata[0]["password"] != hashPass{
		log.Println("Authorize username password mismatch for user: ", username)
		fn_err.Code = common.OAUTH_ERROR_USERNAME_PASS_MISMATCH
	}
	return
}

func GetLoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := common.NewResponse()
		fn_err := common.ParseFormRequest(r)
		if fn_err.Code != 200{
			resp.ErrorStatusCode = fn_err.Code
			resp.WriteError(w)
			return
		}
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if username == "" || password == "" {
			log.Println("Invalid params..")
			resp.ErrorStatusCode = common.OAUTH_ERROR_NULL_USER_OR_PASS
			resp.WriteError(w)
			return
		}
		sess, err := session.GetSession(w, r)
		if err != nil || sess == nil {
			log.Println("Error in getting session. Error:", err, " sess: ", sess)
		} else {
			if sess.Get("user") != nil {
				user := sess.Get("user").(string)
				if user == username {
					log.Println(user, " authorized from session.")
					resp.StatusCode = 200
					resp.Output["msg"]="already logged in"
					resp.WriteJson(w)
					return
				}
				sess, err = session.NewSession(w, r)
				if err != nil || sess == nil {
					log.Println("Error in getting session. Error:", err, " sess: ", sess)
				}
			}
		}
		if err := Authorize(username, password); err.Code != 200{
			resp.ErrorStatusCode = err.Code
			resp.WriteError(w)
			return
		}
		log.Println("Authorization Success for ", username)
		if sess != nil {
			sess.Set("user", username)
			session.SaveSession(sess)
			log.Println("Session has been saved for ", username)
		}
		resp.StatusCode = 200
		resp.Output["msg"]="login success"
		resp.WriteJson(w)
	}
}
func GetLogoutHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]interface{})
		defer func() {
			w.Header().Set("Content-Type", "application/json")
			if resp["error"] != nil {
				w.WriteHeader(500)
			} else {
				w.WriteHeader(200)
			}
			encoder := json.NewEncoder(w)
			err := encoder.Encode(resp)
			if err != nil {
				log.Println("error in writing data")
				w.Write([]byte("unknow error"))
			}
		}()
		session.DeleteSession(w, r)
		log.Println("logout success")
		resp["code"] = 200
		return
	}
}

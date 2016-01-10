package oauth
import(
	"log"
	"net/http"
	"concept-build/server/src/session"
	"concept-build/server/src/config"
	"encoding/json"
)

func Authorize(username, password string) (ok bool, auth_error error){
	ok = true
	return
}

func GetLoginHandler(cfg *config.Config) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request){
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
		if r.Method != "POST"{
			resp["error"] = "Request must be POST"
			return
		}
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		if username == "" || password == ""{
			log.Println("Invalid params..")
			resp["error"] ="Invalid param"
			return
		}
		sess, err := session.GetSession(w, r)
		if err != nil || sess == nil{
			log.Println("Error in getting session. Error:", err, " sess: ", sess)
		}else{
			if sess.Get("user") != nil{
				user := sess.Get("user").(string)
				if user == username{
					log.Println(user, " authorized from session.")
					resp["code"]=200
					return
				}
				sess, err = session.NewSession(w, r)
				if err != nil || sess == nil{
					log.Println("Error in getting session. Error:", err, " sess: ", sess)
				}
			}
		}
		if ok, err := Authorize(username, password); (err != nil || !ok){
			log.Println("Authorization failed: ", err, " Status:", ok)
			resp["error"] = "Authorization failed."
			return
		}
		log.Println("Authorization Success for ", username)
		if sess != nil{
			sess.Set("user", username)
			session.SaveSession(sess)
			log.Println("Session has been saved for ", username)
		}
		resp["code"] = 200
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
		resp["code"]=200
		return
	}
}
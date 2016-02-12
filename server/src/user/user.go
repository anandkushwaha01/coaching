package user
import(
	"log"
	"errors"
	"build-concept/server/src/config"
	"build-concept/server/src/common"
	"build-concept/server/src/session"
	"html/template"
	"net/http"
)
var user_db *common.DBContext
func Init(cfg *config.UserCfg) (error){
	var err error
	user_db, err = common.Dbconnect(cfg.DSN, cfg.MaxIdleConns, cfg.MaxOpenConns)
	if err != nil{
		log.Println("User DB connection failed. Error: ", err)
		return errors.New("User DB connection failed")
	}
	log.Println("User DB connection success")
	return nil
}

func GetHomeHandler(cfg *config.Config) (http.HandlerFunc){
	tmpl_path := cfg.Base.Path
	log.Println("template path: ", tmpl_path)
	home_tmpl, err := template.ParseFiles(tmpl_path+"views/master.tmpl")
	if err != nil{
		log.Println("Error in parsing the template: ", err)
	}
	return func(w http.ResponseWriter, r *http.Request){
		log.Println("Calling home page...")
		data := make(map[string]interface{}, 0)
		if home_tmpl == nil{
			log.Println("template is not initialized")
			resp := common.NewResponse()
			resp.ErrorStatusCode = common.ERROR_TEMPLATE_IS_NULL
			resp.WriteError(w)	
			return
		}
		sess, err := session.GetSession(w, r)
		if err == nil && sess != nil{
			log.Println("Checking user session")
			if sess.Get("user") != nil{
				log.Println("sigining in as user: ", sess.Get("user"))
				data["User"]=sess.Get("user")
				home_tmpl.Execute(w, data)
				return
			}
		}
		home_tmpl.Execute(w, data)
		return
	}
}
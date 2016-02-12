package oauth
import(
	"log"
	"errors"
	"build-concept/server/src/config"
	"build-concept/server/src/common"
	"golang.org/x/crypto/scrypt"
	"github.com/pborman/uuid"
	"encoding/base64"
)

var oauth_db *common.DBContext
func Init(cfg *config.OauthCfg) (error){
	var err error
	oauth_db, err = common.Dbconnect(cfg.DSN, cfg.MaxIdleConns, cfg.MaxOpenConns)
	if err != nil{
		log.Println("Oauth DB connection failed. Error: ", err)
		return errors.New("Oauth DB connection failed")
	}
	log.Println("Oauth DB connection success")
	return nil
}

func FetchStageDataBySec(secret string) (data map[string]string, fn_err common.AppError){
	fn_err = common.AppError{Code:200,}
	if secret == ""{
		log.Println("FetchStageDataBySec null secret error")
		fn_err.Code = common.OAUTH_ERROR_SECRET_NULL
		return
	}
	filter := map[string]interface{}{
		"secret":secret,
	}
	dbdata, err := oauth_db.DbSelect("staged", filter, []string{})
	if err != nil{
		log.Println("Error in fetching stage data for secret:", secret, " error:", err)
		fn_err.Code = common.DB_ERROR_QUERY_EXEC
		return
	} 
	if len(dbdata) <= 0{
		log.Println("no stage data for secret:",secret)
		fn_err.Code =common.OAUTH_ERROR_NO_STAGE_DATA
		return
	}
	data = dbdata[0]
	fn_err.Code = 200
	log.Println("Staged data found for secret: ",secret)
	return
}

func GenPassHash(pass string)(hashPass, salt string){
	token := (uuid.NewRandom())
	salt = base64.URLEncoding.EncodeToString([]byte(token))
	hash,_ := scrypt.Key([]byte(pass), []byte(salt), 16384, 8, 1, 32)
	hashPass = base64.URLEncoding.EncodeToString(hash)
	return
}
func GenPassHashWithSalt(pass, salt string) (hashPass string){
	hash, _ := scrypt.Key([]byte(pass), []byte(salt), 16384, 8, 1, 32)
	hashPass = base64.URLEncoding.EncodeToString(hash)
	return	
}
package user
import(
	"log"
	"errors"
	"concept-build/server/src/config"
	"concept-build/server/src/common"
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
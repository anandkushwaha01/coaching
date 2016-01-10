package session
import(
	"log"
	"net/http"
	"errors"
)
var mgr *Manager
func Init(){
	var err error
	mgr, err = NewManager("redis","oauth", 60*60*48)
	if err != nil{
		log.Println("Error in initializing session", err)
	}
	log.Println("session has been initialized...")
}
func NewSession(w http.ResponseWriter, r *http.Request) (Session, error){
	if mgr == nil{
		log.Println("session is not initialized")
		return	nil, errors.New("nil session manager")
	}
	return mgr.SessionNew(w, r), nil
}
func GetSession(w http.ResponseWriter, r *http.Request) (Session, error){
	if mgr == nil{
		log.Println("session is not initialized")
		return	nil, errors.New("nil session manager")
	}
	return mgr.SessionStart(w, r), nil
}
func DeleteSession(w http.ResponseWriter, r *http.Request){
	if mgr == nil{
		log.Println("session is not initialized")
		return	
	}
	mgr.SessionDestroy(w,r)
}
func SaveSession(sess Session){
	if mgr == nil{
		log.Println("session is not initialized")
		return	
	}
	err := mgr.provider.SessionSave(sess)
	if err != nil{
		log.Println("Error in saving session ",err)
	}
}
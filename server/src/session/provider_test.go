package session
import(
	"testing"
	"encoding/json"
	"time"
)

type TestSessionStore struct {
    Sid          string                      // unique session id
    TimeAccessed time.Time                   // last access time
    Value        map[string]interface{} // session value stored inside
}
type Message struct {
    Name string
    Body string
    Time time.Time
}
func TestJsonEncoding(t *testing.T){
	v := make(map[string]interface{}, 0)
	v["email"] = "anandkushwaha01@gmail.com"
	v["before"]="befoe"
    newsess := &TestSessionStore{Sid: "sid", Value:v, TimeAccessed:time.Now()}
    val, err := json.Marshal(newsess)
    if err != nil{
    	t.Log("Error in marshalling data", err )
    }
    t.Log("val : ", val)
    data := &TestSessionStore{}
    err = json.Unmarshal(val, data)
    t.Log("Unmarshal: ", data)
}
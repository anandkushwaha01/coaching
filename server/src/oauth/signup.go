package oauth

import (
	"concept-build/server/src/config"
	"encoding/json"
	"log"
	"net/http"
)

type UserData struct {
	name     string
	email    string
	phno     string
	password string
}

func GetSignupHandler(cfg *config.Config) http.HandlerFunc {
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
		//TODO
		resp["code"] = 200
		return
	}
}

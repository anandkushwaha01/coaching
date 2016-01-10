package main

import (
	"flag"
	"github.com/paytm/grace"
	"github.com/paytm/logging"
	"log"
	"net/http"
	"concept-build/server/src/oauth"
	"concept-build/server/src/common"
	"concept-build/server/src/session"
	"concept-build/server/src/config"
)

func main() {
	var cfg config.Config
	flag.Parse()
	logging.LogInit()
	ok := config.ReadConfig(&cfg, ".") || config.ReadConfig(&cfg, "/etc/")
	if !ok {
		log.Fatal("failed to read config")
		return
	}
	pool, err := common.InitRedis(cfg.Redis.Address)
	if err != nil{
		log.Println("Redis init failed. Error:", err)
		return
	}
	session.ProviderInit(pool)
	session.Init()
	
	http.Handle("/login", oauth.GetLoginHandler(&cfg))
	http.Handle("/logout", oauth.GetLogoutHandler(&cfg))
	log.Fatal(grace.Serve(":"+cfg.Server.Port, nil))
}

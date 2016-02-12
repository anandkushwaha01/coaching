package main

import (
	"flag"
	"github.com/paytm/grace"
	"github.com/paytm/logging"
	"log"
	"net/http"
	"build-concept/server/src/oauth"
	"build-concept/server/src/user"
	"build-concept/server/src/common"
	"build-concept/server/src/session"
	"build-concept/server/src/config"
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
	err := oauth.Init(&cfg.Oauth)
	if err != nil{
		log.Println("Fatal error. closing the app")
	}
	err = user.Init(&cfg.User)
	if err != nil{
		log.Println("Fatal error. closing the app")
	}
	pool, err := common.InitRedis(cfg.Redis.Address)
	if err != nil{
		log.Println("Redis init failed. Error:", err)
		return
	}
	common.UtilInit(&cfg)
	common.InitError()
	session.ProviderInit(pool)
	session.Init()
	
	http.Handle("/home", user.GetHomeHandler(&cfg))
	http.Handle("/login", oauth.GetLoginHandler(&cfg))
	http.Handle("/logout", oauth.GetLogoutHandler(&cfg))
	http.Handle("/signup", oauth.GetSignupHandler(&cfg))
	http.Handle("/emailverify", oauth.GetEmailVerificationHandler(&cfg))
	log.Fatal(grace.Serve(":"+cfg.Server.Port, nil))
}

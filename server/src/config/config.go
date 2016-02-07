package config

import (
	"gopkg.in/gcfg.v1"
	"log"
)

type Config struct {
	Oauth  OauthCfg
	User   UserCfg
	Server ServerCfg
	Base   BaseCfg
	Redis  RedisCfg
	Smtp 	SmtpCfg
}
type SmtpCfg struct{
	Address string
}
type DBCfg struct {
	DSN          string
	MaxIdleConns int
	MaxOpenConns int
	Querylog     bool
}
type OauthCfg struct{
	DBCfg
}
type UserCfg struct{
	DBCfg
}
type ServerCfg struct {
	Port string
}
type BaseCfg struct {
	Path string
}
type AuthHeaderCfg struct {
	User     string
	Password string
}
type RedisCfg struct {
	Address string
}

func ReadConfig(cfg *Config, path string) bool {
	err := gcfg.ReadFileInto(cfg, path+"/concept-build.conf")
	if err == nil {
		return true
	}
	log.Println("Error: ", err)
	return false
}

package config
import(
	"log"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	DB         	DBCfg
	Server     	ServerCfg
	Base       	BaseCfg
	Redis 		RedisCfg
}
type DBCfg struct {
	DSN          string
	MaxIdleConns int
	MaxOpenConns int
	Querylog     bool
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
type RedisCfg struct{
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
package session

import (
	"build-concept/server/src/common"
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
	"time"
)

type RedisProvider struct {
	lock  sync.Mutex // lock
	Store *common.RedisStorage
}

var pder *RedisProvider

func (pder *RedisProvider) SessionInit(sid string) (Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[string]interface{}, 0)
	newsess := &SessionStore{Sid: sid, TimeAccessed: time.Now(), Value: v}
	return newsess, nil
}

func (pder *RedisProvider) SessionRead(sid string) (Session, error) {
	log.Println("Reading session: ", sid)
	rs := pder.Store
	if rs == nil {
		return nil, errors.New("Redis connection not found")
	}
	data := &SessionStore{}
	data.Value = make(map[string]interface{}, 0)
	err := rs.Get([]byte(sid), data)
	if err != nil {
		log.Println("Error In reading session: ", err)
		return nil, err
	}
	if data.Sid == "" {
		log.Println("no data found")
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return data, nil
}

func (pder *RedisProvider) SessionDestroy(sid string) error {
	rs := pder.Store
	if rs == nil {
		return errors.New("Redis connection not found")
	}
	err := rs.Del([]byte(sid))
	return err
}

func (pder *RedisProvider) SessionSave(ses Session) error {
	sid := ses.SessionID()
	rs := pder.Store
	if rs == nil {
		return errors.New("Redis connection not found")
	}
	// data := ses.value["user_info"].(common.MerchantStatusInfoCtx)
	err := rs.SetEx(ses, []byte(sid), 60*60*48)
	return err
}

func ProviderInit(pool *redis.Pool) {
	log.Println("ProviderInit: starting...")
	pder = &RedisProvider{}
	pder.Store = common.NewRedisStorage(pool, "ses", ":")
	Register("redis", pder)
	log.Println("ProviderInit: finished")
}

package common

import (
	"github.com/dilleeppk7/osin"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	//t.Log("starting redis test")
	t.Logf("%s \n", "starting redis test")
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	if pool == nil {
		t.Error("failed to create pool")
	}
	r := NewRedisStorage(pool, "at", ":")
	if r == nil {
		t.Error("failed to create storage")
	}

	previous_token := &osin.AccessData{
		Client: &osin.DefaultClient{
			Id:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://localhost:14000/appauth",
		},
		AuthorizeData: &osin.AuthorizeData{
			Client: &osin.DefaultClient{
				Id:          "1234",
				Secret:      "aabbccdd",
				RedirectUri: "http://localhost:14000/appauth",
			},
			Code:        "77777",
			ExpiresIn:   3600,
			CreatedAt:   time.Now(),
			RedirectUri: "http://localhost:14000/appauth",
		},

		AccessToken: "88888",
		ExpiresIn:   3600,
		CreatedAt:   time.Now(),
	}
	token := &osin.AccessData{
		Client: &osin.DefaultClient{
			Id:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://localhost:14000/appauth",
		},
		AuthorizeData: &osin.AuthorizeData{
			Code:        "9999",
			ExpiresIn:   3600,
			CreatedAt:   time.Now(),
			RedirectUri: "http://localhost:14000/appauth",
		},
		AccessData:  previous_token,
		AccessToken: "9999",
		ExpiresIn:   3600,
		CreatedAt:   time.Now(),
	}
	r.Set(token, []byte(token.AccessToken))

	var data_got *osin.AccessData
	r.Get([]byte(token.AccessToken), &data_got)
	if data_got.AccessToken != "9999" {
		t.Error("data access token", data_got.AccessToken)

	} else if data_got.AuthorizeData.Code == "9999" {
		t.Error("Authorize code  not equal ")
		t.Logf("before %V \n", token)
		t.Logf("after %V \n", data_got)
	} else if data_got.AccessData.AuthorizeData.Code != "77777" {
		t.Error(" inner object  test failed")
	}

	//t.Log(fmt.Sprintf("%v \n", dddd))

}

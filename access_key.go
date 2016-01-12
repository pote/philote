package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type AccessKey struct {
	Channels map[string]string  `json:"channels"`
	Token string                `json:"-"`
}

func LoadKey(token string) (*AccessKey, error) {
	ak := &AccessKey{Token: token}
	r := RedisPool.Get()

	rawKey, err := redis.String(r.Do("GET", "philote:token:" + token))
	r.Close()
	if err != nil {
		return ak, err
	}
	
	if rawKey  == "" {
		return ak, InvalidTokenError{"unknown token"}
	}


	err = json.Unmarshal([]byte(rawKey), &ak); if err != nil {
		return ak, InvalidTokenError{"invalid token data: " + err.Error()}
	}

	return ak, nil
}

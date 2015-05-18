package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type AccessToken struct {
	Hub      string
	Channels []string
}

func ParseAccessToken(tokenString string) (at *AccessToken, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		redis := RedisPool.Get()
		defer redis.Close()

		key, err := redis.Do("GET", "hubs:"+t.Claims["hub"].(string))

		if key == nil {
			err = errors.New("[SECURITY] Couldn't find hub matching credentials")
		}

		return key, err
	})

	if err != nil {
		return
	}

	if !token.Valid {
		err = errors.New("JWT is invalid")
		return
	}

	rawChannels := token.Claims["channels"].([]interface{})
	channels := make([]string, len(rawChannels))
	for i, channel := range rawChannels {
		channels[i] = channel.(string)
	}

	at = &AccessToken{
		Hub:      token.Claims["hub"].(string),
		Channels: channels,
	}

	return
}

func (at *AccessToken) CanAccess(channel string) bool {
	for _, c := range at.Channels {
		if c == channel {
			return true
		}
	}
	return false
}

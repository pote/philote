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
		// Has to return the Hub's Secret Key (can access t.Claims["hub"] to get
		// the hub's Access Key.)
		return []byte("deadbeefsecret"), nil
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

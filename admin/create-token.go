package main

import (
	"encoding/json"
	"github.com/pote/redisurl"
)

type AccessKey struct {
	Token       string   `json:"-"`
	Read        []string `json:"read"`
	Write       []string `json:"write"`
	AllowedUses int      `json:"allowed_uses"`
	Uses        int      `json:"uses"`
}

func createAccessKey(token string, read, write []string, allowedUses int)  error {
	r, err := redisurl.Connect(); if err != nil {
		return err
	}
	defer r.Close()

	ak := &AccessKey{
		Token: token,
		Read: read,
		Write: write,
		AllowedUses: allowedUses,
		Uses: 0,
	}

	data, err := json.Marshal(ak); if err != nil {
		return err
	}

	_, err =  r.Do("SET", "philote:access_key:" + ak.Token, string(data)); if err != nil {
		return err
	}

	return nil
}

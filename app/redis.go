package main

import (
	"sync"
	"time"
)

type Redis struct {
	datastore sync.Map
}

type RedisValue struct {
	Value      string
	Expiration time.Time
}

func NewRedis() *Redis {
	return &Redis{
		datastore: sync.Map{},
	}
}

func (red *Redis) Set(key string, value string, expiry time.Duration) error {
	var expiration time.Time

	if expiry == 0 {
		expiration = time.Unix(0, 0)
	} else {
		expiration = time.Now().Add(expiry)
	}
	rvalue := RedisValue{
		Value:      value,
		Expiration: expiration,
	}
	red.datastore.Store(key, rvalue)
	return nil
}

func (red *Redis) Get(key string) (string, error) {
	val, found := red.datastore.Load(key)
	if !found {
		return "", nil
	}

	rvalue := val.(RedisValue)
	if rvalue.Expiration != time.Unix(0, 0) && time.Now().After(rvalue.Expiration) {
		red.datastore.Delete(key)
		return "", nil
	}
	return rvalue.Value, nil
}

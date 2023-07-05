package main

import (
	"fmt"
	"sync"
)

type Redis struct {
	datastore sync.Map
}

func NewRedis() *Redis {
	return &Redis{
		datastore: sync.Map{},
	}
}

func (red *Redis) Set(key string, value string) error {
	red.datastore.Store(key, value)
	return nil
}

func (red *Redis) Get(key string) (string, error) {
	val, found := red.datastore.Load(key)
	if !found {
		return "", fmt.Errorf("no such value found")
	}
	value := val.(string)
	return value, nil
}

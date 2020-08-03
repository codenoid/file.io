package storage

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

// Redis storage
type Redis struct {
	Conn *redis.Client
}

// Connect will start redis connection
func (r *Redis) Connect(addr, password, database string) {
	intf, err := strconv.Atoi(database)
	if err != nil {
		panic(err)
	}
	r.Conn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       intf,     // use default DB
	})
}

// Set set bytes file to redis using unique id as key
func (r *Redis) Set(key string, value []byte, expSec int) error {
	exp := time.Duration(expSec) * time.Second
	return r.Conn.Set(key, value, exp).Err()
}

// Get get bytes from redis and write bytes as response (file)
func (r *Redis) Get(key string) ([]byte, error) {
	d := r.Conn.Get(key)
	return d.Bytes()
}

// Del uwu
func (r *Redis) Del(key string) {
	r.Conn.Del(key)
}

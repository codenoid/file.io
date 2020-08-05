package storage

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

type StorageHandler interface {
	Set(key string, value []byte, expired time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string)
}

// Connect to specified `database`, return error if given `database` are invalid or unknown or
// panic if failed to connect
func Connect(database string) (StorageHandler, error) {

	databaseURI, err := url.Parse(database)
	if err != nil {
		return nil, errors.New("unable to parse database URI: " + err.Error())
	}

	// put general information in there
	host := databaseURI.Host
	pass, _ := databaseURI.User.Password()
	db := strings.TrimLeft(databaseURI.Path, "/")

	switch databaseURI.Scheme {
	case "badger":
		stg := Badger{}
		stg.Connect(databaseURI.Path)
		return &stg, nil
	case "redis":
		stg := Redis{}
		stg.Connect(host, pass, db)
		return &stg, nil
	default:
		return nil, errors.New("unknown storage: " + database)
	}

}

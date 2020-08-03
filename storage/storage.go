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

func Connect(database, databaseType string) (StorageHandler, error) {

	switch databaseType {
	case "badger":
		stg := Badger{}
		stg.Connect(database)
		return &stg, nil
	case "redis":
		connURI, err := url.Parse(database)
		if err != nil {
			return nil, errors.New("unable to parse database URI: " + err.Error())
		}

		host := connURI.Host
		pass, _ := connURI.User.Password()
		db := strings.TrimLeft(connURI.Path, "/")
		stg := Redis{}
		stg.Connect(host, pass, db)
		return &stg, nil
	default:
		return nil, errors.New("unknown storage: " + databaseType)
	}

}

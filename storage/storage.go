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
	case "redis":
		connURI, err := url.Parse(database)
		if err != nil {
			return nil, errors.New("unable to parse database URI: " + err.Error())
		}

		host := connURI.Host
		pass, _ := connURI.User.Password()
		db := strings.TrimLeft(connURI.Path, "/")
		rds := Redis{}
		rds.Connect(host, pass, db)
		return &rds, nil
	}

	return nil, errors.New("unknown storage: " + connURI.Scheme)
}

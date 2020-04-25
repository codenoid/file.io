package storage

import "errors"

// Storage allow external use of fileio app functionality
type Storage struct {

	// Address contain redis/mongodb connection uri, only used when
	// we use Redis or MongoDB
	Address string

	// Path contain path to directory that will be used for saving
	// temporary file
	Path string

	// Type define storage name that used for saving temporary data
	Type string

	// Storage Definition
	Redis
}

// Connect will connect
func (s *Storage) Connect(host, username, password, database string) error {
	switch s.Type {
	case "redis":
		s.Redis.Connect(host, password, database)
		return nil
	default:
		return errors.New("storage not found")
	}
}

// Set will save data to specified storage
func (s *Storage) Set(key string, value []byte, expSec int) error {
	switch s.Type {
	case "redis":
		return s.Redis.Set(key, value, expSec)
	default:
		return errors.New("storage not found")
	}
}

// Get ok
func (s *Storage) Get(key string) ([]byte, error) {
	switch s.Type {
	case "redis":
		return s.Redis.Get(key)
	default:
		return nil, errors.New("storage not found")
	}
}

// Del delete data form database
func (s *Storage) Del(key string) {
	switch s.Type {
	case "redis":
		s.Redis.Del(key)
	}
}

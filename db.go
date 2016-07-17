package wtf

import (
	"bytes"
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	levelDB *leveldb.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Database{levelDB: db}, nil
}

func (d Database) Unpack(key string, target interface{}) error {
	data, err := d.levelDB.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(data)).Decode(target)
}

func (d Database) Set(key string, target interface{}) error {
	buf := bytes.Buffer{}

	if err := json.NewEncoder(&buf).Encode(target); err != nil {
		return err
	}

	return d.levelDB.Put([]byte(key), buf.Bytes(), nil)
}

package database

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	*leveldb.DB
}

func NewLevelDB(path string) (DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &LevelDB{db}, err
}

func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	v, err := ldb.DB.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return []byte{}, nil
	}
	return v, err
}

func (ldb *LevelDB) Set(key, value []byte) error {
	return ldb.Put(key, value, nil)
}

func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.DB.Delete(key, nil)
}

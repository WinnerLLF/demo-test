package leveldb

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var db *leveldb.DB
var err error

func init() {
	db, err = initDB("db")
	if err != nil {
		fmt.Printf("initDB err:%v\n", err)
	}
}

func initDB(name string) (*leveldb.DB, error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err = leveldb.OpenFile(fmt.Sprintf("./json/cache/%v", name), o)
	if err != nil {
		fmt.Printf("OpenFile err:%v\n", err)
	}

	return db, err
}

func WriteDB(key, value string) error {
	batch := new(leveldb.Batch)
	batch.Put([]byte(key), []byte(value))
	if err := db.Write(batch, nil); err != nil {
		fmt.Printf("Write err:%v\n", err)
		return err
	}

	return nil
}

func ReadDB(key string) []byte {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Printf("OpenFile err:%v\n", err)
		return []byte{}
	}

	return data
}

func ReadListDB() []string {
	var privateKey []string
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		privateKey = append(privateKey, string(value))
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		fmt.Printf("iter.Error:%v\n", err)
		return privateKey
	}

	return privateKey
}

func DeleteByKey(key string) error {
	if err := db.Delete([]byte(key), nil); err != nil {
		return err
	}

	return nil
}

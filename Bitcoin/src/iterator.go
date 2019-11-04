package main

import (
	"github.com/boltdb/bolt"
	"log"
)

/*
	迭代器
*/
type Iterator struct {
	db      *bolt.DB
	curHash []byte
}

/*
	初始化
*/
func NewIterator(bc *BlockChain) *Iterator {
	return &Iterator{
		db:      bc.db,
		curHash: bc.tail,
	}
}

/*
	游标
*/
func (it *Iterator) Next() *Block {
	var block  *Block
	_=it.db.View(func(tx *bolt.Tx) error {
		b :=tx.Bucket([]byte(bucket_name))
		if b == nil {
			log.Fatal("Bucket 不能为nil")
		}
		blockInfo :=b.Get([]byte(it.curHash))
		block =Deserialize(blockInfo)
		it.curHash = block.PrevHash
		return nil
	})
	return block
}

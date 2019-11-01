package main

import "github.com/boltdb/bolt"

const genesis_chain = "I'm Nakamoto. This is my first creative chain" //创世语
const db_name = "BlockChain.db"                                       //数据库名
const bucket_name = "BlockBucket"                                     //桶名
const last_Hashkey = "lastHashKey"                                    //哈希值

/*
	区块链
*/
type BlockChain struct {
	db   *bolt.DB //数据库句柄
	tail []byte   //末尾哈希值
}

/*
	初始化
*/
func NewBlockChain() *BlockChain {
	db, hash := InitBolt()
	if db == nil {
		return nil
	}
	return &BlockChain{
		db:   db,
		tail: hash,
	}
}

/*
	创建数据库，如果存在则不创建。
*/
func InitBolt() (*bolt.DB, []byte) {
	var hash []byte
	db, err := bolt.Open(db_name, 0600, nil)
	if err != nil {
		return nil, nil
	}
	//事务
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		if b == nil {
			//创建数据库
			b, err = tx.CreateBucket([]byte(bucket_name))
			if err != nil {
				return nil
			}
			coinBase := CoinBaseTX("", genesis_chain)
			//创世
			genesis := NewBlock(nil, []*Transaction{coinBase})
			//添加数据
			_ = b.Put(genesis.Hash, genesis.Serialize())
			//更新Hash值
			_ = b.Put([]byte(last_Hashkey), genesis.Hash)
			//获取hash
			hash = genesis.Hash
		} else {
			hash = b.Get([]byte(last_Hashkey))
		}
		return nil
	})
	return db, hash
}

/*
	添加区块链
*/
func (bc *BlockChain) Add(txs []*Transaction) {
	last := bc.tail
	block := NewBlock(last, txs)
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		if b == nil {
			return nil
		}
		_ = b.Put(block.Hash, block.Serialize())
		_ = b.Put([]byte(last_Hashkey), block.Hash)
		bc.tail = block.Hash
		return nil
	})
}

/*
	UTXO
*/
func (bc *BlockChain)FindUTXO(data string) []TXOutPut {
	var outputs []TXOutPut
	spentMap :=make(map[string]int64)
	it :=NewIterator(bc)	//遍历区块链
	for{
		block :=it.Next()
		for _,tx :=range block.Transactions{	//遍历交易
			for outPutIndex,outPut :=range tx.TXOutPuts{	//遍历outputs

			}
		}
	}
}
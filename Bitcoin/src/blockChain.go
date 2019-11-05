package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
)

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
	var lastHash []byte
	db, err := bolt.Open(db_name, 0600, nil)
	if err != nil {
		return nil
	}
	_=db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		if b == nil {
			b, err = tx.CreateBucket([]byte(bucket_name))
			if err != nil {
				return nil
			}
			coinBaseTx := CoinBaseTX("19qyxmoAxYqk76CfPZotXkVkvdnpsHiGFb", genesis_chain)
			genesis := NewBlock(nil, []*Transaction{coinBaseTx})
			//添加数据
			_ = b.Put(genesis.Hash, genesis.Serialize())
			//更新最后一个hash值
			_ = b.Put([]byte(last_Hashkey), genesis.Hash)
			lastHash = genesis.Hash
		} else {
			lastHash = b.Get([]byte(last_Hashkey))
		}
		return nil
	})

	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

/*
	添加区块链
*/
func (bc *BlockChain) Add(txs []*Transaction) {
	var validTxs []*Transaction
	for _,tx := range txs{
		if bc.VerifyTransaction(tx) {
			validTxs = append(validTxs,tx)
		}else {
			fmt.Printf("签名校验失败的交易:%x\n", tx.Txid)
		}
	}

	last := bc.tail
	block := NewBlock(last, validTxs)
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
	未消费详情
*/
type UTXOInfo struct {
	output TXOutPut
	intdex int64
	txid   []byte
}

/*
	UTXO
*/
func (bc *BlockChain) FindUTXO(publicHash []byte) []UTXOInfo {
	var outputs []UTXOInfo
	spentMap := make(map[string][]int64)
	//1.遍历区块链
	it := NewIterator(bc)
	for {
		block := it.Next()
		//2.遍历交易
		for _, tx := range block.Transactions {
		LABEL:
			//3.遍历outputs
			for outPutIndex, outPut := range tx.TXOutPuts {
				if bytes.Equal(outPut.PublicHash, publicHash) {
					currTxid := string(tx.Txid)
					indexArr := spentMap[currTxid]
					if len(indexArr) != 0 {
						for _, v := range indexArr {
							if v == int64(outPutIndex) {
								continue LABEL
							}
						}
					}
					txoInfo := UTXOInfo{
						output: outPut,
						intdex: int64(outPutIndex),
						txid:   tx.Txid,
					}
					fmt.Printf("OutPut:%x,Index:%d,Value:%f\n",publicHash,outPutIndex,outPut.Value)
					outputs = append(outputs, txoInfo)
				}
			}
			//4.如果不是挖矿交易
			if !tx.isCoinBaseTx() {
				for _, input := range tx.TXInPuts {
					pkHash := getPublicKeySignFromPublic(publicHash)
					if bytes.Equal(input.Public, pkHash) {
						spentKey := string(input.TXID)
						spentMap[spentKey] = append(spentMap[spentKey], input.Index)
					}
				}
			}
		}
		if block.PrevHash == nil {
			break
		}
	}
	return outputs
}

/*
	遍历账本，返回可以支配的钱
*/
func (bc *BlockChain) FindNeedUTXO(publicHash []byte, amount float64) ([]UTXOInfo, float64) {
	txos := bc.FindUTXO(publicHash)
	var sm float64
	var utxos []UTXOInfo
	for _, txo := range txos {
		sm += txo.output.Value
		utxos = append(utxos, txo)
		if sm >= amount {
			break
		}
	}
	return utxos, sm
}

/*
	签名
 */
func (bc *BlockChain)SignTransaction(privateKey *ecdsa.PrivateKey,tx *Transaction) bool {
	fmt.Println("签名:SignTransaction...")
	if tx.isCoinBaseTx() {
		fmt.Println("发现挖矿交易，不需要签名")
		return true
	}
	//1.查到tx所引用的交易的集合
	prevTx := make(map[string]*Transaction)
	//2.遍历
	for _,input := range tx.TXInPuts{
		tx := bc.FindTransactionByTxid(input.TXID)
		if tx == nil {
			return false
		}
		//将交易放到集合
		fmt.Printf("签名交易：txid: %x\n", tx.Txid)
		prevTx[string(input.TXID)] = tx
	}
	return tx.Sign(privateKey,prevTx)
}

/*
	验证
 */
func (bc *BlockChain)VerifyTransaction(tx *Transaction) bool {
	fmt.Println("校验：VerifyTransaction")
	if tx.isCoinBaseTx() {
		return true
	}
	prevtx := make(map[string]*Transaction)
	for _,input := range tx.TXInPuts{
		tx := bc.FindTransactionByTxid(input.TXID)
		if tx ==nil {
			return false
		}
		prevtx[string(input.TXID)] = tx
	}
	return tx.Verify(prevtx)
}

/*
	根据id找到交易本身
 */
func (bc *BlockChain)FindTransactionByTxid(TXid []byte) *Transaction {
	it := NewIterator(bc)
	for {
		block :=it.Next()
		for _,tx := range block.Transactions{
			if bytes.Equal(tx.Txid,TXid) {
				return tx
			}
		}
		if block.PrevHash == nil {
			break
		}
	}
	return nil
}
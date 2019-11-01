package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

/*
	交易结构
*/

type Transaction struct {
	Txid      []byte     //交易ID
	TXInPuts  []TXInPut  //多个交易输入
	TXOutPuts []TXOutPut //多个交易输出
	TimeStamp int64      //时间戳
}

/*
	交易输入
*/
type TXInPut struct {
	TXID      []byte //引用output所在的交易ID
	Index     int64  //引用output所在的下标
	ScriptSig string //解锁脚本
}

/*
	交易输出
*/
type TXOutPut struct {
	LockScript string  //锁定脚本
	Value      float64 //转账金额
}

//挖矿奖励
const reward = 12.5

/*
	挖矿交易
*/
// @ minner:矿池名称
// @ data:创世语等
func CoinBaseTX(minner string, data string) *Transaction {
	input := []TXInPut{
		TXInPut{
			TXID:      nil,
			Index:     -1,
			ScriptSig: data,
		},
	}
	output := []TXOutPut{
		TXOutPut{
			LockScript: minner,
			Value:      reward,
		},
	}
	tx := &Transaction{
		Txid:      nil,
		TXInPuts:  input,
		TXOutPuts: output,
		TimeStamp: time.Now().Unix(),
	}
	tx.SetTXID() //设置交易ID
	return tx
}

/*
	获取交易ID
*/
func (tx *Transaction) SetTXID() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		return
	}
	hash := sha256.Sum256(buff.Bytes())
	tx.Txid = hash[:]
}

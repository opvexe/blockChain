package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
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
	ScriptSig []byte //私钥签名
	Public []byte	//公钥
}

/*
	交易输出
*/
type TXOutPut struct {
	PublicHash []byte  //公钥签名哈希
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
	inputs := []TXInPut{
		TXInPut{
			TXID:      nil,
			Index:     -1,
			ScriptSig: []byte(data),
			Public:nil,
		},
	}
	outputs := []TXOutPut{NewTxOutPut(minner,reward)}
	tx := &Transaction{
		Txid:      nil,
		TXInPuts:  inputs,
		TXOutPuts: outputs,
		TimeStamp: time.Now().Unix(),
	}
	tx.SetTxid() //设置交易ID
	return tx
}

/*
	普通交易
*/
func NewTransaction(from, to string, amount float64, bc *BlockChain) (*Transaction, error) {
	//1.打开钱包
	wm :=NewWalletManager()
	if wm == nil {
		return nil,errors.New("打开钱包失败")
	}
	//2.根据付款人获取钱包
	w,ok :=wm.Wallets[from]
	if !ok {
		return nil,fmt.Errorf("未获取到%s的钱包\n",from)
	}
	//3.获取公钥私钥
	//privateKey :=w.PrivateKey
	publicKey := w.PublicKey
	//4.获取付款人的公钥哈希
	pubicHash :=getPublicKeySignFromPublic(publicKey)
	//5.找到付款人能被合理支配的钱
	ufos, value := bc.FindNeedUTXO(pubicHash, amount)
	if value < amount {
		return nil, errors.New("金额不足")
	}
	//拼接交易
	var inputs []TXInPut
	var outputs []TXOutPut
	//拼接input
	for _, utx := range ufos {
		input := TXInPut{
			TXID:      utx.txid,
			Index:     utx.intdex,
			ScriptSig: nil,
			Public:nil,
		}
		inputs = append(inputs, input)
	}
	//拼接output
	output := NewTxOutPut(to,amount)
	outputs = append(outputs, output)
	//判断是否需要找零,如果有零不写的话，会当成手续费
	if value > amount {
		output1 := NewTxOutPut(from,value-amount)
		outputs = append(outputs, output1)
	}
	tx := &Transaction{
		Txid:      nil,
		TXInPuts:  inputs,
		TXOutPuts: outputs,
		TimeStamp: time.Now().Unix(),
	}

	return tx, nil
}

/*
	创建output
 */
func NewTxOutPut(address string,amount float64) TXOutPut {
	//1.根据地址获取公钥签名
	publicKey :=getPublicKeySignFromAddress(address)
	output :=TXOutPut{
		PublicHash: publicKey,
		Value:      amount,
	}
	return output
}

/*
	设置交易id
*/
func (tx *Transaction) SetTxid() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		fmt.Println("SetTxid Encode err", err)
		return
	}
	hash := sha256.Sum256(buff.Bytes())
	tx.Txid = hash[:]
}

/*
	判断是否是挖矿交易
*/
func (tx *Transaction) isCoinBaseTx() bool {
	input := tx.TXInPuts[0]
	if len(tx.TXInPuts) == 1 && input.TXID == nil && input.Index == -1 {
		return true
	}
	return false
}

package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"math/big"
	"strings"
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
			Public:publicKey,
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
	拷贝副本
 */
func (tx *Transaction)TrimmedTransactionCopy() *Transaction {
	var inputs []TXInPut
	var outputs []TXOutPut
	for _,input := range tx.TXInPuts{
		inputNew := TXInPut{
			TXID:      input.TXID,
			Index:     input.Index,
			ScriptSig: nil,
			Public:    nil,
		}
		inputs = append(inputs, inputNew)
	}
	copy(outputs,tx.TXOutPuts)
	txCopy := &Transaction{
		Txid:      nil,
		TXInPuts:  inputs,
		TXOutPuts: outputs,
		TimeStamp: tx.TimeStamp,
	}
	return txCopy
}
/*
	签名
 */
func (tx *Transaction)Sign(privateKey *ecdsa.PrivateKey,prevTx map[string]*Transaction) bool {
	fmt.Println("开始签名:Sign...")
	//1.获取交易副本
	txCopy := tx.TrimmedTransactionCopy()
	//2.遍历副本
	for i,input := range txCopy.TXInPuts{
		prevtx :=prevTx[string(input.TXID)]
		if prevtx == nil {
			return false
		}
		output :=prevtx.TXOutPuts[input.Index]
		txCopy.TXInPuts[i].Public = output.PublicHash
		txCopy.SetTxid()
		//3.对当前交易做哈希处理，得到需要签名的数据
		hashData := txCopy.Txid
		fmt.Printf(">>>>>签名内容:%x\n",hashData)
		//5.使用私钥进行签名
		r,s,err:=ecdsa.Sign(rand.Reader,privateKey,hashData[:])
		if err!=nil {
			return false
		}
		//6.将签名赋值给原始交易
		signature := append(r.Bytes(),s.Bytes()...)
		tx.TXInPuts[i].ScriptSig = signature
		//7.将当前的input的public字段设置成nil
		txCopy.TXInPuts[i].Public = nil
		txCopy.Txid = nil
	}
	fmt.Println("交易签名成功")
	return true
}

/*
	验证
 */
func (tx *Transaction)Verify(prevTx map[string]*Transaction) bool {
	fmt.Println("开始验证:Verify...")
	//1.生成副本
	txCopy := tx.TrimmedTransactionCopy()
	//2.遍历副本
	for i,input := range txCopy.TXInPuts{
		prevtx := prevTx[string(input.TXID)]
		if prevtx == nil {
			return false
		}
		//2.对交易做哈希处理
		output :=prevtx.TXOutPuts[input.Index]
		txCopy.TXInPuts[i].Public = output.PublicHash
		txCopy.SetTxid()
		//3.使用签名，公钥，数据，进行校验
		hashData := txCopy.Txid //数据
		signData :=tx.TXInPuts[i].ScriptSig
		publicKey :=tx.TXInPuts[i].Public

		fmt.Printf("===>>>> 校验哈希:%x\n",hashData)

		//还原签名,r,s
		var r,s big.Int
		r.SetBytes(signData[:len(signData)/2])
		s.SetBytes(signData[len(signData)/2:])

		//还原公钥
		var x,y big.Int
		x.SetBytes(publicKey[:len(publicKey)/2])
		y.SetBytes(publicKey[len(publicKey)/2:])

		publicRaw := ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     &x,
			Y:     &y,
		}
		//验证
		if !ecdsa.Verify(&publicRaw,hashData[:],&r,&s) {
			fmt.Println("==>>>校验失败")
			return false
		}
		txCopy.Txid = nil
		txCopy.TXInPuts[i].Public = nil
	}
	fmt.Println("校验成功")
	return true
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

/*
	打印
 */
func (tx *Transaction)String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Txid))

	for i, input := range tx.TXInPuts {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Public))
	}

	for i, output := range tx.TXOutPuts{
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PublicHash))
	}

	return strings.Join(lines, "\n")
}
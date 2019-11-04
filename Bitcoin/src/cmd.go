package main

import (
	"fmt"
	"os"
	"strconv"
)

/*
	使用说明描述
*/
const DesUse = `	
	./bc add <data>  "区块数据"
	./bc print       "打印区块"
	./bc balance [from]	 "总金额"
	./bc send [from] [to] [amount] [minner] [data] "转账给某人"
	./bc create 	"创建钱包"
	./bc list 		"查询所有钱包地址"
`

type Cmd struct {
	bc *BlockChain
}

/*
	初始化
*/
func NewCmd(bc *BlockChain) *Cmd {
	return &Cmd{bc: bc}
}

/*
	解析
*/
func (c *Cmd) Run() {
	cmd := os.Args
	if len(cmd) < 2 {
		fmt.Println(DesUse)
		return
	}
	switch cmd[1] {
	case "add":
		if len(cmd) != 3 {
			fmt.Println(DesUse)
			return
		}
		c.add(cmd[2])
	case "send":
		if len(cmd) != 7 {
			fmt.Println("cmd 参数无效!")
			return
		}
		from := cmd[2]
		to := cmd[3]
		amountStr := cmd[4]
		amount, _ := strconv.ParseFloat(amountStr, 64)
		minner := cmd[5]
		data := cmd[6]
		c.send(from, to, minner, data, amount)
	case "print":
		c.print()
	case "balance":
		if len(cmd) != 3 {
			fmt.Println(DesUse)
			return
		}
		c.getBalance(cmd[2])
	case "list":
		c.listAddress()
	case "create":
		c.createWallet()
	default:
		fmt.Println(DesUse)
	}
}

/*
	添加
*/
func (c *Cmd) add(d string) {
	//c.bc.Add([]byte(d))
}

/*
	获取总额
*/
func (c *Cmd) getBalance(address string) {
	//1.根据地址获取公钥哈希
	publicHash := getPublicKeySignFromAddress(address)
	//2.获取未使用的钱
	utxo := c.bc.FindUTXO(publicHash)
	var t float64
	for _, txo := range utxo {
		t += txo.output.Value
	}
	fmt.Printf("%x总金额%f\n", publicHash, t)
}

/*
	打印
*/
func (c *Cmd) print() {
	it := NewIterator(c.bc)
	for {
		v := it.Next()
		fmt.Println("Version:", v.Version)
		fmt.Printf("PreHash:%x\n", v.PrevHash)
		fmt.Printf("Hash:%x\n", v.Hash)
		fmt.Println("MerkleRoot:", string(v.MerkleRoot))
		fmt.Println("TimeStamp:", v.TimeStamp)
		fmt.Println("Bits:", v.Bits)
		fmt.Println("Nonce:", v.Nonce)
		fmt.Println("Data:", v.Transactions[0].TXInPuts[0].ScriptSig)
		pow := NewProofWork(v)
		fmt.Printf("校验工作量:%v\n", pow.isValid())
		if v.PrevHash == nil {
			fmt.Println("区块链扫描完毕")
			break
		}
		fmt.Println("")
	}
}

/*
	打印所有地址
*/
func (c *Cmd) listAddress() {
	wm := NewWalletManager()
	if wm == nil {
		return
	}
	for _, value := range wm.listAddress() {
		fmt.Println("address:", value)
	}
}

/*
	创建钱包
*/
func (c *Cmd) createWallet() {
	wm := NewWalletManager()
	if wm == nil {
		return
	}
	address, err := wm.CreateWallet()
	if err != nil {
		return
	}
	fmt.Println("钱包地址:", address)
}

/*
	转账
*/
func (c *Cmd) send(from, to, minner, data string, amount float64) {
	fmt.Printf("%x转账给%x金额%f,挖矿人:%x,矿石语:%s\n", from, to, amount, minner, data)
	//1.创建挖矿交易
	coin := CoinBaseTX(minner, data)
	//2.创建交易
	txs := []*Transaction{coin}
	//3.创建普通交易
	otx, err := NewTransaction(from, to, amount, c.bc)
	if err != nil {
		fmt.Println("cmd send NewTransaction:", err)
		return
	} else {
		txs = append(txs, otx)
	}
	//添加到区块链上
	c.bc.Add(txs)
}

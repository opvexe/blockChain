package main

import (
	"fmt"
	"os"
)

/*
	使用说明描述
*/
const DesUse = `	
	./bc add <data>  "区块数据"
	./bc print       "打印区块"`

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
		}
		c.add(cmd[2])
	case "print":
		c.print()
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

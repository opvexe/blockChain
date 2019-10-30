package model

import (
	"bytes"
	"crypto/sha256"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"math/big"
)

/*
	工作量计算
*/
type ProofWork struct {
	block  *Block   //Block块
	target *big.Int //难度值
}

/*
	初始化
*/
func NewProofWork(block *Block) *ProofWork {
	pow := &ProofWork{
		block: block,
	}
	//系统自动调节获取一个难度的哈希值
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	bigTmp := big.Int{}
	bigTmp.SetString(targetStr, 16)
	pow.target = &bigTmp
	return pow
}

/*
	准备函数
*/
func (pow *ProofWork) prepareData(n int64) []byte {
	b := pow.block
	//获取加密数据
	hashData := [][]byte{
		[]byte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		dig2Byte(b.TimeStamp),
		dig2Byte(b.Bits),
		dig2Byte(n),
	}
	return bytes.Join(hashData, []byte(""))
}

/*
	开始挖矿
*/
func (pow *ProofWork) Run() (int64,[]byte) {
	var (
		n int64    //随机数
		h [32]byte //hash值
	)
	for{
		fmt.Printf("挖矿中:%x\n",h)
		data := pow.prepareData(n)
		h = sha256.Sum256(data)
		t := big.Int{}
		t.SetBytes(h[:])
		//比较
		if t.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功:%x,随机数:%d\n",h,n)
			break
		}else {
			n++
		}
	}
	return n,h[:]
}

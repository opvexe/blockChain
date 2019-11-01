package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"time"
)

/*
	区块
*/
type Block struct {
	Version      string         //版本号
	PrevHash     []byte         //前区块哈希值
	Hash         []byte         //当前哈希值
	MerkleRoot   []byte         //Merkle根，该区块交易的数据计算
	TimeStamp    int64          //时间戳
	Bits         int64          //难度值
	Nonce        int64          //随机数
	Transactions []*Transaction //交易结构体
}

/*
	初始化
*/
func NewBlock(pvHash []byte, txs []*Transaction) *Block {
	b := &Block{
		Version:      "Bitcoin 1.3.0.0",
		PrevHash:     pvHash,
		Hash:         nil,
		MerkleRoot:   nil,
		TimeStamp:    time.Now().Unix(),
		Bits:         0,
		Nonce:        0,
		Transactions: txs,
	}
	//挖矿
	pow := NewProofWork(b)
	n, h := pow.Run()
	b.Nonce = n
	b.Hash = h
	return b
}

/*
	编码
*/
func (b *Block) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(b)
	if err != nil {
		return nil
	}
	return buff.Bytes()
}

/*
	解码
*/
func Deserialize(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(block)
	if err != nil {
		return nil
	}
	return &block
}

/*
	将时间戳转换成byte
*/
func dig2Byte(i int64) []byte {
	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, i)
	if err != nil {
		return nil
	}
	return buff.Bytes()
}

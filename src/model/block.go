package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

/*
	区块
*/
type Block struct {
	version    string //版本号
	prevHash   []byte //前区块哈希值
	hash       []byte //当前哈希值
	merkleRoot []byte //Merkle根，该区块交易的数据计算
	timeStamp  int64  //时间戳
	bits       int64  //难度值
	nonce      int64  //随机数
	data       []byte //区块体
}

/*
	初始化
*/
func NewBlock(pvHash, data []byte) *Block {
	b := &Block{
		version:    "1.0",
		prevHash:   pvHash,
		hash:       nil,
		merkleRoot: nil,
		timeStamp:  time.Now().Unix(),
		bits:       0,
		nonce:      0,
		data:       data,
	}
	//	处理挖矿数据
	by := [][]byte{
		[]byte(b.version),
		b.prevHash,
		b.merkleRoot,
		dig2Byte(b.timeStamp),
		dig2Byte(b.bits),
		dig2Byte(b.nonce),
	}
	str := bytes.Join(by, []byte(""))
	hash := sha256.Sum256(str)
	b.hash = hash[:]
	return b
}

/*
	将时间戳转换成byte
*/
func dig2Byte(i int64) []byte {
	var buff bytes.Buffer
	err:=binary.Write(&buff,binary.LittleEndian,i)
	if err!=nil {
		return nil
	}
	return buff.Bytes()
}

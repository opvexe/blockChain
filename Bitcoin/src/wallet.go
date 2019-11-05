package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const payloadLen = 25 //base58解码后的长度
const checkSumLen = 4 //checksum的长度
/*
	钱包
*/
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey //钱包私钥
	PublicKey  []byte            //钱包公钥
}

/*
	创建钱包
*/
func NewWallet() *Wallet {
	//椭圆曲线加密
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("NewWallet GenerateKey :", err)
		return nil
	}
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

/*
	创建签名地址
*/
func (wm *Wallet) getAddress() string {
	//1.获取公钥哈希
	publicHash := getPublicKeySignFromPublic(wm.PublicKey)
	//2.获取版本号
	payload := append([]byte{byte(00)}, publicHash...)
	//3.截取校验码
	chesum := chekSum(payload)
	payload = append(payload, chesum...)
	return base58.Encode(payload)
}

/*
	根据公钥获取公钥签名
*/
func getPublicKeySignFromPublic(public []byte) []byte {
	//1.sha256加密
	hash256 := sha256.Sum256(public)
	//2.ripemd160加密
	ripemd := ripemd160.New()
	_, err := ripemd.Write(hash256[:])
	if err != nil {
		fmt.Println("getPublicKeySignFromPublic ripemd160:", err)
		return nil
	}
	ripemdHash := ripemd.Sum(nil)
	return ripemdHash
}

/*
	根据地址获取公钥签名
*/
func getPublicKeySignFromAddress(address string) []byte {
	//1.base58解密
	decode58 := base58.Decode(address)
	if len(decode58) != payloadLen {
		fmt.Println("getPublicKeySignFromAddress base58 Decode")
		return nil
	}
	//获取公钥哈希
	publicHash := decode58[1 : payloadLen-checkSumLen]
	return publicHash
}

/*
	获取校验码
*/
func chekSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]
	return checksum
}
/*
	校验地址的有效性
 */
func isAvailAddress(address string) bool {
	//1.对传入的地址base58解密
	decoderInfo := base58.Decode(address)
	//2.校验
	if len(decoderInfo)!=25 {
		fmt.Println("==>>>>地址无效:isAvailAddress")
		return false
	}
	checksum1 := chekSum(decoderInfo[:payloadLen-checkSumLen])
	checksum2:=decoderInfo[payloadLen-checkSumLen:]
	fmt.Printf("checksum1:%x\n",checksum1)
	fmt.Printf("checksum2:%x\n",checksum2)
	return bytes.Equal(checksum1,checksum2)
}
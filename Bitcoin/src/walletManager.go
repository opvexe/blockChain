package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

/*
	钱包管理器
*/
type WalletManager struct {
	Wallets map[string]*Wallet
}

/*
	创建结构
*/
func NewWalletManager() *WalletManager {
	var wm WalletManager
	wm.Wallets = make(map[string]*Wallet)
	err := wm.LoadFileName()
	if err != nil {
		return nil
	}
	return &wm
}

/*
	创建钱包
*/
func (wm *WalletManager) CreateWallet() (string, error) {
	//1.创建钱包
	w := NewWallet()
	if w == nil {
		return "", errors.New("WalletManager CreateWallet 钱包创建失败")
	}
	//2.获取地址
	address := w.getAddress()
	wm.Wallets[address] = w
	err := wm.SaveToFile()
	if err != nil {
		return "", err
	}
	return address, nil
}

const FileName = "wallet.dat"

/*
	保存至文件
*/
func (wm *WalletManager) SaveToFile() error {
	//1.序列化
	gob.Register(elliptic.P256())
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(wm)
	if err != nil {
		fmt.Println("WalletManager Encode:", err)
		return err
	}
	//写入
	err = ioutil.WriteFile(FileName, buff.Bytes(), 0600)
	if err != nil {
		fmt.Println("WalletManager WriteFile:", err)
		return err
	}
	return nil
}

/*
	获取文件
*/
func (wm *WalletManager) LoadFileName() error {
	//1.判断文件是否存在
	if !isFileExits(FileName) {
		fmt.Println("WalletManager SaveToFile 不存在")
		return nil
	}
	//2.读取文件
	buff, err := ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("WalletManager ReadFile:", err)
		return err
	}
	//3.反序列化
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(buff))
	err = decoder.Decode(wm)
	if err != nil {
		fmt.Println("WalletManager Decode:", err)
		return err
	}
	return nil
}

/*
	判断文件是否存在
*/
func isFileExits(fileName string) bool {
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		return false
	}
	return true
}

/*
	获取所有地址
 */
func (wm *WalletManager)listAddress() (address []string) {
	for key,_ :=range wm.Wallets{
		address = append(address, key)
	}
	sort.Strings(address)
	return
}

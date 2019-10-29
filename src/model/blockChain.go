package model

/*
	区块链
*/
type BlockChain struct {
	Chain []*Block
}

/*
	初始化
*/
func NewBlockChain() *BlockChain {
	genesis := NewBlock(nil, []byte("I'm Nakamoto. This is my first creative chain."))
	return &BlockChain{Chain: []*Block{
		genesis,
	}}
}

/*
	添加
*/
func (bc *BlockChain) Add(data []byte) {
	//获取前一个区块的哈希值
	pvHash := bc.Chain[len(bc.Chain)-1].hash
	b := NewBlock(pvHash, data)
	bc.Chain = append(bc.Chain, b)
}

package blockchain

import (
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/cache"
)

type Chain struct {
	curBlock *protocol.Block	//当前记录的高度
	db *cache.DBCache
}

func New(db *cache.DBCache) *Chain {
	return &Chain{
		db:db,
	}
}


func (ts *Chain) load()  {
	//todo 从数据库加载Chain，如果没有数据，则初始化默认 Chain
	block ,_:=ts.db.GetLastBlock()
	if block != nil{
		ts.curBlock = block
		return
	}

	//初始化创世块
	zeroBlock:= &protocol.Block{
		ParentHash: "0x00000000",
		BlockNum:0,
		Txs:[]*protocol.Tx{},
		Difficulty:"0001",
		Nonce:block.Nonce,
		TimeStamp:block.TimeStamp,
	}
	DeferBlockMgt.Complete(zeroBlock)

	//db
	ts.db.SetLastBlock(zeroBlock)
	ts.db.AddBlock(zeroBlock)

}



func (ts *Chain) Run()  {
	ts.load()
}


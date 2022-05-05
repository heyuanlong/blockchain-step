package blockchain

import (
	"heyuanlong/blockchain-step/p2p"
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/cache"
	"time"
)

type Chain struct {
	curBlock *protocol.Block	//当前记录的高度
	db *cache.DBCache
	p2p p2p.P2pI
}

func New(db *cache.DBCache,p2p p2p.P2pI) *Chain {
	return &Chain{
		db:db,
		p2p:p2p,
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


func (ts *Chain) checkHeight()  {
	for {
		//todo 检测高度是否落后了
		//请求获取落后的块


		time.Sleep(time.Second*10)
	}
}

func (ts *Chain) Run()  {
	ts.load()

	//接收区块
	//接收的区块放入区块池
	//chain从区块池获取区块
	//挖矿

	go ts.checkHeight()
	ts.p2p.RegisterOnReceive("block",ts.msgOnRecv)


}


func (ts *Chain) msgOnRecv(MsgType string, msgBytes []byte, p *p2p.Peer){

}




























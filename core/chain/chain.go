package chain

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	kblock "heyuanlong/blockchain-step/core/block"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/p2p"
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/cache"
	"time"
)

type Chain struct {
	curBlock *protocol.Block	//当前记录的高度
	db *cache.DBCache
	p2p p2p.P2pI

	notifyNewBlock chan *protocol.Block
}

func New(db *cache.DBCache,p2p p2p.P2pI) *Chain {
	return &Chain{
		db:db,
		p2p:p2p,
		notifyNewBlock : make(chan *protocol.Block,50),
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
	kblock.DeferBlockMgt.Complete(zeroBlock)

	ts.curBlock = zeroBlock

	//db
	ts.db.SetLastBlock(zeroBlock)
	ts.db.AddBlock(zeroBlock)

}

func (ts *Chain) getLastBlock()  {
	msg := p2p.BroadcastMsg{
		MsgType:types.MSG_TYPE_GETLASTBLOCK,
		Msg: []byte{},
	}

	ts.p2p.Broadcast(&msg)
}

//func (ts *Chain) checkHeight()  {
//	for {
//		//todo 检测高度是否落后了
//		//请求获取落后的块
//
//
//		time.Sleep(time.Second*10)
//	}
//}

func (ts *Chain) Run()  {
	ts.load()
	ts.p2p.RegisterOnReceive("block",ts.msgOnRecv)
	ts.getLastBlock()

	//接收区块
	//接收的区块放入区块池
	//chain从区块池获取区块
	//挖矿

	go ts.readBlockPool()
	ts.dig()
}

//chain从区块池获取区块
func (ts *Chain) readBlockPool()  {
	t :=time.NewTicker(time.Second)

	for{
		select {
			case <- t.C:
				//todo 查找池里有没有高度+1的区块
				
			case block := <- ts.notifyNewBlock:
				ts.dealNewBlock(block)

		}
	}
}

func (ts *Chain) dealNewBlock(block *protocol.Block)  {
	if ts.curBlock.BlockNum != (block.BlockNum - 1){
		return
	}

	kblock.DeferBlockMgt.Complete(block)

	ts.curBlock = block

	//db
	ts.db.SetLastBlock(block)
	ts.db.AddBlock(block)

	//todo 通知挖矿reset
}

//挖矿
func (ts *Chain) dig()  {

}

//-------------------------------------------------------------------
//p2p 回调处理函数
func (ts *Chain) msgOnRecv(msgType string, msgBytes []byte, p *p2p.Peer){
	switch msgType {
	case types.MSG_TYPE_BLOCK:
		ts.msgDealBlock(msgBytes,p)
	}
}

func (ts *Chain) msgDealBlock(msgBytes []byte, p *p2p.Peer){
	//接收的区块放入区块池

	block := &protocol.Block{}
	err :=proto.Unmarshal(msgBytes,block)
	if err !=nil{
		log.Error(err)
	}

	err = kblock.DeferBlockMgt.AddToPool(block)
	if err !=nil{
		log.Error(err)
	}

	//todo 可能会卡住
	ts.notifyNewBlock <- block
}



























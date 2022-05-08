package chain

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/common"
	kblock "heyuanlong/blockchain-step/core/block"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/p2p"
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/cache"
	"time"
)

type Chain struct {
	curBlock *protocol.Block //当前记录的高度
	db       *cache.DBCache
	p2p      p2p.P2pI

	notifyNewBlock chan *protocol.Block
	notifyDig      chan struct{}
}

func New(db *cache.DBCache, p2p p2p.P2pI) *Chain {
	return &Chain{
		db:             db,
		p2p:            p2p,
		notifyNewBlock: make(chan *protocol.Block, 50),
		notifyDig:      make(chan struct{}, 500),
	}
}

func (ts *Chain) load() {
	//从数据库加载Chain，如果没有数据，则初始化默认 Chain
	block, _ := ts.db.GetLastBlock()
	if block != nil {
		ts.curBlock = block
		return
	}

	//初始化创世块
	zeroBlock := &protocol.Block{
		ParentHash: "0x00000000",
		BlockNum:   0,
		Txs:        []*protocol.Tx{},
		Difficulty: "0001",
		Nonce:      block.Nonce,
		TimeStamp:  block.TimeStamp,
	}
	kblock.DeferBlockMgt.Complete(zeroBlock)

	ts.curBlock = zeroBlock

	//db
	ts.db.SetLastBlock(zeroBlock)
	ts.db.AddBlock(zeroBlock)

}

func (ts *Chain) getLastBlock() {
	msg := p2p.BroadcastMsg{
		MsgType: types.MSG_TYPE_GETLASTBLOCK,
		Msg:     []byte{},
	}

	ts.p2p.Broadcast(&msg)
}

func (ts *Chain) Run() {
	ts.load()
	ts.p2p.RegisterOnReceive("block", ts.msgOnRecv)
	ts.getLastBlock()

	//接收区块
	//接收的区块放入区块池
	//chain从区块池获取区块
	//挖矿

	go ts.readBlockPool()
	ts.digRun()
}

//chain从区块池获取区块
func (ts *Chain) readBlockPool() {
	t := time.NewTicker(time.Second)

	for {
		select {
		case <-t.C:
			//查找池里有没有高度+1的区块
			//todo 取区块需要修改，这里没有处理分叉的可能
			b := kblock.DeferBlockMgt.FindByNumber(ts.curBlock.BlockNum + 1)
			if b != nil {
				ts.dealNewBlock(b)
			}

		case block := <-ts.notifyNewBlock:
			ts.dealNewBlock(block)

		}
	}
}

func (ts *Chain) dealNewBlock(block *protocol.Block) {
	if ts.curBlock.BlockNum != (block.BlockNum - 1) {
		//todo 待处理分叉
		return
	}

	kblock.DeferBlockMgt.Complete(block)

	ts.curBlock = block

	//db
	ts.db.SetLastBlock(block)
	ts.db.AddBlock(block)

	//区块池移除
	kblock.DeferBlockMgt.DelFromPool(block)

	//通知挖矿reset
	ts.notifyDig <- struct{}{}
}

//挖矿
func (ts *Chain) digRun() {
	digBlock := ts.buildDigBlock()
	for {
		select {
		case <-ts.notifyDig:
			digBlock = ts.buildDigBlock()
		default:
			if ts.dig(digBlock) {
				//广播
				ts.broadcast(digBlock)

				//放入区块池
				ts.addToPool(digBlock)
			}

		}
	}
}

//构建待挖矿的区块
func (ts *Chain) buildDigBlock() *protocol.Block {
	block := &protocol.Block{
		ParentHash: ts.curBlock.Hash,
		BlockNum:   ts.curBlock.BlockNum + 1,
		Difficulty: "00",
		Nonce:      0,
		TimeStamp:  0,
	}

	//加载交易
	block.Txs = tx.DeferTxMgt.Gets(200)

	kblock.DeferBlockMgt.Complete(block)
	return block
}

//simple pow
func (ts *Chain) dig(block *protocol.Block) bool {
	t := time.Now()
	block.TimeStamp = uint64(t.Unix())
	n := uint64(t.UnixNano())
	for i := 0; i < 1000; i++ {
		block.Nonce = n + uint64(i)
		hash := kblock.DeferBlockMgt.Hash(block)
		if string(hash[0:2]) == block.Difficulty {
			block.Hash = common.Bytes2HexWithPrefix(hash)
			return true
		}
	}

	return false
}

//p2p 回调处理函数-------------------------------------------------------------------

func (ts *Chain) msgOnRecv(msgType string, msgBytes []byte, p *p2p.Peer) {
	switch msgType {
	case types.MSG_TYPE_BLOCK:
		ts.msgDealBlock(msgBytes, p)
	}
}

func (ts *Chain) msgDealBlock(msgBytes []byte, p *p2p.Peer) {
	//接收的区块放入区块池

	block := &protocol.Block{}
	err := proto.Unmarshal(msgBytes, block)
	if err != nil {
		log.Error(err)
	}

	//如果区块太旧，就丢弃
	if ts.curBlock.BlockNum > (block.BlockNum + 100){
		return
	}

	ts.addToPool(block)
}
//--------------------------------------------------------------------------------

//发送到区块池，并通知chain去取
func (ts *Chain) addToPool(block *protocol.Block) {
	err := kblock.DeferBlockMgt.AddToPool(block)
	if err != nil {
		log.Error(err)
	}

	//todo 可能会卡住
	ts.notifyNewBlock <- block
}

//
func (ts *Chain) broadcast(block *protocol.Block) {

	msg, err := proto.Marshal(block)
	if err != nil {
		log.Error(err)
	}

	m := &p2p.BroadcastMsg{
		types.MSG_TYPE_BLOCK,
		msg,
	}

	ts.p2p.Broadcast(m)
}

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

	notifyHaveBlockToPool chan *protocol.Block //通知 goroutine 有新块到了区块池
	notifyDig      chan struct{}			   //通知 goroutine 重置挖矿数据

	peerBlockHeight uint64				//对等网络的标记高度
	peerHeightMap map[string]uint64		//对等网络的高度

}

func New(db *cache.DBCache, p2p p2p.P2pI) *Chain {
	return &Chain{
		db:             db,
		p2p:            p2p,
		notifyHaveBlockToPool: make(chan *protocol.Block, 50),
		notifyDig:      make(chan struct{}, 500),

		peerHeightMap:make(map[string]uint64),
	}
}

func (ts *Chain) Run() {
	//加载Chain数据
	ts.loadChain()

	//注册p2p数据回调函数
	ts.p2p.RegisterOnReceive(types.MSG_TYPE_BLOCK, ts.msgOnRecv)
	ts.p2p.RegisterOnReceive(types.MSG_TYPE_RESP_LASTBLOCK, ts.msgOnRecv)

	//广播请求获取对等节点的块
	go ts.getLastBlock()

	//从区块池取区块
	go ts.readBlockPool()

	//挖矿
	ts.digRun()
}


func (ts *Chain) loadChain() {
	//从数据库加载Chain数据
	block, _ := ts.db.GetLastBlock()
	if block != nil {
		ts.curBlock = block
		return
	}

	//没有Chain数据
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
	ts.commit(zeroBlock)

}

//广播请求获取对等节点的最新块
func (ts *Chain) getLastBlock() {
	msg := p2p.BroadcastMsg{
		MsgType: types.MSG_TYPE_REQ_LASTBLOCK,
		Msg:     []byte{},
	}

	t := time.NewTicker(time.Second * 1)

	for {
		ts.p2p.Broadcast(&msg)

		//发现落后于对等节点，请求block
		//todo 极有可能多次重复广播请求同样的区块号的区块
		if ts.curBlock.BlockNum  <  (ts.peerBlockHeight) {
			peers , _ := ts.p2p.Peers()
			peersLen := len(peers)
			if peersLen == 0 {
				log.Error("peers len == 0")
				return
			}

			peerIndex := 0
			for i := ts.curBlock.BlockNum + 1; i < ts.peerBlockHeight; i++ {

				//广播请求指定区块
				ts.broadcastReqBlock(i,peers[peerIndex])

				peerIndex++
				if peerIndex == peersLen{
					peerIndex = 0
				}
			}
		}

		<-t.C
	}
}

//从区块池获取区块
func (ts *Chain) readBlockPool() {
	t := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-t.C:
			ts.dealNewBlock()

		case  <-ts.notifyHaveBlockToPool:
			ts.dealNewBlock()

		}
	}
}

func (ts *Chain) dealNewBlock() {

	block := kblock.DeferBlockMgt.GetFisrt()
	if block == nil {
		return
	}

	//todo 待处理分叉

	if ts.curBlock.BlockNum >=  (block.BlockNum) {
		return
	}

	kblock.DeferBlockMgt.Complete(block)

	ts.curBlock = block
	ts.commit(block)

	//区块池移除
	kblock.DeferBlockMgt.DelFromPool(block)

	//通知挖矿reset
	ts.notifyDig <- struct{}{}
}

func (ts *Chain) commit(block *protocol.Block) {
	//db
	ts.db.SetLastBlock(block)
	ts.db.AddBlock(block)
}


//挖矿----------------------------------------------------------------------------
func (ts *Chain) digRun() {
	digBlock := ts.buildDigBlock()
	for {
		select {
		case <-ts.notifyDig:
			digBlock = ts.buildDigBlock()
		default:
			if ts.dig(digBlock) {
				//广播挖到的区块
				ts.broadcastDigBlock(digBlock)

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
	case types.MSG_TYPE_RESP_LASTBLOCK:
		ts.msgDealRespLastBlock(msgBytes, p)
	}
}


//block
func (ts *Chain) msgDealBlock(msgBytes []byte, p *p2p.Peer) {
	//接收的区块放入区块池

	block := &protocol.Block{}
	err := proto.Unmarshal(msgBytes, block)
	if err != nil {
		log.Error(err)
	}

	//todo 要多个分叉里判断
	//如果区块太旧，就丢弃
	if ts.curBlock.BlockNum > (block.BlockNum + 100){
		return
	}

	ts.addToPool(block)
}

//respLastBlock
func (ts *Chain) msgDealRespLastBlock(msgBytes []byte, p *p2p.Peer) {

	block := &protocol.Block{}
	err := proto.Unmarshal(msgBytes, block)
	if err != nil {
		log.Error(err)
	}

	ts.peerHeightMap[p.ID] = block.BlockNum

	//todo 计算标记高度

	ts.peerBlockHeight = block.BlockNum
}

//--------------------------------------------------------------------------------

//发送到区块池，并通知chain去取
func (ts *Chain) addToPool(block *protocol.Block) {
	err := kblock.DeferBlockMgt.AddToPool(block)
	if err != nil {
		log.Error(err)
	}

	//todo 可能会卡住
	ts.notifyHaveBlockToPool <- block
}

//广播挖到的区块
func (ts *Chain) broadcastDigBlock(block *protocol.Block) {

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

//广播请求指定区块
func (ts *Chain) broadcastReqBlock(number uint64,p *p2p.Peer) {

	obj := &protocol.BlockNumber{
		BlockNum:   number,
	}
	msg, err := proto.Marshal(obj)
	if err != nil {
		log.Error(err)
	}

	m := &p2p.BroadcastMsg{
		types.MSG_TYPE_REQ_BLOCKBYNUMBER,
		msg,
	}

	ts.p2p.BroadcastToPeer(m,p)
}

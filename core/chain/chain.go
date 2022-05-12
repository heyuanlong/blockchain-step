package chain

import (
	"errors"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/common"
	kblock "heyuanlong/blockchain-step/core/block"
	"heyuanlong/blockchain-step/core/config"
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
	blockMgt *kblock.BlockMgt

	notifyHaveBlockToPool chan *protocol.Block //通知 goroutine 有新块到了区块池
	notifyDig      chan struct{}			   //通知 goroutine 重置挖矿数据

	peerBlockHeight uint64				//对等网络的标记高度
	peerHeightMap map[string]uint64		//对等网络的高度

}

func New(db *cache.DBCache, p2p p2p.P2pI,blockMgt *kblock.BlockMgt) *Chain {
	return &Chain{
		db:             db,
		p2p:            p2p,
		blockMgt:            blockMgt,

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
		Difficulty: "000",
		Nonce:      0,
		TimeStamp:  uint64(time.Now().Unix()),
	}
	ts.blockMgt.Complete(zeroBlock)


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
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
			log.Error(err)
		}
	}()

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

	block := ts.blockMgt.GetFisrt()
	if block == nil {
		log.Error("not find new block")
		return
	}
	log.Info("dealNewBlock:",block.BlockNum)

	//todo 待处理分叉

	if ts.curBlock.BlockNum >=  (block.BlockNum) {
		log.Info("ts.curBlock.BlockNum >=  (block.BlockNum):",ts.curBlock.BlockNum , (block.BlockNum))
		//区块池移除
		ts.blockMgt.DelFromPool(block)
		return
	}

	ts.blockMgt.Complete(block)


	if err := ts.commit(block);err != nil{
		log.Error("commitBlock fail:",block.BlockNum,",err:",err)

		ts.blockMgt.DelFromPool(block)

	}else{
		log.Info("commitBlock:",block.BlockNum)

		ts.curBlock = block
		ts.blockMgt.DelFromPool(block)
	}





	//通知挖矿reset
	ts.notifyDig <- struct{}{}

}

func (ts *Chain) checkTx(txObj *protocol.Tx) error {
	senderObj ,err :=ts.db.GetAccount(txObj.Sender.Address)
	if err != nil{
		log.Error(err)
		return err
	}
	if senderObj.Id.Address == "" {
		log.Error("sender not find in chain block")
		return errors.New("sender not find in chain block")
	}
	if senderObj.Nonce != txObj.Nonce{
		log.Error("nonce 错误")
		return errors.New("nonce 错误")
	}
	if  txObj.Amount > senderObj.Balance {
		log.Error("发送的金额大于余额")
		return errors.New("发送的金额大于余额")
	}

	return nil
}

func (ts *Chain) commit(block *protocol.Block) error {

	for _, txObj := range block.Txs {
		if err := ts.checkTx(txObj);err!= nil{
			//txPool移除
			tx.DeferTxMgt.DelFromPool(txObj)
			return err
		}
	}

	//db
	if err := ts.db.SetLastBlock(block);err != nil{
		log.Info(err)
		return err
	}
	if err := ts.db.AddBlock(block);err != nil{
		log.Info(err)
		return err
	}

	//todo 快照与回滚
	for _, txObj := range block.Txs {

		//tx 数据提交
		senderAccount ,_ :=ts.db.GetAccount(txObj.Sender.Address)
		toAccount ,_ :=ts.db.GetAccount(txObj.To.Address)

		//说明账户不存在 需要新建一个账户
		if toAccount == nil{
			toAccount = &protocol.Account{
				Id:          txObj.To,
				Balance:     0,
				Nonce:     0,
				AccountType: int32(protocol.AccountType_Normal),
			}
		}

		senderAccount.Balance -= txObj.Amount
		senderAccount.Nonce += 1
		toAccount.Balance += txObj.Amount

		ts.db.AddAccount(senderAccount)
		ts.db.AddAccount(toAccount)


		//todo 交易收据验证

		//交易存储
		ts.db.AddTx(txObj)

		//txPool移除
		tx.DeferTxMgt.DelFromPool(txObj)
	}


	//挖矿奖励
	if block.Miner != nil && block.Miner.Address != ""{
		minerAccount ,_ :=ts.db.GetAccount(block.Miner.Address)
		if minerAccount == nil{
			minerAccount = &protocol.Account{
				Id:          block.Miner,
				Balance:     0,
				Nonce:     	 0,
				AccountType: int32(protocol.AccountType_Normal),
			}
		}
		minerAccount.Balance += types.DIG_REWARD
		ts.db.AddAccount(minerAccount)
	}


	return nil
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
		Difficulty: "000000",
		Nonce:      0,
		TimeStamp:  0,
		Miner: &protocol.Address{Address: config.Config.Miner},
	}

	block.Txs = tx.DeferTxMgt.Gets(200)
	log.Info("加载交易:",len(block.Txs))

	ts.blockMgt.Complete(block)
	return block
}

//simple pow
func (ts *Chain) dig(block *protocol.Block) bool {
	t := time.Now()
	block.TimeStamp = uint64(t.Unix())
	n := uint64(t.UnixNano())
	for i := 0; i < 1000; i++ {
		block.Nonce = n + uint64(i)
		hash := ts.blockMgt.Hash(block)
		hashHex := common.Bytes2Hex(hash)

		if hashHex[0:len(block.Difficulty)] == block.Difficulty {
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
	//todo 校验里面的tx
	err := ts.blockMgt.AddToPool(block)
	if err != nil {
		log.Error(err,block.BlockNum)
	}

	//todo 可能会卡住
	ts.notifyHaveBlockToPool <- block
}

//广播挖到的区块
func (ts *Chain) broadcastDigBlock(block *protocol.Block) {

	msg, err := proto.Marshal(block)
	if err != nil {
		log.Error(err)
		return
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

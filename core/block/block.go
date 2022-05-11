package block

import (
	"crypto/sha256"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/core/types"
	"heyuanlong/blockchain-step/p2p"
	"heyuanlong/blockchain-step/protocol"
	"google.golang.org/protobuf/proto"
	"heyuanlong/blockchain-step/storage/cache"
	"sync"
)



type BlockMgt struct {
	sync.RWMutex
	db *cache.DBCache
	p2p      p2p.P2pI

	blockPool *blockPoolStruct
}

func NewBlockMgt(db *cache.DBCache, p2p p2p.P2pI) *BlockMgt{
	return &BlockMgt{
		db:db,
		p2p:p2p,
		blockPool :newBlockPool(),
	}

}


func (ts *BlockMgt) Run() {
	//注册p2p数据回调函数
	ts.p2p.RegisterOnReceive(types.MSG_TYPE_REQ_BLOCKBYNUMBER, ts.msgOnRecv)
}


func (ts *BlockMgt) Complete(block *protocol.Block){

	//MerkleRoot
	block.TxsRoot = ts.MerkleRoot(block)

	//block.Hash
	block.Hash = common.Bytes2HexWithPrefix(ts.Hash(block))
}

func (ts *BlockMgt) Hash(block *protocol.Block) ([]byte) {
	t := &protocol.Block{
		ParentHash:block.ParentHash,
		Txs:block.Txs,
		Difficulty:block.Difficulty,
		Nonce:block.Nonce,
		TimeStamp:block.TimeStamp,
		TxsRoot:block.TxsRoot,
	}
	b, _ := proto.Marshal(t)

	sh := sha256.New()
	sh.Write(b)
	hash := sh.Sum(nil)

	return hash
}


func (ts *BlockMgt) MerkleRoot(block *protocol.Block) []byte {
	txHashs := make([][]byte, 0, len(block.Txs))
	for i := range block.Txs {
		hash ,_ := tx.DeferTxMgt.Hash(block.Txs[i])
		txHashs = append(txHashs,hash )
	}
	return common.Merkel(txHashs)
}

//-----------------------------------------------------------------

func (ts *BlockMgt) AddToPool(block *protocol.Block) error {
	return ts.blockPool.AddToPool(block)
}

func (ts *BlockMgt) DelFromPool(block *protocol.Block) error {
	return ts.blockPool.DelFromPool(block)
}
func (ts *BlockMgt) IsInPool(block *protocol.Block) bool {
	return ts.blockPool.IsInPool(block)
}
func (ts *BlockMgt) GetFisrt() *protocol.Block {
	return ts.blockPool.GetFisrt()
}


//p2p 回调处理函数-------------------------------------------------------------------

func (ts *BlockMgt) msgOnRecv(msgType string, msgBytes []byte, p *p2p.Peer) {
	switch msgType {
	case types.MSG_TYPE_REQ_BLOCKBYNUMBER:
		ts.msgDealReqBlockByNumber(msgBytes, p)
	}
}

//block
func (ts *BlockMgt) msgDealReqBlockByNumber(msgBytes []byte, p *p2p.Peer) {
	//接收的区块放入区块池

	blockNumber := &protocol.BlockNumber{}
	err := proto.Unmarshal(msgBytes, blockNumber)
	if err != nil {
		log.Error(err)
	}

	block,err := ts.db.GetBlockByNumber(blockNumber.BlockNum)
	if err != nil{
		log.Error(err)
	}
	if block == nil{
		return
	}

	msg, err := proto.Marshal(block)
	if err != nil {
		log.Error(err)
		return
	}
	m := &p2p.BroadcastMsg{
		types.MSG_TYPE_BLOCK,
		msg,
	}

	ts.p2p.BroadcastToPeer(m,p)
}
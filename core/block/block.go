package block

import (
	"crypto/sha256"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/protocol"
	"google.golang.org/protobuf/proto"
	"sync"
)

var DeferBlockMgt BlockMgt

type BlockMgt struct {
	sync.RWMutex
	poolCap      int
	blockPool map[string]*protocol.Block
}

func (ts *BlockMgt) Complete(block *protocol.Block){

	//MerkleRoot
	block.TxsRoot = ts.MerkleRoot(block)

	//block.Hash
	hash := ts.Hash(block)
	block.Hash = common.Bytes2HexWithPrefix(hash)
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

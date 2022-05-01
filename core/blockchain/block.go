package blockchain

import (
	"crypto/sha256"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/core/tx"
	"heyuanlong/blockchain-step/protocol"
	"google.golang.org/protobuf/proto"
)

var DeferBlockMgt BlockMgt

type BlockMgt struct {

}
func (ts *BlockMgt) Complete(block *protocol.Block){

	//MerkleRoot
	block.TxsRoot = ts.MerkleRoot(block)

	//block.Hash
	hash,_ := ts.Hash(block)
	block.Hash = common.Bytes2HexWithPrefix(hash)
}

func (ts *BlockMgt) Hash(block *protocol.Block) ([]byte,error) {
	t := &protocol.Block{
		ParentHash:block.ParentHash,
		Txs:block.Txs,
		Difficulty:block.Difficulty,
		Nonce:block.Nonce,
		TimeStamp:block.TimeStamp,
		TxsRoot:block.TxsRoot,
	}
	b, err := proto.Marshal(t)
	if err != nil {
		return []byte{}, err
	}
	sh := sha256.New()
	sh.Write(b)
	hash := sh.Sum(nil)

	return hash,nil
}


func (ts *BlockMgt) MerkleRoot(block *protocol.Block) []byte {
	txHashs := make([][]byte, 0, len(block.Txs))
	for i := range block.Txs {
		hash ,_ := tx.DeferTxMgt.Hash(block.Txs[i])
		txHashs = append(txHashs,hash )
	}
	return common.Merkel(txHashs)
}

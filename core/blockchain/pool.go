package blockchain

import (
	"errors"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/protocol"
)

func (ts *BlockMgt) AddToPool(block *protocol.Block) (error) {
	hash ,_ := ts.Hash(block)
	blockid := common.Bytes2HexWithPrefix(hash)

	ts.Lock()
	defer ts.Unlock()

	_, ok := ts.blockPool[blockid]
	if ok {
		return errors.New("existing")
	}

	if len(ts.blockPool) >= ts.poolCap{
		return errors.New("pool cap is full")
	}

	ts.blockPool[blockid] = block
	return nil
}

func (ts *BlockMgt) DelFromPool(block *protocol.Block) (error) {
	hash ,_ := ts.Hash(block)
	blockid := common.Bytes2HexWithPrefix(hash)

	ts.Lock()
	defer ts.Unlock()

	_, ok := ts.blockPool[blockid]
	if !ok {
		return errors.New("not find")
	}

	delete(ts.blockPool, blockid)

	return nil
}

func (ts *BlockMgt) IsInPool(block *protocol.Block) (bool) {
	ts.RLock()
	defer ts.RUnlock()

	hash ,_ := ts.Hash(block)
	blockid := common.Bytes2HexWithPrefix(hash)

	_, ok := ts.blockPool[blockid]
	return ok
}

// 从交易池获取一定数量的交易
func (ts *BlockMgt) Gets(num int) []*protocol.Block {
	ts.RLock()
	defer ts.RUnlock()

	blocks := make([]*protocol.Block, 0)
	i := 0
	for _, v := range ts.blockPool {
		blocks = append(blocks, v)
		i++
		if i>= num{
			break
		}
	}

	return blocks
}
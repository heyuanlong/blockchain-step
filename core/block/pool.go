package block

import (
	"errors"
	"heyuanlong/blockchain-step/protocol"
	"sort"
	"sync"
)

type blockSort []*protocol.Block

func (ts blockSort) Len() int           { return len(ts) }
func (ts blockSort) Less(i, j int) bool { return ts[i].BlockNum < ts[j].BlockNum }
func (ts blockSort) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }

type blockPoolStruct struct {
	sync.RWMutex
	poolCap      int
	blocks blockSort
}

func newBlockPool() *blockPoolStruct {
	return &blockPoolStruct{
		poolCap:1000, //todo 可能初次拉取块的时候，这里的容量不够大
		blocks: make([]*protocol.Block, 0),
	}
}

func (ts *blockPoolStruct) AddToPool(block *protocol.Block) error {
	ts.Lock()
	defer ts.Unlock()

	for _, v := range ts.blocks {
		if v.Hash == block.Hash {
			return errors.New("existing")
		}
	}

	if len(ts.blocks) >= ts.poolCap {
		return errors.New("pool cap is full")
	}

	ts.blocks = append(ts.blocks,block)
	sort.Sort(ts.blocks)

	return nil
}

func (ts *blockPoolStruct) DelFromPool(block *protocol.Block) error {
	ts.Lock()
	defer ts.Unlock()

	index := -1
	for i, v := range ts.blocks {
		if v.Hash == block.Hash {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("not find")
	}

	ts.blocks = append(ts.blocks[:index], ts.blocks[index+1:]...)
	return nil
}


func (ts *blockPoolStruct) IsInPool(block *protocol.Block) bool {
	ts.RLock()
	defer ts.RUnlock()

	for _, v := range ts.blocks {
		if v.Hash == block.Hash {
			return true
		}
	}

	return false
}

//
func (ts *blockPoolStruct) GetFisrt() *protocol.Block {
	ts.RLock()
	defer ts.RUnlock()

	if len(ts.blocks) > 0{
		return ts.blocks[0]
	}

	return nil
}

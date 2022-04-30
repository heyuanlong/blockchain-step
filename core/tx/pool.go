package tx

import (
	"errors"
	"heyuanlong/blockchain-step/common"
	"heyuanlong/blockchain-step/protocol"
)

func (ts *TxMgt) AddToPool(tx *protocol.Tx) (error) {
	hash ,_ := ts.Hash(tx)
	txid := common.Bytes2HexWithPrefix(hash)

	ts.Load()
	defer ts.Unlock()

	_, ok := ts.txPool[txid]
	if ok {
		return errors.New("existing")
	}

	if len(ts.txPool) >= ts.poolCap{
		return errors.New("pool cap is full")
	}

	ts.txPool[txid] = tx
	return nil
}

func (ts *TxMgt) DelFromPool(tx *protocol.Tx) (error) {
	hash ,_ := ts.Hash(tx)
	txid := common.Bytes2HexWithPrefix(hash)

	ts.Load()
	defer ts.Unlock()

	_, ok := ts.txPool[txid]
	if !ok {
		return errors.New("not find")
	}

	delete(ts.txPool, txid)

	return nil
}

func (ts *TxMgt) IsInPool(tx *protocol.Tx) (bool) {
	ts.RLock()
	defer ts.RUnlock()

	hash ,_ := ts.Hash(tx)
	txid := common.Bytes2HexWithPrefix(hash)

	_, ok := ts.txPool[txid]
	return ok
}

// 从交易池获取一定数量的交易
func (ts *TxMgt) Gets(num int) []*protocol.Tx {
	ts.RLock()
	defer ts.RUnlock()

	txs := make([]*protocol.Tx, 0)
	i := 0
	for _, v := range ts.txPool {
		txs = append(txs, v)
		i++
		if i>= num{
			break
		}
	}

	return txs
}
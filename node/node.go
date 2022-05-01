package node

import (
	"heyuanlong/blockchain-step/core/blockchain"
	"heyuanlong/blockchain-step/storage/cache"
	"os"
	"os/signal"
	"syscall"
)

// Node
type Node struct {
	db *cache.DBCache
	chain *blockchain.Chain
}


func New() *Node {
	// 创建缓存数据库
	db := cache.New("./.datadir")

	//chain
	chain:=blockchain.New(db)

	return &Node{
		db:db,
		chain:chain,
	}
}


func (ts *Node) Run() {
	go ts.chain.Run()



	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT) // 2,3,15
	<- ch

	//todo 清理
}
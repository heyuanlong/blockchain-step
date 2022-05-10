package node

import (
	kblock "heyuanlong/blockchain-step/core/block"
	chain2 "heyuanlong/blockchain-step/core/chain"
	"heyuanlong/blockchain-step/core/config"
	"heyuanlong/blockchain-step/p2p/http"
	"heyuanlong/blockchain-step/storage/cache"
	"os"
	"os/signal"
	"syscall"
)

// Node
type Node struct {
	db *cache.DBCache

	blockMgt *kblock.BlockMgt
	chain *chain2.Chain
}


func New(confFile string) *Node {
	config.InitConfParamByFile(confFile)
	// 创建缓存数据库
	db := cache.New(config.Config.DataDir)
	// p2p
	p:=http.New(config.Config.Node.NodeAddrs,config.Config.Node.Local,config.Config.Node.NodeId)
	//blockMgt
	blockMgt := kblock.NewBlockMgt(db,p)
	//chain
	chain:= chain2.New(db,p,blockMgt)

	return &Node{
		db:db,
		blockMgt:blockMgt,
		chain:chain,
	}
}


func (ts *Node) Run() {
	go ts.blockMgt.Run()
	go ts.chain.Run()



	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT) // 2,3,15
	<- ch

	//todo 清理
}





























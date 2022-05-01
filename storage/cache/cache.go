package cache

import (
	"errors"
	lru "github.com/hashicorp/golang-lru"
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/database"
	"path"
	"google.golang.org/protobuf/proto"
)

// 增加一个缓存层
type DBCache struct {
	blockCache *lru.Cache
	blockDB    database.DB

	blockNumberCache *lru.Cache
}

func New(filepath string) *DBCache {

	blockCache, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	blockDB, err := database.NewLevelDB(path.Join(filepath, "./stepblock/block.db"))
	if err != nil {
		panic(err)
	}

	blockNumberCache, err := lru.New(1024)
	if err != nil {
		panic(err)
	}


	dbCache := &DBCache{
		blockCache: blockCache,
		blockDB:    blockDB,
		blockNumberCache:    blockNumberCache,
	}

	return dbCache
}

func (ts *DBCache) GetLastBlock() (*protocol.Block,error) {
	v ,err := ts.blockDB.Get([]byte("last_block"))
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil , errors.New("not find")
	}

	var block protocol.Block
	err = proto.Unmarshal([]byte(v), &block)
	return &block, err
}

func (ts *DBCache) SetLastBlock(block *protocol.Block) (error) {
	v, _ := proto.Marshal(block)
	return ts.blockDB.Set([]byte("last_block"),v)
}

func (ts *DBCache) AddBlock(block *protocol.Block) (error) {
	ts.blockCache.Add(block.Hash,block)
	ts.blockNumberCache.Add(block.BlockNum,block)

	v, _ := proto.Marshal(block)
	return ts.blockDB.Set([]byte(block.Hash),v)
}



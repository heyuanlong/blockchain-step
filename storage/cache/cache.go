package cache

import (
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"heyuanlong/blockchain-step/protocol"
	"heyuanlong/blockchain-step/storage/database"
	"path"
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

//-----------------------------------------------------------------
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
	ts.blockNumberCache.Add(block.BlockNum,block.Hash)

	v, _ := proto.Marshal(block)
	if err :=  ts.blockDB.Set([]byte(block.Hash),v); err != nil {
		return err
	}
	if err :=  ts.blockDB.Set([]byte(fmt.Sprintf("%d",block.BlockNum)),[]byte(block.Hash)); err != nil {
		return err
	}


	return nil
}

func (ts *DBCache) GetBlockByHash(hash string) (*protocol.Block,error) {
	v,ok := ts.blockCache.Get(hash)
	if ok {
		return v.(*protocol.Block), nil
	}

	value ,err := ts.blockDB.Get([]byte(hash))
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, errors.New("GetBlockByHash:链数据有问题?")
	}

	var block protocol.Block
	err = proto.Unmarshal(value, &block)
	return &block, err
}


func (ts *DBCache) GetBlockByNumber(number uint64) (*protocol.Block,error) {
	hash , ok :=ts.blockNumberCache.Get(number)
	if !ok {
		value ,err := ts.blockDB.Get([]byte(fmt.Sprintf("%d",number)))
		if err != nil {
			return nil, err
		}
		if len(value) == 0 {
			return nil, nil
		}
		hash = value
	}

	log.Println(hash.(string))


	return ts.GetBlockByHash(hash.(string))
}


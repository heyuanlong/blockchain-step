package cache

import (
	"path"
	lru "github.com/hashicorp/golang-lru"
	"heyuanlong/blockchain-step/storage/database"
)



// 增加一个缓存层
type DBCache struct {

	blockCache      *lru.Cache
	blockDB database.DB

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


	dbCache := &DBCache{
		blockCache:blockCache,
		blockDB:blockDB,
	}

	return dbCache
}
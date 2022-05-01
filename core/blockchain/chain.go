package blockchain

type chain struct {
	blockCurHeight int64	//当前记录的高度
}


func (ts *chain) Load()  {
	//todo 从数据库加载chain，如果没有数据，则初始化默认 chain
}



func (ts *chain) Run()  {

}


package p2p

import "sync"

type P2pI interface {
	// 向所有的节点广播消息
	Broadcast( msg *BroadcastMsg) error
	// 广播到指定的peer
	BroadcastToPeer( msg *BroadcastMsg, p *Peer) error
	// 广播 除了指定的peer
	BroadcastExceptPeer( msg *BroadcastMsg, p *Peer) error

	// 移除某个peer
	RemovePeer(p *Peer) error
	RegisterOnReceive(MsgType string, callBack OnReceive) error
	Start() error
	// 返回所有存在的peers
	Peers() ([]*Peer, error)
}

type BroadcastMsg struct {
	MsgType string                 `json:"msg_type"`
	Msg     []byte                 `json:"msg"`
}

// OnReceive 注册接收消息回到
type OnReceive func(MsgType string, msgBytes []byte, p *Peer)

type Peer struct {
	ID      string // 定义peerid  每个peerid应该是唯一的
	Address string // 地址
}

type PeerBooks struct {
	sync.RWMutex
	sets map[string]*Peer
}

func NewPeerBooks() *PeerBooks {
	return &PeerBooks{
		sets: make(map[string]*Peer),
	}
}

func (pb *PeerBooks) AddPeer(p *Peer) {
	if p == nil {
		return
	}
	pb.Lock()
	pb.sets[p.ID] = p
	pb.Unlock()
}

func (pb *PeerBooks) FindPeer(id string) *Peer {
	pb.RLock()
	defer pb.RUnlock()
	v := pb.sets[id]
	return v
}

func (pb *PeerBooks) RemovePeer(id string) {
	pb.Lock()
	defer pb.Unlock()
	delete(pb.sets, id)
}

func (pb *PeerBooks) GetAll() []*Peer {
	pb.RLock()
	defer pb.RUnlock()

	m := make([]*Peer,0,len(pb.sets))
	for _, v := range pb.sets {
		m = append(m,v)
	}
	return m
}
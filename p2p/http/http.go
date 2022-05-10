package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/common"
	"io/ioutil"
	"net/http"
	"sync"
	"heyuanlong/blockchain-step/p2p"

)

type HTTPNetWork struct {
	LocalAddress string   // 本机地址
	NodeID       string   // 节点ID
	peerBooks    *p2p.PeerBooks

	msgQueue     chan *HTTPMsg
	recvCB       map[string]p2p.OnReceive
	sync.RWMutex
}

type HTTPMsg struct {
	*p2p.BroadcastMsg
	*p2p.Peer
}

func New(nodeAddrs []string, local string, nodeID string ) p2p.P2pI {
	obj := &HTTPNetWork{
		LocalAddress:local,
		NodeID: nodeID,
		peerBooks: p2p.NewPeerBooks(),

		msgQueue:     make(chan *HTTPMsg, 1000),
		recvCB: make(map[string]p2p.OnReceive),
	}

	for _, v := range nodeAddrs {
		obj.peerBooks.AddPeer(&p2p.Peer{
			ID :v,
			Address :v,
		})
	}

	return obj
}


func (ts *HTTPNetWork ) Start() error{
	router := gin.Default()

	router.GET("/broadcast", ts.commonHander)
	router.POST("/broadcast", ts.commonHander)


	go ts.Recv()

	router.Run(ts.LocalAddress)
	return nil
}


func (ts *HTTPNetWork ) RegisterOnReceive(MsgType string, callBack p2p.OnReceive) error{
	ts.Lock()
	ts.recvCB[MsgType] = callBack
	ts.Unlock()

	return nil
}


func (ts *HTTPNetWork ) Broadcast( msg *p2p.BroadcastMsg) error{
	requestBody, _ := json.Marshal(msg)
	header := map[string]string{
		"peer_id":ts.NodeID,
		"peer_address":"http://"+ts.LocalAddress,
	}

	for _, peer := range ts.peerBooks.GetAll() {
		go func(addr string) {
			_, err :=common.HttpDo(addr,"POST", map[string]string{},  header ,requestBody,5, map[string]interface{}{})
			if err != nil {
				log.Errorf("P2P 广播出错, err: %v", err)
			}
		}(peer.Address+"/broadcast")
	}

	return nil
}

// 广播到指定的peer
func (ts *HTTPNetWork ) BroadcastToPeer( msg *p2p.BroadcastMsg, p *p2p.Peer) error{
	requestBody, _ := json.Marshal(msg)
	header := map[string]string{
		"peer_id":ts.NodeID,
		"peer_address":"http://"+ts.LocalAddress,
	}

	_, err :=common.HttpDo(p.Address+"/broadcast","POST", map[string]string{},  header ,requestBody,5, map[string]interface{}{})
	if err != nil {
		log.Errorf("P2P 广播出错, err: %v", err)
	}

	return nil
}

// 广播 除了指定的peer
func (ts *HTTPNetWork ) BroadcastExceptPeer( msg *p2p.BroadcastMsg, p *p2p.Peer) error{
	requestBody, _ := json.Marshal(msg)
	header := map[string]string{
		"peer_id":ts.NodeID,
		"peer_address":"http://"+ts.LocalAddress,
	}

	for _, peer := range ts.peerBooks.GetAll() {
		if peer.Address == p.Address {
			continue
		}
		go func(addr string) {
			_, err :=common.HttpDo(addr,"POST", map[string]string{},  header ,requestBody,5, map[string]interface{}{})
			if err != nil {
				log.Errorf("P2P 广播出错, err: %v", err)
			}
		}(peer.Address+"/broadcast")
	}

	return nil
}

// 移除某个peer
func (ts *HTTPNetWork ) RemovePeer(p *p2p.Peer) error{
	ts.peerBooks.RemovePeer(p.ID)

	return nil
}

// 返回所有存在的peers
func (ts *HTTPNetWork ) Peers() ([]*p2p.Peer, error){
	peers := make([]*p2p.Peer, 0)

	return peers,nil
}


//hander------------------------------------------------------------

func (ts *HTTPNetWork ) commonHander(c *gin.Context) {
	content, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Debugf("读取请求内容出错 %s", err.Error())
		return
	}

	var revMsg p2p.BroadcastMsg
	if err := json.Unmarshal(content, &revMsg); err != nil {
		log.Debugf("解码请求内容出错 %s", err.Error())
		return
	}
	peer := p2p.Peer{
		ID:      c.GetHeader("peer_id"),
		Address: c.GetHeader("peer_address"),
	}

	select {
	case ts.msgQueue <- &HTTPMsg{
		&revMsg,
		&peer,
	}:
	default:
	}

	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}

func (ts *HTTPNetWork) Recv() {
	log.Debugf("开始接收消息")
	for {
		select {
		case msg := <-ts.msgQueue:
			onReceive := ts.recvCB[msg.MsgType]
			if onReceive != nil {
				go onReceive(msg.MsgType, msg.Msg, msg.Peer)
			} else {
				log.Debugf("当前消息ID没有相对应的处理模块 msgID: %s", msg.MsgType)
			}
		}
	}
}
//-----------------------------------------------------------------
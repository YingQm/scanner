package node

import (
	"fmt"
	"gitlab.33.cn/wallet/monitor/sync"
	"gitlab.33.cn/wallet/monitor/types"
	"net"
	"strings"
	sync2 "sync"
	"time"
)

type CoinNodePort struct {
	CoinType    string
	Node        []string
	NotAvalible chan types.Message
}

type CoinNodeSync struct {
	CoinType    string
	Node        []string
	TargetAddr  string
	Type        int64
	MaxDiff     int64
	NotAvalible chan types.Message
}

func NewCoinNodePort(ctype string, addrs []string) *CoinNodePort {
	coin := &CoinNodePort{
		CoinType:    ctype,
		Node:        addrs,
		NotAvalible: make(chan types.Message, len(addrs)),
	}
	return coin
}

func (c *CoinNodePort) RunScan() {
	for _, v := range c.Node {
		_, err := net.Dial("tcp", v)
		if err != nil {
			c.NotAvalible <- types.Message{CoinType: c.CoinType, Addr: v, ErrMsg: err.Error()}
		}
	}
	close(c.NotAvalible)
}

func NewCoinNodeSync(node *types.NodeSync) *CoinNodeSync {
	addrArr := strings.Split(strings.TrimSpace(node.Nodes), ",")
	coin := &CoinNodeSync{
		CoinType:    node.CoinType,
		Node:        addrArr,
		TargetAddr:  node.Target,
		Type:        node.NodeType,
		MaxDiff:     node.MaxDiff,
		NotAvalible: make(chan types.Message, len(addrArr)),
	}
	return coin
}

func (c *CoinNodeSync) GetBlockHeight(cfg *types.Config) []map[string]interface{} {
	timeFormat := "20060102 15:04:05"
	heightArr := make([]map[string]interface{}, 0)

	wg := sync2.WaitGroup{}
	wg.Add(len(c.Node))

	for _, addr := range c.Node {
		go func(addr string) {
			fmt.Printf("%s %s GetBlockHeight start:%s\n", c.CoinType, addr, time.Now().Format(timeFormat))
			height := sync.GetBlockHeight(c.CoinType, addr, c.Type, cfg)
			heightArr = append(heightArr, height)
			wg.Done()
			fmt.Printf("%s %s GetBlockHeight end:%s\n", c.CoinType, addr, time.Now().Format(timeFormat))
		}(addr)

	}
	wg.Wait()
	return heightArr
}

func (c *CoinNodeSync) RunScan() {
	for _, addr := range c.Node {
		_, err := net.Dial("tcp", addr)
		if err != nil {
			c.NotAvalible <- types.Message{CoinType: c.CoinType, Addr: addr, ErrMsg: err.Error()}
		}
	}
	close(c.NotAvalible)
}

func (c *CoinNodeSync) RunSync(cfg *types.Config) {
	var msgArr []types.Message
	switch c.CoinType {
	case "ETH", "ETC":
		msgArr = sync.EthrumSync(c.Node, c.TargetAddr, c.MaxDiff)
	case "BTC":
		msgArr = sync.BtcSync(c.Node, c.TargetAddr, c.Type, c.MaxDiff, cfg)
	case "USDT":
		msgArr = sync.UsdtSync(c.Node, c.TargetAddr, cfg)

	default:

	}
	for _, msg := range msgArr {
		c.NotAvalible <- msg
	}
	close(c.NotAvalible)
}

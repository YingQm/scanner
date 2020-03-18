package util

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"gitlab.33.cn/wallet/monitor/node"
	"gitlab.33.cn/wallet/monitor/types"
	"strings"
	"time"
)

type Cache struct {
	cfg             *types.Config
	Xdb             *DbHandler
	NodePortLru     *lru.Cache
	NodeSyncLru     *lru.Cache
	NodeHeightLru   *lru.Cache
	NodeParallelLru *lru.Cache
}

func NewCache(cfg *types.Config, handler *DbHandler) *Cache {
	nodePortLru, err := lru.New(512)
	if err != nil {
		panic(err)
	}
	nodeSyncLru, err := lru.New(512)
	if err != nil {
		panic(err)
	}
	nodeHeightLru, err := lru.New(512)
	if err != nil {
		panic(err)
	}
	nodeParallelLru, err := lru.New(512)
	if err != nil {
		panic(err)
	}
	cache := &Cache{
		Xdb:             handler,
		cfg:             cfg,
		NodePortLru:     nodePortLru,
		NodeSyncLru:     nodeSyncLru,
		NodeHeightLru:   nodeHeightLru,
		NodeParallelLru: nodeParallelLru,
	}
	cache.initParalleLru()
	return cache
}

func (c *Cache) initParalleLru() {
	c.NodeParallelLru.Add("DTC", "daitaochain.coins")
	c.NodeParallelLru.Add("BECC", "beechain.coins")
}

func (c *Cache) UpdateCache() {
	c.UpdateNodePortInfo()
	c.UpdateNodeSyncInfo()
	c.UpdateNodeHeight()
	updateAppTicker := time.NewTicker(10 * time.Minute)
	updateHeightTicker := time.NewTicker(5 * time.Minute)
	defer updateAppTicker.Stop()
	for {
		select {
		case <-updateAppTicker.C:
			c.UpdateNodePortInfo()
			c.UpdateNodeSyncInfo()

		case <-updateHeightTicker.C:
			c.UpdateNodeHeight()
		}
	}
}

func (c *Cache) UpdateNodePortInfo() {
	log.Info("Cache UpdateNodePortInfo")
	nodeinfos := c.Xdb.FetchNodePortAddr()

	for key, val := range nodeinfos {
		urls := strings.TrimSpace(val)
		log.Info("Cache NodePortInfo", "key", key, "value", urls)
		c.NodePortLru.Add(key, urls)
	}
}

func (c *Cache) UpdateNodeSyncInfo() {
	log.Info("Cache UpdateNodeSyncInfo")
	nodeSyncArr := c.Xdb.FetchNodeSyncAddr()
	for _, v := range nodeSyncArr {
		log.Info("Cache NodeSyncInfo", "cointype", v.CoinType, "nodes", v.Nodes, "nodetype", v.NodeType)
		key := fmt.Sprintf("%s-%d", v.CoinType, v.NodeType)
		c.NodeSyncLru.Add(key, v)
	}
}

//缓存USDT的节点高度
func (c *Cache) UpdateNodeHeight() {

	coinTypes := c.NodeSyncLru.Keys()
	for _, coin := range coinTypes {
		coinType := coin.(string)
		find := strings.Contains(coinType, "USDT")
		if !find {
			continue
		}
		coinNodeInfo, ok := c.NodeSyncLru.Get(coinType)
		if !ok {
			continue
		}
		log.Info("Cache UpdateCoinNodeHeight", "coinType", coinType, "nodeinfo", coinNodeInfo)
		nodeinfo := coinNodeInfo.(types.NodeSync)
		coin := node.NewCoinNodeSync(&nodeinfo)
		heightMapArr := coin.GetBlockHeight(c.cfg)
		c.NodeHeightLru.Add(coinType, heightMapArr)
	}
	//遍历打印缓存的信息
	log.Info("Cache PrintCacheNodeHeight")
	cointypes := c.NodeHeightLru.Keys()
	for _, ctype := range cointypes {
		heightInfo, ok := c.NodeHeightLru.Get(ctype)
		if !ok {
			continue
		}
		log.Info("Cache NodeHeight", "cointype", ctype, "heightinfo", heightInfo)
	}
}

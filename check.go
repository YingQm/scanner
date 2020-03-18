package main

import (
	"fmt"
	"gitlab.33.cn/wallet/monitor/node"
	"gitlab.33.cn/wallet/monitor/send"
	"gitlab.33.cn/wallet/monitor/types"
	"strings"
	sync2 "sync"
	"time"
)

func RunCheck() {
	for {
		log.Info("start check", " check interval in minute(s)", cfg.IntervalTime)

		go CheckNodePort()
		go CheckNodeSync(cfg)
		time.Sleep((time.Duration)(cfg.IntervalTime) * time.Minute)
	}
}

func CheckNodePort() {
	var coinsArr = []*node.CoinNodePort{}
	wg := sync2.WaitGroup{}
	var checkResult string

	coinTypes := cache.NodePortLru.Keys()
	log.Info("CheckNodePort", "NodePortLru keys", coinTypes)
	for _, v := range coinTypes {
		cointype := v.(string)
		nodeAddrs, ok := cache.NodePortLru.Get(cointype)
		if !ok {
			continue
		}
		nodeAddrsStr := nodeAddrs.(string)
		addrs := strings.Split(nodeAddrsStr, ",")
		if len(addrs) == 0 {
			continue
		}
		coin := node.NewCoinNodePort(cointype, addrs)
		coinsArr = append(coinsArr, coin)
		wg.Add(1)
		go func() {
			coin.RunScan()
			wg.Done()
		}()
	}

	wg.Wait()

	for _, coin := range coinsArr {
		tmpResult := ""
		for message := range coin.NotAvalible {
			tmpResult += message.Addr + "<br><hr />\n"
		}
		if len(tmpResult) > 0 {
			tmpResult = fmt.Sprintf("%s 端口不通的节点列表:<br>\n %s\n", coin.CoinType, tmpResult)
			checkResult += tmpResult + "<hr />\n"
		}
	}
	if len(checkResult) > 0 {
		send.SendEmail(cfg, "提醒:服务器端口不通", checkResult)
	}
}

func CheckNodeSync(cfg *types.Config) {
	var coinsArr = []*node.CoinNodeSync{}
	var syncResult string
	wg := sync2.WaitGroup{}
	keys := cache.NodeSyncLru.Keys()
	log.Info("CheckNodeSync", "NodeSyncLru keys", keys)
	for _, key := range keys {
		coinNode, ok := cache.NodeSyncLru.Get(key.(string))
		if !ok {
			continue
		}
		nodeinfo := coinNode.(types.NodeSync)
		nodeAddrs := strings.Split(nodeinfo.Nodes, ",")
		log.Info("CheckNodeSync", "cointype", nodeinfo.CoinType, "nodeaddrs:", nodeAddrs, "targetAddr", nodeinfo.Target)
		coin := node.NewCoinNodeSync(&nodeinfo)
		coinsArr = append(coinsArr, coin)
		wg.Add(1)
		go func() {
			coin.RunSync(cfg)
			wg.Done()
		}()
	}

	wg.Wait()
	for _, coin := range coinsArr {
		tmpResult := ""
		for message := range coin.NotAvalible {
			tmpResult += message.Addr + " " + message.ErrMsg + "<br><hr />\n"
		}
		if len(tmpResult) > 0 {
			tmpResult = fmt.Sprintf("%s 同步异常的节点列表:<br>\n %s\n", coin.CoinType, tmpResult)
			syncResult += tmpResult + "<hr />\n"
		} else {
			result := fmt.Sprintf("check %s nothing abnormal detected", coin.CoinType)
			log.Info(result)
		}
	}
	if len(syncResult) > 0 {
		send.SendEmail(cfg, "提醒:服务器节点不同步", syncResult)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.33.cn/wallet/monitor/node"
	"gitlab.33.cn/wallet/monitor/sync"
	"gitlab.33.cn/wallet/monitor/types"
	"net"
	"net/http"
	"strings"
	"time"
)

// request url "http:localhost:9988/getblockheight?cointype=btc"
func GetBlockHeightHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	reqIpPort := r.RemoteAddr
	ip, _, err := net.SplitHostPort(reqIpPort)
	if err != nil {
		log.Error("http", "Error while parse request ip:", err)
	}
	log.Info("GetBlockHeight", " request ip", ip)

	if rateLimiter != nil {
		if jobChan != nil {
			select {
			case jobChan <- struct{}{}:
				if rateLimiter.Allow(ip) {
					handleGetBlockHeight(w, r)
				} else {
					limitNum := cfg.RateLimit.MaxCount
					RespError(w, fmt.Sprintf("rate limit exceeded,retry later(maximum %d times per minute every ip)", limitNum))
				}
				<-jobChan
			default:
				log.Info("http", "jobChan", "chan full")
				RespError(w, fmt.Sprintf("rate limit exceeded,retry later(maximum 10 times per minute)"))
			}
			return
		}
		if rateLimiter.Allow(ip) {
			handleGetBlockHeight(w, r)
			return
		}
		RespError(w, "rate limit exceeded,retry later(maximum 10 times per minute)")
		return
	}

	handleGetBlockHeight(w, r)
}

func handleGetBlockHeight(w http.ResponseWriter, r *http.Request) {
	cointype := r.Form.Get("cointype")
	if cointype == "" {
		RespError(w, "bad request cointype missing")
		return
	}

	var err error
	coinType := strings.ToUpper(cointype)
	heightMapArr := make([]map[string]interface{}, 0)

	//usdt的高度从缓存中读取，其他的币种实时查询
	if coinType == "USDT" {
		result, ok := cache.NodeHeightLru.Get("USDT-0")
		if !ok {
			log.Info("USDT cache not found ")
			heightMapArr, err = GetNodeBlockHeight(coinType, cfg)
			if err == nil {
				cache.NodeHeightLru.Add("USDT-0", heightMapArr)
			}
		} else {
			log.Info("get USDT height from cache ")
			heightMapArr = result.([]map[string]interface{})
		}
	} else {
		heightMapArr, err = GetNodeBlockHeight(coinType, cfg)
		if err != nil {
			log.Error("GetBlockHeight", " error", err.Error())
		}
	}
	resp := types.ClientResponse{Id: 1}
	if err != nil {
		log.Info("GetNodeBlockHeight", "query", coinType, "error", err.Error())
		resp.Error = err.Error()
	}

	resp.Result = heightMapArr

	jsonBytes, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Error("GetBlockHeight", "MarshalIndent err", err.Error())
	}
	w.Write(jsonBytes)
}

func GetNodeBlockHeight(searchCoinType string, cfg *types.Config) ([]map[string]interface{}, error) {
	var (
		err error
	)

	log.Info("GetNodeBlockHeight", "query", searchCoinType, "start time", time.Now().Format("20060102 15:04:05"))
	heightMapArr := make([]map[string]interface{}, 0)
	coinTypes := cache.NodeSyncLru.Keys()
	var findSearchCoin bool
	for _, coin := range coinTypes {
		coinType := coin.(string)
		if strings.Contains(coinType, searchCoinType) {
			findSearchCoin = true
			coinNodeUrls, ok := cache.NodeSyncLru.Get(coinType)
			if !ok {
				continue
			}
			log.Info("GetNodeBlockHeight", "coinType", coinType, "nodeinfo", coinNodeUrls)
			nodeinfo := coinNodeUrls.(types.NodeSync)
			coin := node.NewCoinNodeSync(&nodeinfo)
			heightMap := coin.GetBlockHeight(cfg)
			heightMapArr = append(heightMapArr, heightMap...)
		}
	}

	coinParallelTypes := cache.NodeParallelLru.Keys()
	for _, pcoin := range coinParallelTypes {
		pcoinType := pcoin.(string)
		if strings.Contains(pcoinType, searchCoinType) {
			findSearchCoin = true
			pcoinSymbol, ok := cache.NodeParallelLru.Get(pcoinType)
			if !ok {
				continue
			}
			log.Info("GetParallelNodeBlockHeight", "coinType", pcoinType, "Symbol", pcoinSymbol)
			height := sync.GetBlockHeight(searchCoinType, pcoinSymbol.(string), 0, cfg)
			heightArr := make([]map[string]interface{}, 0)
			heightArr = append(heightArr, height)
			heightMapArr = append(heightMapArr, heightArr...)
		}
	}

	if !findSearchCoin {
		err = errors.New(fmt.Sprintf("%s not surpport", searchCoinType))
		log.Info("GetNodeBlockHeight", "query", searchCoinType, "error", err.Error())
	}
	log.Info("GetNodeBlockHeight", "query", searchCoinType,
		"end time", time.Now().Format("20060102 15:04:05"), "result", heightMapArr)
	return heightMapArr, err
}

func RespError(w http.ResponseWriter, errMessage string) {
	resp := types.ClientResponse{Id: 1}
	resp.Result = make([]map[string]interface{}, 0)
	resp.Error = errMessage
	result, _ := json.Marshal(resp)
	w.Write(result)
}

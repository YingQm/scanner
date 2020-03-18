package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"log"
	"strconv"
	"strings"
)

func BtcSync(serviceAddr []string, targetUrl string, nodeType int64, topDiff int64, cfg *types.Config) []types.Message {
	var msgArr []types.Message
	//如果没有设置target节点返回空
	if targetUrl == "" {
		fmt.Printf("targetAddr empty")
		return msgArr
	}
	heightTarget, err := GetTargetBlockHeight(targetUrl)
	if err != nil {
		fmt.Printf("BTC GetBlockHeight \"%s\" err %s\n", targetUrl, err.Error())
		msg := types.Message{CoinType: "BTC", Addr: targetUrl, ErrMsg: err.Error()}
		msgArr = append(msgArr, msg)
		return msgArr
	}

	for i := 0; i < len(serviceAddr); i++ {
		var (
			height uint64
			err    error
		)
		//去除地址中的空格
		addr := strings.TrimSpace(serviceAddr[i])
		if nodeType == 1 {
			height, err = GetInsightBlockHeight(addr)
		} else {
			height, err = GetBtcBlockHeight(addr, cfg)
		}
		if err != nil {
			errMsg := fmt.Sprintf("GetBlockHeight error %s \n", err.Error())
			msg := types.Message{CoinType: "BTC", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		} else if heightTarget > height && heightTarget-height > uint64(topDiff) {
			errMsg := fmt.Sprintf("%s 落后 %d 个高度, 目标节点高度=%d, 当前节点高度=%d \n", addr, heightTarget-height, heightTarget, height)
			msg := types.Message{CoinType: "BTC", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		}

	}
	return msgArr
}

func GetInsightBlockHeight(url string) (uint64, error) {
	getstr := url + "/insight-api" + "/status?q=getinfo"
	resp, err := SendToServerTls("GET", getstr, nil)
	if err != nil {
		return 0, err
	}

	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("info").Interface() != nil {
		height := js.Get("info").Get("blocks").MustUint64()
		fmt.Printf("BTC %s GetInsightBlockHeight response:%d\n", url, height)
		return height, nil
	}

	return 0, errors.New(js.Get("error").MustString())
}

func GetBtcBlockHeight(url string, cfg *types.Config) (uint64, error) {
	postdata := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0","method":"getblockcount", "params":[]}`)
	resp, err := SendToServer_v1("POST", url, strings.NewReader(postdata), "", "", cfg)
	if err != nil {
		return 0, err
	}
	fmt.Printf("BTC %s GetBtcBlockHeight resp:%s\n", url, string(resp))
	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("result").Interface() != nil {
		height := js.Get("result").MustInt64()

		return uint64(height), nil
	}
	return 0, errors.New("GetBlockHeight Err")
}

//blockchain.info 限速10秒1次
func GetTargetBlockHeight(url string) (uint64, error) {
	// "https://blockchain.info/q/getblockcount"
	requestURL := url + "/q/getblockcount"
	result, err := SendToServerTls("GET", requestURL, nil)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		return 0, err
	}
	fmt.Println("blockchain.info response height:", string(result))
	retStr := string(result)
	height, err := strconv.ParseUint(retStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}

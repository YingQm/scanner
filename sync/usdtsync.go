package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"strings"
)

func UsdtSync(serviceAddr []string, mainAddr string, cfg *types.Config) []types.Message {
	var msgArr []types.Message
	if mainAddr == "" {
		fmt.Printf("targetAddr empty")
		return msgArr
	}
	heightMain, err := GetTargetBlockHeight(mainAddr)
	if err != nil {
		fmt.Printf("GetBlockHeight \"%s\" err %s\n", mainAddr, err.Error())
		msg := types.Message{CoinType: "USDT", Addr: mainAddr, ErrMsg: err.Error()}
		msgArr = append(msgArr, msg)
		return msgArr
	}

	for i := 0; i < len(serviceAddr); i++ {
		addr := strings.TrimSpace(serviceAddr[i])
		height, err := getUsdtBlockHeight(addr, cfg)
		if err != nil {
			errMsg := fmt.Sprintf("GetBlockHeight error %s \n", err.Error())
			msg := types.Message{CoinType: "USDT", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		} else if heightMain > height && heightMain-height > 12 {
			errMsg := fmt.Sprintf("%s 落后 %d 个高度, 目标节点高度=%d, 当前节点高度=%d \n", addr, heightMain-height, heightMain, height)
			msg := types.Message{CoinType: "USDT", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		}
	}
	return msgArr
}

func getUsdtBlockHeight(url string, cfg *types.Config) (uint64, error) {
	postdata := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0","method":"omni_getinfo", "params":[]}`)
	resp, err := SendToServer_v1("POST", url, strings.NewReader(postdata), cfg.Omniname, cfg.Omnipasswd, cfg)
	if err != nil {
		return 0, err
	}

	fmt.Printf("USDT %s getUsdtBlockHeight response:%s\n", url, string(resp))

	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("result").Interface() != nil {
		height := js.Get("result").Get("block").MustInt64()
		return uint64(height), nil
	}
	if js.Get("error").Interface() != nil {
		errMsg := js.Get("error").Get("message").MustString()
		return 0, errors.New(errMsg)
	}

	return 0, errors.New("GetBlockHeight Err")
}

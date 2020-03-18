package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"math/big"
	"strings"
)

// "https://mainnet.infura.io" http://39.107.57.133:8545","http://39.107.60.254:8545","http://47.106.117.142:8801
func EthSync(ServiceAddr []string) string {
	heightMain, err := getEthBlockHeight("https://mainnet.infura.io")
	if err != nil {
		return fmt.Sprintf("GetBlockHeight \"https://mainnet.infura.io\" err %s", err.Error())
	}

	result := ""
	for i := 0; i < len(ServiceAddr); i++ {
		height, err := getEthBlockHeight(ServiceAddr[i])
		if err != nil {
			result += fmt.Sprintf("%s GetBlockHeight err %s \n", ServiceAddr[i], err.Error())
		} else if heightMain > height && heightMain-height > 12 {
			result += fmt.Sprintf("%s %d height behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
		}
	}

	return result
}

func getEthBlockHeight(url string) (uint64, error) {
	postdata := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0","method":"eth_blockNumber", "params":[]}`)
	resp, err := SendToServerTls("POST", url, strings.NewReader(postdata))
	if err != nil {
		return 0, err
	}

	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("result").Interface() != nil {
		hexstr := js.Get("result").MustString()
		bn := big.NewInt(0)
		var ok bool
		bn, ok = bn.SetString(hexstr[2:], 16)
		if !ok {
			return 0, errors.New("hextodec err")
		}
		fmt.Printf("ETH %s getEthBlockHeight respone:%d\n", url, bn.Uint64())
		return bn.Uint64(), nil
	}
	return 0, errors.New("GetBlockHeight Err")
}

func EthrumSync(serviceAddr []string, mainAddr string, maxdiff int64) []types.Message {
	var msgArr []types.Message
	if mainAddr == "" {
		fmt.Printf("targetAddr empty")
		return msgArr
	}
	heightMain, err := getEthBlockHeight(mainAddr)
	if err != nil {
		fmt.Printf("GetBlockHeight \"%s\" err %s\n", mainAddr, err.Error())
		msg := types.Message{CoinType: "ETH", Addr: mainAddr, ErrMsg: err.Error()}
		msgArr = append(msgArr, msg)
		return msgArr
	}

	for i := 0; i < len(serviceAddr); i++ {
		addr := strings.Trim(serviceAddr[i], " ")
		height, err := getEthBlockHeight(addr)
		if err != nil {
			errMsg := fmt.Sprintf("GetBlockHeight error %s \n", err.Error())
			msg := types.Message{CoinType: "ETH", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		} else if heightMain > height && heightMain-height > uint64(maxdiff) {
			errMsg := fmt.Sprintf("%s 落后 %d 个高度, 目标节点高度=%d, 当前节点高度=%d \n", addr, heightMain-height, heightMain, height)
			msg := types.Message{CoinType: "ETH", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		}
	}
	return msgArr
}

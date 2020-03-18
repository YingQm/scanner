package sync

import (
<<<<<<< HEAD
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"math/big"
=======
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
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
<<<<<<< HEAD
			result += fmt.Sprintf("%s %d height behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
=======
			result += fmt.Sprintf("%s %d levels behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
		}
	}

	return result
}

func getEthBlockHeight(url string) (uint64, error) {
	postdata := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0","method":"eth_blockNumber", "params":[]}`)
<<<<<<< HEAD
	resp, err := SendToServerTls("POST", url, strings.NewReader(postdata))
=======
	resp, err := sendToServerTls("POST", url, strings.NewReader(postdata))
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
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
<<<<<<< HEAD
		fmt.Printf("ETH %s getEthBlockHeight respone:%d\n", url, bn.Uint64())
=======
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
		return bn.Uint64(), nil
	}
	return 0, errors.New("GetBlockHeight Err")
}

<<<<<<< HEAD
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
=======
func sendToServerTls(method, url string, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Close = true
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	rbs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return rbs, nil
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
}

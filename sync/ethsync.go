package sync

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
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
			result += fmt.Sprintf("%s %d levels behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
		}
	}

	return result
}

func getEthBlockHeight(url string) (uint64, error) {
	postdata := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0","method":"eth_blockNumber", "params":[]}`)
	resp, err := sendToServerTls("POST", url, strings.NewReader(postdata))
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
		return bn.Uint64(), nil
	}
	return 0, errors.New("GetBlockHeight Err")
}

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
}

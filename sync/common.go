package sync

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gitlab.33.cn/wallet/monitor/types"
	"io"
	"io/ioutil"
	"net/http"
)

func GetBlockHeight(coinType string, url string, nodeType int64, cfg *types.Config) map[string]interface{} {
	var height uint64
	var err error
	fmt.Printf("GetBlockHeight cointype:%s, url:%s\n", coinType, url)
	result := make(map[string]interface{}, 0)
	switch coinType {
	case "BTC":
		if nodeType == 1 {
			height, err = GetInsightBlockHeight(url)
		} else {
			height, err = GetBtcBlockHeight(url, cfg)
		}
	case "USDT":
		height, err = getUsdtBlockHeight(url, cfg)
	case "ETH", "ETC":
		height, err = getEthBlockHeight(url)
	case "BTY":
		height, err = GetBtyBlockHeight(url)
	case "EOS":
		height, err = GetEosBlockHeight(url)
	case "BNB":
		height, err = GetBnbBlockHeight(url)
	case "DCR":
		if nodeType == 1 {
			height, err = GetDrcInsightBlockHeight(url)
		} else {
			height, err = GetDrcBlockHeight(url, cfg)
		}
	case "DTC":
		height, err = GetBtyParallelHeight(url)

	default:
		err = errors.New(fmt.Sprintf("%s not surpport", coinType))
		fmt.Println("default", err.Error())
	}
	result["cointype"] = coinType
	result["url"] = url
	result["height"] = height
	if err != nil {
		result["error"] = err.Error()
	} else {
		result["error"] = ""
	}
	return result
}

func SendToServerTls(method, url string, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		//Timeout:   10 * time.Second,
	}

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

func SendToServer_v1(method, url string, body io.Reader, user, pwd string, cfg *types.Config) ([]byte, error) {
	client := http.DefaultClient
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(cfg.Rpcname, cfg.Rpcpasswd)

	if user != "" || pwd != "" {
		request.SetBasicAuth(user, pwd)

	}
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

func SendToServerTls_v1(method, url string, body io.Reader, user, pwd string, cfg *types.Config) ([]byte, error) {

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
	request.SetBasicAuth(cfg.Rpcname, cfg.Rpcpasswd)

	if user != "" || pwd != "" {
		request.SetBasicAuth(user, pwd)

	}
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

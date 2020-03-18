package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"strings"
)

// serviceAddr=["https://47.106.117.142:8804"]  explorAddr= ["http://47.106.117.142:3003"]
func DcrSync(ServiceAddr []string, cfg *types.Config) string {
	heightMain, err := GetDrcInsightBlockHeight("http://47.106.117.142:3003")
	if err != nil {
		return fmt.Sprintf("GetBlockHeight \"https://mainnet.infura.io\" err %s", err.Error())
	}

	result := ""
	for i := 0; i < len(ServiceAddr); i++ {
		height, err := GetDrcBlockHeight(ServiceAddr[i], cfg)
		if err != nil {
			result += fmt.Sprintf("%s GetBlockHeight err %s \n", ServiceAddr[i], err.Error())
		} else if heightMain > height && heightMain-height > 12 {
			result += fmt.Sprintf("%s %d levels behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
		}
	}

	return result
}

func GetDrcInsightBlockHeight(url string) (uint64, error) {
	getstr := fmt.Sprintf(`/api/status?q=getinfo`)

	var resp []byte
	var senderr error
	var err error
	resp, err = SendToServerTls("GET", url+getstr, nil)
	if err != nil {
		senderr = err
	}

	if resp == nil {
		return 0, senderr
	}

	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("info").Interface() != nil {
		return js.Get("info").Get("blocks").MustUint64(), nil
	}
	return 0, errors.New(js.Get("error").MustString())
}

func GetDrcBlockHeight(url string, cfg *types.Config) (uint64, error) {
	postdata := fmt.Sprintf(`{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`)
	fmt.Printf("name:%s, passwd:%s, poststr:%s\n", cfg.Rpcname, cfg.Rpcpasswd, postdata)

	resp, err := SendToServerTls_v1("POST", url, strings.NewReader(postdata), cfg.Rpcname, cfg.Rpcpasswd, cfg)
	if err != nil {
		return 0, err
	}
	fmt.Printf("%s GetDrcBlockHeight resp:%s\n", url, string(resp))
	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}
	if js.Get("result").Interface() != nil {
		height := js.Get("result").Get("blocks").MustUint64()
		fmt.Printf("DCR GetBlockHeight return height:%d\n", height)
		return height, nil
	}
	if js.Get("error").Interface() != nil {
		errMsg := js.Get("error").Get("message").MustString()
		return 0, errors.New(errMsg)
	}

	return 0, errors.New("DCR GetBlockHeight err")
}

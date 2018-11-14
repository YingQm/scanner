package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

// serviceAddr=["https://47.106.117.142:8804"]  explorAddr= ["http://47.106.117.142:3003"]
func DcrSync(ServiceAddr []string) string {
	heightMain, err := getDrcBlockHeight("http://47.106.117.142:3003")
	if err != nil {
		return fmt.Sprintf("GetBlockHeight \"https://mainnet.infura.io\" err %s", err.Error())
	}

	result := ""
	for i := 0; i < len(ServiceAddr); i++ {
		height, err := getDrcBlockHeight(ServiceAddr[i])
		if err != nil {
			result += fmt.Sprintf("%s GetBlockHeight err %s \n", ServiceAddr[i], err.Error())
		} else if heightMain > height && heightMain-height > 12 {
			result += fmt.Sprintf("%s %d levels behind, heightMain=%d height=%d \n", ServiceAddr[i], heightMain-height, heightMain, height)
		}
	}

	return result
}

func getDrcBlockHeight(url string) (uint64, error) {
	getstr := fmt.Sprintf(`/api/status?q=getinfo`)

	var resp []byte
	var senderr error
	var err error
	resp, err = sendToServerTls("GET", url+getstr, nil)
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

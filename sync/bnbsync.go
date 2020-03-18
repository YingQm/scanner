package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"strconv"
)

func GetBnbBlockHeight(url string) (uint64, error) {
	method := "/status"
	Url := url + method

	fmt.Println("request url:", Url)
	resp, err := SendToServerTls("GET", Url, nil)
	if err != nil {
		return 0, err
	}

	js, err := simplejson.NewJson(resp)
	if err != nil {
		return 0, err
	}

	if js.Get("result").Get("sync_info").Interface() != nil {
		heightString := js.Get("result").Get("sync_info").Get("latest_block_height").MustString()
		height, err := strconv.ParseUint(heightString, 0, 64)
		if err != nil {
			fmt.Println("BNB GetBlockHeight", "ParseUintError:", err.Error())
			return 0, err
		}
		fmt.Println("BNB GetBlockHeight return height:", height)
		return height, nil
	}

	return 0, errors.New("BNB GetBlockHeight Err")

}

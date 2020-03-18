package sync

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

func GetEosBlockHeight(url string) (uint64, error) {
	posturl := url + "/v1/chain/get_info"

	resp, err := SendToServerTls("POST", posturl, nil)
	if err != nil {
		return 0, err
	}

	js, err := simplejson.NewJson([]byte(resp))
	if err != nil {
		return 0, err
	}

	if js.Get("head_block_num").Interface() != nil {
		height := js.Get("head_block_num").MustUint64()
		fmt.Printf("EOS %s getEosBlockHeight response:%d\n", url, height)
		return height, nil
	}

	return 0, errors.New("GetBlockHeight Err")

}

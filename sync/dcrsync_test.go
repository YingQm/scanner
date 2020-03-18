package sync

import (
	"fmt"
	"gitlab.33.cn/wallet/monitor/types"
	"testing"
)

func TestGetDcrBlockHeight(t *testing.T) {
	url := "https://47.106.117.142:8804"

	cfg := new(types.Config)
	cfg.Rpcname = "fuzamei"
	cfg.Rpcpasswd = "fuzamei"
	height, err := GetDrcBlockHeight(url, cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("resp height:", height)
	t.Log("GetEosBlockHeight return height:", height)
}

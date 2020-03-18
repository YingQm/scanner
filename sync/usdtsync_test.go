package sync

import (
	"fmt"
	"gitlab.33.cn/wallet/monitor/types"
	"testing"
)

func TestGetUsdtBlockHeight(t *testing.T) {
	url := "http://183.129.226.76:8332"

	cfg := new(types.Config)
	cfg.Rpcname = "fuzamei"
	cfg.Rpcpasswd = "fuzamei"
	height, err := getUsdtBlockHeight(url, cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("resp height:", height)
	t.Log("GetUsdtBlockHeight return height:", height)
}

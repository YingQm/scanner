package sync

import (
	"flag"
	"gitlab.33.cn/wallet/monitor/config"
	"gitlab.33.cn/wallet/monitor/types"
	"testing"
)

var configpath = flag.String("f", "../config.toml", "configfile")
var cfg *types.Config

var url string

func init() {
	url = "http://47.106.117.142:8802"
	cfg = config.InitCfg(*configpath)
}

func TestGetTargetBlockHeight(t *testing.T) {
	url := "https://blockchain.info"
	height, err := GetTargetBlockHeight(url)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("target url return height:", height)
}

func TestGetInsightBlockHeight(t *testing.T) {
	url := "http://119.3.17.220:32815"
	height, err := GetInsightBlockHeight(url)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("insight-api return height:", height)
}

func TestGetBtcBlockHeight(t *testing.T) {
	height, err := GetBtcBlockHeight(url, cfg)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("insight-api return height:", height)

}

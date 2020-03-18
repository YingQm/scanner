package test

import (
	"encoding/json"
	"fmt"
	"gitlab.33.cn/wallet/monitor/sync"
	"gitlab.33.cn/wallet/monitor/types"
	"time"

	sync2 "sync"
	"testing"
	//"time"
)

var host string
var cfg *types.Config

func init() {
	host = "http://localhost:9988"
	cfg = new(types.Config)
	cfg.Rpcname = "fuzamei"
	cfg.Rpcpasswd = "fuzamei"
}

func TestRateLimit(t *testing.T) {
	nSec := 5
	wg := sync2.WaitGroup{}
	for i := 0; i < nSec; i++ {
		wg.Add(7)
		go func() {
			TestGetBTCBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetETHBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetUSDTBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetETCBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetEOSBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetDCRBlockHeight(t)
			wg.Done()
		}()
		go func() {
			TestGetBNBBlockHeight(t)
			wg.Done()
		}()
		time.Sleep(1 * time.Second)
	}
	wg.Wait()

}

func TestGetBTCBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "btc")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}
	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetBTCBlockHeight", result)
}

func TestGetETHBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "eth")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}
	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetETHBlockHeight", result)
}
func TestGetUSDTBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "usdt")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetUSDTBlockHeight", result)
}

func TestGetETCBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "etc")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetETCBlockHeight", result)
}

func TestGetEOSBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "eos")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetEOSBlockHeight", result)
}

func TestGetBTYBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "bty")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetBTYBlockHeight", result)
}

func TestGetDCRBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "dcr")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetDCRBlockHeight", result)
}

func TestGetBNBBlockHeight(t *testing.T) {
	requestURL := fmt.Sprintf("%s/getblockheight?cointype=%s", host, "bnb")
	resp, err := sync.SendToServerTls("GET", requestURL, nil)
	if err != nil {
		t.Logf("error: %s\n", err.Error())
		return
	}

	result := make(map[string]interface{}, 0)
	json.Unmarshal(resp, &result)

	t.Log("TestGetBNBBlockHeight", result)
}

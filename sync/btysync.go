package sync

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"gitlab.33.cn/wallet/monitor/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WalletRawReq struct {
	Cointype    string      `json:"cointype"`
	TokenSymbol string      `json:"symbol"`
	RawData     interface{} `json:"rawdata"`
	Tag         string      `json:"tag"`
}

type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	Id     uint64         `json:"id"`
}

type clientResponse struct {
	Id     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

func BtySync(serviceAddr []string, targetUrl string, maxdiff int64) []types.Message {
	var msgArr []types.Message
	if targetUrl == "" {
		return msgArr
	}
	heightTarget, err := GetBtyBlockHeight(targetUrl)
	if err != nil {
		fmt.Printf("GetBlockHeight \"%s\" err %s\n", targetUrl, err.Error())
		msg := types.Message{CoinType: "BTY", Addr: targetUrl, ErrMsg: err.Error()}
		msgArr = append(msgArr, msg)
		return msgArr
	}

	for i := 0; i < len(serviceAddr); i++ {
		addr := strings.TrimSpace(serviceAddr[i])
		height, err := GetBtyBlockHeight(addr)
		if err != nil {
			errMsg := fmt.Sprintf("GetBlockHeight error %s \n", err.Error())
			msg := types.Message{CoinType: "BTY", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		} else if heightTarget > height && heightTarget-height > uint64(maxdiff) {
			errMsg := fmt.Sprintf("%s 落后 %d 个高度, 目标节点高度=%d, 当前节点高度=%d \n", addr, heightTarget-height, heightTarget, height)
			msg := types.Message{CoinType: "BTY", Addr: addr, ErrMsg: errMsg}
			msgArr = append(msgArr, msg)
		}
	}
	return msgArr
}

func GetBtyBlockHeight(url string) (uint64, error) {
	poststr := fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"Chain33.GetLastHeader","params":[]}`)

	resp, err := SendToServerTls("POST", url, strings.NewReader(poststr))
	if err != nil {
		return 0, err
	}

	js, err := simplejson.NewJson([]byte(resp))
	if err != nil {
		return 0, err
	}

	if js.Get("result").Interface() != nil {
		height := js.Get("result").Get("height").MustUint64()
		fmt.Printf("BTY %s getBtyBlockHeight response:%d\n", url, height)
		return height, nil
	}

	return 0, errors.New(js.Get("error").MustString())
}

func GetBtyParallelHeight(TokenSymbol string) (uint64, error) {
	url := "https://114.55.101.159:8084"
	var reply interface{}
	rawdata := `{"jsonrpc":"2.0","id":2,"method":"Chain33.GetLastHeader","params":[]}`
	args := WalletRawReq{Cointype: "BTY", TokenSymbol: TokenSymbol, RawData: hex.EncodeToString([]byte(rawdata))}
	err := CallTls(url, "Wallet.Transport", &args, &reply)
	if err != nil {
		log.Fatal("Transport error:", err)
		return 0, err
	}

	jbs, err := json.Marshal(reply)
	if err != nil {
		log.Fatal("Marshal error:", err)
		return 0, err
	}

	js, err := simplejson.NewJson(jbs)
	if err != nil {
		return 0, err
	}

	if js.Get("height").Interface() != nil {
		height := js.Get("height").MustUint64()
		return height, nil
	}

	return 0, errors.New(js.Get("error").MustString())
}

func CallTls(url, method string, params, resp interface{}) error {
	req := &clientRequest{}
	req.Method = method
	req.Params[0] = params
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("post", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer request.Body.Close()

	httpcliTls := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	postresp, err := httpcliTls.Do(request)
	if err != nil {
		return err
	}
	defer postresp.Body.Close()
	b, err := ioutil.ReadAll(postresp.Body)
	if err != nil {
		return err
	}
	log.Println("response", string(b), "")
	cresp := &clientResponse{}
	err = json.Unmarshal(b, &cresp)
	if err != nil {
		return err
	}
	if cresp.Error != nil || cresp.Result == nil {
		x, ok := cresp.Error.(string)
		if !ok {
			return fmt.Errorf("invalid error %v", cresp.Error)
		}
		if x == "" {
			x = "unspecified error"
		}
		return fmt.Errorf(x)
	}
	return json.Unmarshal(*cresp.Result, resp)
}

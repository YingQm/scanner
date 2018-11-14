package main

import (
	"flag"
	"fmt"
	"github.com/YingQm/scanner/sync"
	"github.com/YingQm/scanner/config"
	"github.com/YingQm/scanner/send"
	"net"
	"strings"
	"time"
)

//关闭程序
var clo chan bool = make(chan bool)
var cfg *config.Config
var result chan string = make(chan string)
var scanResultOld string
var scanResult string

func RunScan(ipports []string) {
	for {
		scanResult = ""
		for i := 0; i < len(ipports); i++ {
			go scan(ipports[i])
		}

		for i := 0; i < len(ipports); i++ {
			s, ok := <-result
			if ok && len(s) > 0 {
				scanResult += (string)(s + "\n")
			}
		}

		if len(scanResult) > 0 && !StringSliceReflectEqual(scanResult, scanResultOld) {
			send.SendEmail(cfg, "提醒", "端口不开放:"+scanResult)
		}
		scanResultOld = scanResult

		time.Sleep((time.Duration)(cfg.IntervalTime) * time.Minute)
	}
}

func StringSliceReflectEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	as := strings.Split(a, "\n")
	bs := strings.Split(b, "\n")
	if len(as) != len(bs) {
		return false
	}

	for i := 0; i < len(as); i++ {
		bfind := false
		for j := 0; j < len(bs); j++ {
			if as[i] == bs[j] {
				bfind = true
				bs = append(bs[:j], bs[j+1:]...)
				break
			}
		}

		if !bfind {
			return false
		}
	}

	return true
}

func scan(address string) {
	_, err := net.Dial("tcp", address)
	if err != nil {
		//	fmt.Println(address)
		result <- address
	} else {
		result <- ""
	}
}

func SendToEmail() {
	for {
		if len(scanResultOld) > 0 {
			send.SendEmail(cfg, "提醒", "端口不开放:"+scanResultOld)
		}
		time.Sleep((time.Duration)(cfg.SendTime) * time.Minute)
	}
}

func RunSync() {
	for {
		syncresult := sync.EthSync(cfg.EthServiceAddr)
		if len(syncresult) > 0 {
			fmt.Println("syncresult", syncresult)
			send.SendEmail(cfg, "提醒", "节点不同步:"+syncresult)
		}
		time.Sleep((time.Duration)(cfg.SendTime) * time.Minute)
	}
}

func main() {
	configpath := flag.String("f", "config.toml", "configfile")
	cfg = config.InitCfg(*configpath)

	if len(cfg.IpPosts) > 0 {
		go RunScan(cfg.IpPosts)
		go SendToEmail()
	}

	if len(cfg.EthServiceAddr) > 0 {
		go RunSync()
	}

	//等待退出指令
	<-clo
	fmt.Println("Exit")
}

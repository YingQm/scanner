package main

import (
	"flag"
	"fmt"
<<<<<<< HEAD
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/wallet/monitor/config"
	lg "gitlab.33.cn/wallet/monitor/log"
	"gitlab.33.cn/wallet/monitor/types"
	"gitlab.33.cn/wallet/monitor/util"
	"net/http"
=======
	"github.com/YingQm/scanner/sync"
	"github.com/YingQm/scanner/config"
	"github.com/YingQm/scanner/send"
	"net"
	"strings"
	"time"
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
)

//关闭程序
var clo chan bool = make(chan bool)
var cfg *types.Config
var xdb *util.DbHandler
var cache *util.Cache
var configpath = flag.String("f", "./config.toml", "configfile")
var log = l.New("module", "monitor")
var rateLimiter *util.RateLimiter
var jobChan chan struct{}

<<<<<<< HEAD
func main() {
	flag.Parse()
	var err error
=======
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
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb

	cfg = config.InitCfg(*configpath)
	lg.SetFileLog(cfg.Log)
	log.Info("Log", "loglevel", cfg.Log.Loglevel, "logConsoleLevel", cfg.Log.LogConsoleLevel, "logFile", cfg.Log.LogFile,
		"ratelimit", cfg.RateLimit.LimitOpen, "maxConCurrentNum", cfg.RateLimit.MaxConCurrentNum,
		"ratelimit interval", cfg.RateLimit.TimeInterval, "ratelimit maxcount", cfg.RateLimit.MaxCount)
	log.Info("cfg", "noderpcuser", cfg.Rpcname, "noderpcpassword", cfg.Rpcpasswd, "omniuser", cfg.Omniname,
		"omnipasswd", cfg.Omnipasswd)

	xdb, err = util.NewMysql(cfg)
	if err != nil {
		panic(err)
	}

<<<<<<< HEAD
	limitRateOpen := cfg.RateLimit.LimitOpen
	if limitRateOpen {
		rateLimiter, err = util.NewRateLimiter(cfg.RateLimit.TimeInterval, cfg.RateLimit.MaxCount)
		if err != nil {
			panic(err)
		}
		if cfg.RateLimit.MaxConCurrentNum > 0 {
			jobChan = make(chan struct{}, cfg.RateLimit.MaxConCurrentNum)
=======
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
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
		}
	}
	cache = util.NewCache(cfg, xdb)
	go cache.UpdateCache()
	go RunCheck()

	http.HandleFunc("/getblockheight", GetBlockHeightHandler)
	http.ListenAndServe(":9988", nil)

<<<<<<< HEAD
=======
	if len(cfg.IpPosts) > 0 {
		go RunScan(cfg.IpPosts)
		go SendToEmail()
	}

	if len(cfg.EthServiceAddr) > 0 {
		go RunSync()
	}

>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
	//等待退出指令
	<-clo
	fmt.Println("Exit")
}

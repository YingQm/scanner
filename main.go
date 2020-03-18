package main

import (
	"flag"
	"fmt"
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/wallet/monitor/config"
	lg "gitlab.33.cn/wallet/monitor/log"
	"gitlab.33.cn/wallet/monitor/types"
	"gitlab.33.cn/wallet/monitor/util"
	"net/http"
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

func main() {
	flag.Parse()
	var err error

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

	limitRateOpen := cfg.RateLimit.LimitOpen
	if limitRateOpen {
		rateLimiter, err = util.NewRateLimiter(cfg.RateLimit.TimeInterval, cfg.RateLimit.MaxCount)
		if err != nil {
			panic(err)
		}
		if cfg.RateLimit.MaxConCurrentNum > 0 {
			jobChan = make(chan struct{}, cfg.RateLimit.MaxConCurrentNum)
		}
	}
	cache = util.NewCache(cfg, xdb)
	go cache.UpdateCache()
	go RunCheck()

	http.HandleFunc("/getblockheight", GetBlockHeightHandler)
	http.ListenAndServe(":9988", nil)

	//等待退出指令
	<-clo
	fmt.Println("Exit")
}
